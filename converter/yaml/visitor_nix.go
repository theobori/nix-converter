package yaml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/orivej/go-nix/nix/parser"
	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
)

type NixVisitor struct {
	i       common.Indentation
	p       *parser.Parser
	node    *parser.Node
	options *converter.ConverterOptions
}

func NewNixVisitor(p *parser.Parser, node *parser.Node, options *converter.ConverterOptions) *NixVisitor {
	return &NixVisitor{
		i:       *common.NewDefaultIndentation(),
		node:    node,
		p:       p,
		options: options,
	}
}

func (n *NixVisitor) visitSet(node *parser.Node) (string, error) {
	if len(node.Nodes) == 0 {
		return "{}", nil
	}

	// Build nested structure from attribute paths while preserving order
	type nodeInfo struct {
		keys  []string
		value *parser.Node
	}
	items := []nodeInfo{}

	for _, child := range node.Nodes {
		attrPathNode := child.Nodes[0]
		valueNode := child.Nodes[1]

		var keys []string
		if attrPathNode.Type == parser.AttrPathNode && len(attrPathNode.Nodes) > 1 {
			// Multi-level path like package.meta.desc
			for _, keyNode := range attrPathNode.Nodes {
				key, err := n.visit(keyNode)
				if err != nil {
					return "", err
				}
				keys = append(keys, key)
			}
		} else {
			// Single key
			key, err := n.visit(attrPathNode)
			if err != nil {
				return "", err
			}
			keys = []string{key}
		}

		items = append(items, nodeInfo{keys: keys, value: valueNode})
	}

	// Build nested map structure with order preservation
	type mapNode struct {
		children  map[string]*mapNode
		childKeys []string // Preserve insertion order
		value     *parser.Node
	}
	root := &mapNode{children: make(map[string]*mapNode), childKeys: []string{}}

	for _, item := range items {
		current := root
		for i, key := range item.keys {
			if i == len(item.keys)-1 {
				// Leaf node
				if current.children[key] == nil {
					current.children[key] = &mapNode{value: item.value}
					current.childKeys = append(current.childKeys, key)
				} else {
					current.children[key].value = item.value
				}
			} else {
				// Intermediate node
				if current.children[key] == nil {
					current.children[key] = &mapNode{children: make(map[string]*mapNode), childKeys: []string{}}
					current.childKeys = append(current.childKeys, key)
				} else if current.children[key].children == nil {
					current.children[key].children = make(map[string]*mapNode)
					current.children[key].childKeys = []string{}
				}
				current = current.children[key]
			}
		}
	}

	// Convert to YAML output
	var buildYAML func(m *mapNode) ([]string, error)
	buildYAML = func(m *mapNode) ([]string, error) {
		e := []string{}
		keys := m.childKeys
		if n.options.SortIterators.SortHashmap {
			keys = make([]string, len(m.childKeys))
			copy(keys, m.childKeys)
			slices.Sort(keys)
		}

		for _, key := range keys {
			node := m.children[key]
			safeKey := MakeNameSafe(key, n.options.UnsafeKeys)
			keyString := n.i.IndentValue() + safeKey + ":"

			var valueNode *parser.Node
			if node.value != nil {
				valueNode = node.value
			}

			if valueNode != nil {
				// Has a value node
				switch valueNode.Type {
				case parser.SetNode, parser.ListNode:
					if len(valueNode.Nodes) == 0 {
						value, err := n.visit(valueNode)
						if err != nil {
							return nil, err
						}
						e = append(e, keyString+" "+value)
					} else {
						n.i.Indent()
						value, err := n.visit(valueNode)
						if err != nil {
							return nil, err
						}
						n.i.UnIndent()
						e = append(e, keyString+"\n"+value)
					}
				default:
					value, err := n.visit(valueNode)
					if err != nil {
						return nil, err
					}
					e = append(e, keyString+" "+value)
				}
			} else if len(node.children) > 0 {
				// Has nested children (from expanded paths)
				n.i.Indent()
				nested, err := buildYAML(node)
				if err != nil {
					return nil, err
				}
				n.i.UnIndent()
				e = append(e, keyString+"\n"+strings.Join(nested, "\n"))
			}
		}

		return e, nil
	}

	lines, err := buildYAML(root)
	if err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func (n *NixVisitor) visitList(node *parser.Node) (string, error) {
	if len(node.Nodes) == 0 {
		return "[]", nil
	}

	e := []string{}
	for _, child := range node.Nodes {
		n.i.Indent()

		s, err := n.visit(child)
		if err != nil {
			return "", err
		}

		n.i.UnIndent()

		e = append(e, n.i.IndentValue()+"- "+strings.TrimLeft(s, " "))
	}

	if n.options.SortIterators.SortList {
		slices.Sort(e)
	}

	return strings.Join(e, "\n"), nil
}

func (n *NixVisitor) visitUnaryNegative(node *parser.Node) (string, error) {
	result, err := n.visit(node.Nodes[0])
	if err != nil {
		return "", nil
	}

	return "-" + result, nil
}

func (n *NixVisitor) visitParens(node *parser.Node) (string, error) {
	// Empty Nix parens are not allowed
	return n.visit(node.Nodes[0])
}

func (n *NixVisitor) visitIndentedString(node *parser.Node) (string, error) {
	if len(node.Nodes) == 0 || len(node.Nodes[0].Tokens) == 0 {
		return "\"\"", nil
	}

	raw := n.p.TokenString(node.Nodes[0].Tokens[0])
	processed := nix.ProcessIndentedString(raw)

	// If it's a single line, just return as quoted string
	if !strings.Contains(processed, "\n") {
		return common.MakeStringSafe(processed), nil
	}

	// Use YAML block scalar for multiline strings
	lines := strings.Split(strings.TrimSuffix(processed, "\n"), "\n")
	result := "|-\n"
	n.i.Indent()
	for _, line := range lines {
		result += n.i.IndentValue() + line + "\n"
	}
	n.i.UnIndent()

	return strings.TrimSuffix(result, "\n"), nil
}

func (n *NixVisitor) visit(node *parser.Node) (string, error) {
	switch node.Type {
	case parser.SetNode:
		return n.visitSet(node)
	case parser.ListNode:
		return n.visitList(node)
	case parser.TextNode:
		return nix.VisitText(n.p, node)
	case parser.AttrPathNode:
		return nix.VisitAttrPathNode(n.p, node)
	case parser.IDNode:
		return nix.VisitID(n.p, node)
	case parser.StringNode:
		return nix.VisitString(n.p, node)
	case parser.IStringNode:
		return n.visitIndentedString(node)
	case parser.IntNode:
		return nix.VisitInt(n.p, node)
	case parser.FloatNode:
		return nix.VisitFloat(n.p, node)
	case parser.OpNode + 57378: // The negative unary operator
		return n.visitUnaryNegative(node)
	case parser.ApplyNode:
		return nix.VisitApply(n.p, node)
	case parser.ParensNode:
		return n.visitParens(node)
	default:
		return "", fmt.Errorf("unsupported node type: %s", node.Type.String())
	}
}

func (n *NixVisitor) Visit() (string, error) {
	return n.visit(n.node)
}

func FromNix(data string, options *converter.ConverterOptions) (string, error) {
	p, err := parser.ParseString(data)
	if err != nil {
		return "", err
	}

	out, err := NewNixVisitor(p, p.Result, options).Visit()
	if err != nil {
		return "", err
	}

	return out, nil
}
