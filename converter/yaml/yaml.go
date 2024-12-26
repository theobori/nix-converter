package yaml

import (
	"strings"

	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
	"gopkg.in/yaml.v3"
)

type YAMLVisitor struct {
	indentLevel int
	indentValue string
	node        *yaml.Node
}

func NewYAMLVisitor(node *yaml.Node) *YAMLVisitor {
	return &YAMLVisitor{
		node: node,
	}
}

func (y *YAMLVisitor) indent() {
	y.indentValue, y.indentLevel = common.Indent(y.indentLevel, nix.IndentSize)
}

func (y *YAMLVisitor) unIndent() {
	y.indentValue, y.indentLevel = common.UnIndent(y.indentLevel, nix.IndentSize)
}

func (y *YAMLVisitor) visitMapping(node *yaml.Node) string {
	e := []string{}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]

		y.indent()
		e = append(e, y.indentValue+key+" = "+y.visit(value)+";")
		y.unIndent()
	}

	return "{\n" + strings.Join(e, "\n") + "\n" + y.indentValue + "}"
}

func (y *YAMLVisitor) visitSequence(node *yaml.Node) string {
	e := []string{}
	for _, item := range node.Content {
		y.indent()
		e = append(e, y.indentValue+y.visit(item))
		y.unIndent()
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + y.indentValue + "]"
}

func (y *YAMLVisitor) visitScalar(node *yaml.Node) string {
	if node.Tag == "!!int" || node.Tag == "!!float" || node.Tag == "!!bool" {
		return node.Value
	}

	return "\"" + node.Value + "\""
}

func (y *YAMLVisitor) visit(node *yaml.Node) string {
	switch node.Kind {
	case yaml.MappingNode:
		return y.visitMapping(node)
	case yaml.SequenceNode:
		return y.visitSequence(node)
	case yaml.ScalarNode:
		return y.visitScalar(node)
	default:
		return ""
	}
}

func (y *YAMLVisitor) Eval() string {
	return y.visit(y.node)
}

func ToNix(data string) (string, error) {
	var node yaml.Node

	err := yaml.Unmarshal([]byte(data), &node)
	if err != nil {
		return "", err
	}

	out := NewYAMLVisitor(node.Content[0]).Eval()

	return out, nil
}
