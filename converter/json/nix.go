package json

import (
	"fmt"
	"strings"

	"github.com/orivej/go-nix/nix/parser"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
)

const IndentSize = 2

type NixVisitor struct {
	i    common.Indentation
	p    *parser.Parser
	node *parser.Node
}

func NewNixVisitor(p *parser.Parser, node *parser.Node) *NixVisitor {
	return &NixVisitor{
		i:    *common.NewDefaultIndentation(),
		node: node,
		p:    p,
	}
}

func (n *NixVisitor) visitSet(node *parser.Node) (string, error) {
	e := []string{}
	for _, child := range node.Nodes {
		n.i.Indent()
		key, err := n.visit(child.Nodes[0])
		if err != nil {
			return "", err
		}

		value, err := n.visit(child.Nodes[1])
		if err != nil {
			return "", err
		}

		e = append(e, n.i.IndentValue()+"\""+key+"\""+": "+value)
		n.i.UnIndent()
	}

	return "{\n" + strings.Join(e, ",\n") + "\n" + n.i.IndentValue() + "}", nil
}

func (n *NixVisitor) visitList(node *parser.Node) (string, error) {
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

	return "[\n" + strings.Join(e, ",\n") + "\n" + n.i.IndentValue() + "]", nil
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
	case parser.StringNode, parser.IStringNode:
		return nix.VisitString(n.p, node)
	case parser.IntNode:
		return nix.VisitInt(n.p, node)
	case parser.FloatNode:
		return nix.VisitFloat(n.p, node)
	default:
		return "", fmt.Errorf("unsupported node type: %s", node.Type.String())
	}
}

func (n *NixVisitor) Visit() (string, error) {
	return n.visit(n.node)
}

func FromNix(data string) (string, error) {
	p, err := parser.ParseString(data)
	if err != nil {
		return "", err
	}

	out, err := NewNixVisitor(p, p.Result).Visit()
	if err != nil {
		return "", err
	}

	return out, nil
}
