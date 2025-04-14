package yaml

import (
	"fmt"
	"strings"

	"github.com/orivej/go-nix/nix/parser"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
)

const YAMLIndentSize = 2

func isYAMLString(s string) bool {
	for i := range s {
		if !common.IsAlphaNumeric(s[i]) && s[i] != ' ' {
			return false
		}
	}

	return true
}

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
	if len(node.Nodes) == 0 {
		return "{}", nil
	}

	e := []string{}
	for _, child := range node.Nodes {
		key, err := n.visit(child.Nodes[0])
		if err != nil {
			return "", err
		}

		valueNode := child.Nodes[1]
		keyString := n.i.IndentValue() + key + ": "

		switch valueNode.Type {
		case parser.SetNode, parser.ListNode:
			if len(valueNode.Nodes) == 0 {
				value, err := n.visit(valueNode)
				if err != nil {
					return "", err
				}
				e = append(e, keyString+value)
			} else {
				e = append(e, keyString)

				n.i.Indent()
				value, err := n.visit(valueNode)
				if err != nil {
					return "", err
				}
				n.i.UnIndent()

				e = append(e, value)
			}
		default:
			value, err := n.visit(valueNode)
			if err != nil {
				return "", err
			}
			e = append(e, keyString+value)
		}
	}

	return strings.Join(e, "\n"), nil
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

	return strings.Join(e, "\n"), nil
}

func (n *NixVisitor) visitString(node *parser.Node) (string, error) {
	token := n.p.TokenString(node.Nodes[0].Tokens[0])

	if !isYAMLString(token) {
		token = "\"" + token + "\""
	}

	return token, nil
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
		return n.visitString(node)
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
