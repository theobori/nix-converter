package xml

import (
	"encoding/xml"
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

type xmlNode struct {
	valueNode *parser.Node
	children  map[string]*xmlNode
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
	switch node.Type {
	case parser.IDNode:
		return nix.VisitID(n.p, node)
	case parser.StringNode:
		if len(node.Nodes) == 0 {
			return "", nil
		}
		return n.p.TokenString(node.Nodes[0].Tokens[0]), nil
	default:
		return "", fmt.Errorf("unsupported key node type: %s", node.Type.String())
	}
}

func (n *NixVisitor) buildXML(node *xmlNode) (string, error) {
	var parts []string

	for _, k := range node.order {
		child := node.children[k]
		open := n.i.IndentValue() + "<" + k + ">"
		closeTag := n.i.IndentValue() + "</" + k + ">"

		n.i.Indent()
		var body string
		var err error
		if child.valueNode != nil {
			body, err = n.visit(child.valueNode)
		} else {
			body, err = n.buildXML(child)
		}
		n.i.UnIndent()
		if err != nil {
			return "", err
		}

		parts = append(parts, open+"\n"+body+"\n"+closeTag)
	}

	if n.options.SortIterators.SortHashmap {
		slices.Sort(parts)
	}

	return strings.Join(parts, "\n"), nil
}

func (n *NixVisitor) visitSet(node *parser.Node) (string, error) {
	root := &xmlNode{children: make(map[string]*xmlNode)}

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
				curr.children[k] = &xmlNode{children: make(map[string]*xmlNode)}
				curr.order = append(curr.order, k)
			}
			if i == len(keys)-1 {
				curr.children[k].valueNode = valueNode
			} else {
				curr = curr.children[k]
			}
		}
	}

	return n.buildXML(root)
}

func (n *NixVisitor) visitList(node *parser.Node) (string, error) {
	e := []string{}
	for _, child := range node.Nodes {
		n.i.Indent()
		s, err := n.visit(child)
		n.i.UnIndent()
		if err != nil {
			return "", err
		}

		e = append(e, n.i.IndentValue()+"<item>\n"+s+"\n"+n.i.IndentValue()+"</item>")
	}

	if n.options.SortIterators.SortList {
		slices.Sort(e)
	}

	return strings.Join(e, "\n"), nil
}

func (n *NixVisitor) visitUnaryNegative(node *parser.Node) (string, error) {
	result, err := n.visitNode(node.Nodes[0])
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
	return makeXMLString(processed), nil
}

func (n *NixVisitor) visit(node *parser.Node) (string, error) {
	s, err := n.visitNode(node)
	if err != nil {
		return "", err
	}

	switch node.Type {
	case parser.SetNode, parser.ListNode, parser.ParensNode:
		return s, nil
	}

	return n.i.IndentValue() + s, nil
}

func (n *NixVisitor) visitNode(node *parser.Node) (string, error) {
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

func makeXMLString(s string) string {
	b, _ := xml.Marshal(s)
	return string(b)
}
