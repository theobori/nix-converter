package nix

import (
	"github.com/orivej/go-nix/nix/parser"
)

const IndentSize = 2

func VisitText(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Tokens[0]), nil
}

func VisitAttrPathNode(p *parser.Parser, node *parser.Node) (string, error) {
	tokens := node.Nodes[0].Tokens

	var token int
	// Case when the value is between double quotes
	if len(tokens) > 1 {
		token = node.Nodes[0].Nodes[0].Tokens[0]
	} else {
		token = tokens[0]
	}

	return p.TokenString(token), nil
}

func VisitID(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Tokens[0]), nil
}

func VisitString(p *parser.Parser, node *parser.Node) (string, error) {
	return "\"" + p.TokenString(node.Nodes[0].Tokens[0]) + "\"", nil
}

func VisitInt(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Tokens[0]), nil
}

func VisitFloat(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Tokens[0]), nil
}
