package yaml

import (
	"fmt"
	"strings"

	"github.com/theobori/nix-converter/internal/common"
	"gopkg.in/yaml.v3"
)

type YAMLVisitor struct {
	i    common.Indentation
	node *yaml.Node
}

func NewYAMLVisitor(node *yaml.Node) *YAMLVisitor {
	return &YAMLVisitor{
		i:    *common.NewDefaultIndentation(),
		node: node,
	}
}

func (y *YAMLVisitor) visitMapping(node *yaml.Node) string {
	e := []string{}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]

		y.i.Indent()
		e = append(e, y.i.IndentValue()+key+" = "+y.visit(value)+";")
		y.i.UnIndent()
	}

	return "{\n" + strings.Join(e, "\n") + "\n" + y.i.IndentValue() + "}"
}

func (y *YAMLVisitor) visitSequence(node *yaml.Node) string {
	e := []string{}
	for _, item := range node.Content {
		y.i.Indent()
		e = append(e, y.i.IndentValue()+y.visit(item))
		y.i.UnIndent()
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + y.i.IndentValue() + "]"
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

	if len(node.Content) == 0 {
		return "", fmt.Errorf("empty node")
	}

	out := NewYAMLVisitor(node.Content[0]).Eval()

	return out, nil
}
