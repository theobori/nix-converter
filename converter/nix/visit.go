package nix

import (
	"strings"

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

// ProcessIndentedString processes a Nix indented string (IStringNode)
// by stripping common leading whitespace and handling escape sequences
func ProcessIndentedString(raw string) string {
	if raw == "" {
		return ""
	}

	lines := strings.Split(raw, "\n")
	
	// Remove first line if it's empty or only whitespace (common for indented strings)
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	
	// Remove last line if it's empty or only whitespace
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	
	if len(lines) == 0 {
		return ""
	}
	
	// Find the minimum indentation (ignoring empty lines)
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
	
	// Strip the common indentation and process escape sequences
	result := make([]string, len(lines))
	for i, line := range lines {
		// For lines that are only whitespace, make them empty
		if strings.TrimSpace(line) == "" {
			result[i] = ""
			continue
		}
		
		// Strip common indentation
		if len(line) >= minIndent {
			line = line[minIndent:]
		}
		
		// Process Nix indented string escape sequences
		// ''${ -> ${
		line = strings.ReplaceAll(line, "''${", "${")
		// ''\ -> '
		line = strings.ReplaceAll(line, "''\\", "'")
		// ''' -> ''
		line = strings.ReplaceAll(line, "'''", "''")
		
		result[i] = line
	}
	
	return strings.Join(result, "\n") + "\n"
}
