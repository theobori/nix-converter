package nix

import (
	"fmt"

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
		return n.p.TokenString(node.Nodes[0].Tokens[0]), nil
	case parser.IStringNode:
		if len(node.Nodes) == 0 || len(node.Nodes[0].Tokens) == 0 {
			return "", nil
		}
		raw := n.p.TokenString(node.Nodes[0].Tokens[0])
		return ProcessIndentedString(raw), nil
	case parser.IntNode:
		return VisitInt(n.p, node)
	case parser.FloatNode:
		return VisitFloat(n.p, node)
	default:
		return nil, fmt.Errorf("unsupported node type: %s", node.Type.String())
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
