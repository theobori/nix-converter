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

func (fn *NixVisitor) visitSet(node *parser.Node) (string, error) {
	e := []string{}
	for _, child := range node.Nodes {
		fn.i.Indent()
		key, err := fn.visit(child.Nodes[0])
		if err != nil {
			return "", err
		}

		value, err := fn.visit(child.Nodes[1])
		if err != nil {
			return "", err
		}

		e = append(e, fn.i.IndentValue()+"\""+key+"\""+": "+value)
		fn.i.UnIndent()
	}

	return "{\n" + strings.Join(e, ",\n") + "\n" + fn.i.IndentValue() + "}", nil
}

func (fn *NixVisitor) visitList(node *parser.Node) (string, error) {
	e := []string{}
	for _, child := range node.Nodes {
		fn.i.Indent()
		s, err := fn.visit(child)
		if err != nil {
			return "", err
		}

		e = append(e, fn.i.IndentValue()+s)
		fn.i.UnIndent()
	}

	return "[\n" + strings.Join(e, ",\n") + "\n" + fn.i.IndentValue() + "]", nil
}

func (fn *NixVisitor) visit(node *parser.Node) (string, error) {
	switch node.Type {
	case parser.SetNode:
		return fn.visitSet(node)
	case parser.ListNode:
		return fn.visitList(node)
	case parser.TextNode:
		return nix.VisitText(fn.p, node)
	case parser.AttrPathNode:
		return nix.VisitAttrPathNode(fn.p, node)
	case parser.IDNode:
		return nix.VisitID(fn.p, node)
	case parser.StringNode:
		return nix.VisitString(fn.p, node)
	case parser.IntNode:
		return nix.VisitInt(fn.p, node)
	case parser.FloatNode:
		return nix.VisitFloat(fn.p, node)
	default:
		return "", fmt.Errorf("unauthorized node type: %s", node.Type.String())
	}
}

func (fn *NixVisitor) Eval() (string, error) {
	return fn.visit(fn.node)
}

func FromNix(data string) (string, error) {
	p, err := parser.ParseString(data)
	if err != nil {
		return "", err
	}

	out, err := NewNixVisitor(p, p.Result).Eval()
	if err != nil {
		return "", err
	}

	return out, nil
}
