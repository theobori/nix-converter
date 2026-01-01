package nix

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/orivej/go-nix/nix/parser"
)

type NixVisitor struct {
	p    *parser.Parser
	node *parser.Node
}

func NewNixVisitor(p *parser.Parser, node *parser.Node) *NixVisitor {
	return &NixVisitor{
		node: node,
		p:    p,
	}
}

func (n *NixVisitor) visitSet(node *parser.Node) (any, error) {
	out := map[string]any{}

	for _, child := range node.Nodes {
		// Handle nested attribute paths (e.g., package.meta.desc)
		attrPathNode := child.Nodes[0]
		value, err := n.visit(child.Nodes[1])
		if err != nil {
			return map[string]any{}, err
		}

		// Extract all keys in the path
		keys := make([]string, len(attrPathNode.Nodes))
		for i, keyNode := range attrPathNode.Nodes {
			keyVal, err := n.visit(keyNode)
			if err != nil {
				return map[string]any{}, err
			}
			keyStr, ok := keyVal.(string)
			if !ok {
				return nil, fmt.Errorf("unable to convert attr key to string")
			}
			keys[i] = keyStr
		}

		// Create nested structure
		current := out
		for i := 0; i < len(keys)-1; i++ {
			if _, exists := current[keys[i]]; !exists {
				current[keys[i]] = map[string]any{}
			}
			var ok bool
			current, ok = current[keys[i]].(map[string]any)
			if !ok {
				return nil, fmt.Errorf("conflicting attribute path")
			}
		}
		current[keys[len(keys)-1]] = value
	}

	return out, nil
}

func (n *NixVisitor) visitList(node *parser.Node) (any, error) {
	out := []any{}

	for _, child := range node.Nodes {
		item, err := n.visit(child)
		if err != nil {
			return []any{}, err
		}

		out = append(out, item)
	}

	return out, nil
}

func (n *NixVisitor) visitUnaryNegative(node *parser.Node) (any, error) {
	result, err := n.visit(node.Nodes[0])
	if err != nil {
		return "", err
	}

	t := reflect.TypeOf(result)
	switch t.Kind() {
	case reflect.Int64:
		return -result.(int64), nil
	case reflect.Float64:
		return -result.(float64), nil
	default:
		return nil, fmt.Errorf("unsupported go type: %s", t.Kind().String())
	}
}

func (n *NixVisitor) visitString(node *parser.Node) (any, error) {
	if len(node.Nodes) == 0 {
		return "", nil
	}

	return n.p.TokenString(node.Nodes[0].Tokens[0]), nil
}

func (n *NixVisitor) visitIndentedString(node *parser.Node) (string, error) {
	if len(node.Nodes) == 0 || len(node.Nodes[0].Tokens) == 0 {
		return "", nil
	}
	raw := n.p.TokenString(node.Nodes[0].Tokens[0])
	return ProcessIndentedString(raw), nil
}

func (n *NixVisitor) visitParens(node *parser.Node) (any, error) {
	// Empty Nix parens are not allowed
	return n.visit(node.Nodes[0])
}

func (n *NixVisitor) visit(node *parser.Node) (any, error) {
	switch node.Type {
	case parser.SetNode:
		return n.visitSet(node)
	case parser.ListNode:
		return n.visitList(node)
	case parser.TextNode:
		return VisitText(n.p, node)
	case parser.AttrPathNode:
		return VisitAttrPathNode(n.p, node)
	case parser.IDNode:
		return VisitID(n.p, node)
	case parser.StringNode:
		return n.visitString(node)
	case parser.IStringNode:
		return n.visitIndentedString(node)
	case parser.IntNode:
		return VisitIntRaw(n.p, node)
	case parser.FloatNode:
		return VisitFloatRaw(n.p, node)
	case parser.OpNode + 57378: // The negative unary operator
		return n.visitUnaryNegative(node)
	case parser.ApplyNode:
		return VisitApplyRaw(n.p, node)
	case parser.ParensNode:
		return n.visitParens(node)
	default:
		return nil, fmt.Errorf("unsupported node type: %s", node.Type)
	}
}

func ProcessIndentedString(raw string) string {
	if !strings.Contains(raw, "\n") {
		return raw
	}

	lines := strings.Split(raw, "\n")

	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent == -1 {
		minIndent = 0
	}

	result := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			result[i] = ""
			continue
		}

		if len(line) >= minIndent {
			line = line[minIndent:]
		}

		// Handle escapes
		line = strings.ReplaceAll(line, "''${", "${")
		line = strings.ReplaceAll(line, "''\\", "'")
		result[i] = line
	}

	return strings.Join(result, "\n")
}

func (n *NixVisitor) Visit() (any, error) {
	return n.visit(n.node)
}

func GoValue(data string) (any, error) {
	p, err := parser.ParseString(data)
	if err != nil {
		return nil, err
	}

	out, err := NewNixVisitor(p, p.Result).Visit()
	if err != nil {
		return nil, err
	}

	return out, nil
}
