package nix

import (
	"strconv"

	"github.com/orivej/go-nix/nix/parser"
	"github.com/theobori/nix-converter/internal/common"
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
		if len(node.Nodes[0].Nodes) == 0 {
			return "", nil
		}
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
	if len(node.Nodes) == 0 {
		return common.MakeStringSafe(""), nil
	}

	token := p.TokenString(node.Nodes[0].Tokens[0])

	return common.MakeStringSafe(token), nil
}

func VisitInt(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Tokens[0]), nil
}

func VisitIntRaw(p *parser.Parser, node *parser.Node) (int64, error) {
	s, err := VisitInt(p, node)
	if err != nil {
		return 0.0, nil
	}

	return strconv.ParseInt(s, 10, 64)
}

func VisitFloat(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Tokens[0]), nil
}

func VisitFloatRaw(p *parser.Parser, node *parser.Node) (float64, error) {
	s, err := VisitFloat(p, node)
	if err != nil {
		return 0.0, nil
	}

	return strconv.ParseFloat(s, 64)
}

func VisitApply(p *parser.Parser, node *parser.Node) (string, error) {
	return p.TokenString(node.Nodes[0].Tokens[0]) + p.TokenString(node.Nodes[1].Tokens[0]), nil
}

func VisitApplyRaw(p *parser.Parser, node *parser.Node) (float64, error) {
	s, err := VisitApply(p, node)
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(s, 64)
}
