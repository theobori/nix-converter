package json

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/orivej/go-nix/nix/parser"
	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
)

const IndentSize = 2

type NixVisitor struct {
	i       common.Indentation
	p       *parser.Parser
	node    *parser.Node
	options *converter.ConverterOptions
}

type jsonNode struct {
	valueNode *parser.Node
	children  map[string]*jsonNode
	order     []string
}

func NewNixVisitor(p *parser.Parser, node *parser.Node, options *converter.ConverterOptions) *NixVisitor {
	return &NixVisitor{
		i:       *common.NewDefaultIndentation(),
		node:    node,
		p:       p,
		options: options,
	}
}

func (n *NixVisitor) visitKey(node *parser.Node) (string, error) {
	if node.Type == parser.IDNode {
		return n.p.TokenString(node.Tokens[0]), nil
	}
	if node.Type == parser.StringNode {
		if len(node.Nodes) == 0 {
			return "", nil
		}
		return n.p.TokenString(node.Nodes[0].Tokens[0]), nil
	}
	return "", fmt.Errorf("unsupported key node type: %s", node.Type.String())
}

func (n *NixVisitor) buildJSON(node *jsonNode) (string, error) {
	var parts []string

	keys := node.order

	n.i.Indent()
	for _, k := range keys {
		child := node.children[k]
		keyStr := n.i.IndentValue() + makeJSONString(k) + ": "
		if child.valueNode != nil {
			valStr, err := n.visit(child.valueNode)
			if err != nil {
				return "", err
			}
			parts = append(parts, keyStr+valStr)
		} else {
			childStr, err := n.buildJSON(child)
			if err != nil {
				return "", err
			}
			parts = append(parts, keyStr+childStr)
		}
	}
	n.i.UnIndent()

	return "{\n" + strings.Join(parts, ",\n") + "\n" + n.i.IndentValue() + "}", nil
}

func (n *NixVisitor) visitSet(node *parser.Node) (string, error) {
	if len(node.Nodes) == 0 {
		return "{}", nil
	}

	root := &jsonNode{children: make(map[string]*jsonNode)}

	for _, child := range node.Nodes {

		attrPathNode := child.Nodes[0]
		valueNode := child.Nodes[1]

		var keys []string
		for _, kNode := range attrPathNode.Nodes {
			k, err := n.visitKey(kNode)
			if err != nil {
				return "", err
			}
			keys = append(keys, k)
		}

		curr := root
		for i, k := range keys {
			if curr.children[k] == nil {
				curr.children[k] = &jsonNode{children: make(map[string]*jsonNode)}
				curr.order = append(curr.order, k)
			}
			if i == len(keys)-1 {
				curr.children[k].valueNode = valueNode
			} else {
				curr = curr.children[k]
			}
		}
	}

	return n.buildJSON(root)
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

		e = append(e, n.i.IndentValue()+s)
		n.i.UnIndent()
	}

	if n.options.SortIterators.SortList {
		slices.Sort(e)
	}

	return "[\n" + strings.Join(e, ",\n") + "\n" + n.i.IndentValue() + "]", nil
}

func (n *NixVisitor) visitUnaryNegative(node *parser.Node) (string, error) {
	result, err := n.visit(node.Nodes[0])
	if err != nil {
		return "", err
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
	return makeJSONString(processed), nil
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

func makeJSONString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
