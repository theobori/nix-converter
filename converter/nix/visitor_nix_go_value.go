package nix

import (
	"fmt"
	"reflect"

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
		key, err := n.visit(child.Nodes[0])
		if err != nil {
			return map[string]any{}, err
		}

		keyString, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("unable to convert attr key to string")
		}

		value, err := n.visit(child.Nodes[1])
		if err != nil {
			return map[string]any{}, err
		}

		out[keyString] = value
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
		return "", nil
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
	case parser.StringNode, parser.IStringNode:
		return n.visitString(node)
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
