package yaml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/internal/common"
	"gopkg.in/yaml.v3"
)

type YAMLVisitor struct {
	anchors map[string]string
	i       common.Indentation
	node    *yaml.Node
	options *converter.ConverterOptions
}

func NewYAMLVisitor(node *yaml.Node, options *converter.ConverterOptions) *YAMLVisitor {
	return &YAMLVisitor{
		anchors: make(map[string]string),
		i:       *common.NewDefaultIndentation(),
		node:    node,
		options: options,
	}
}

func newAnchorIndentation() *common.Indentation {
	i := common.NewDefaultIndentation()
	i.Indent()
	return i
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

	if y.options.SortIterators.SortHashmap {
		slices.Sort(e)
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

	if y.options.SortIterators.SortList {
		slices.Sort(e)
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + y.i.IndentValue() + "]"
}

func (y *YAMLVisitor) visitScalar(node *yaml.Node) string {
	if node.Tag == "!!int" || node.Tag == "!!float" || node.Tag == "!!bool" {
		return node.Value
	}

	return "\"" + node.Value + "\""
}

func (y *YAMLVisitor) visitAlias(node *yaml.Node) string {
	return node.Alias.Anchor
}

func (y *YAMLVisitor) visit(node *yaml.Node) string {
	var output string
	indent := y.i

	if node.Anchor != "" {
		y.i = *newAnchorIndentation()
	}

	switch node.Kind {
	case yaml.MappingNode:
		output = y.visitMapping(node)
	case yaml.SequenceNode:
		output = y.visitSequence(node)
	case yaml.ScalarNode:
		output = y.visitScalar(node)
	case yaml.AliasNode:
		output = y.visitAlias(node)
	default:
		output = ""
	}

	y.i = indent
	if node.Anchor == "" {
		return output
	}

	y.anchors[node.Anchor] = output
	return node.Anchor
}

func (y *YAMLVisitor) Visit() string {
	firstPass := y.visit(y.node)
	if len(y.anchors) == 0 {
		return firstPass
	}

	secondPass := "let\n"
	y.i = *newAnchorIndentation()
	for k, v := range y.anchors {
		secondPass += y.i.IndentValue() + k + " = " + v + ";\n"
	}
	secondPass += "in " + firstPass

	return secondPass
}

func ToNix(data string, options *converter.ConverterOptions) (string, error) {
	var node yaml.Node

	err := yaml.Unmarshal([]byte(data), &node)
	if err != nil {
		return "", err
	}

	if len(node.Content) == 0 {
		return "", fmt.Errorf("empty node")
	}

	out := NewYAMLVisitor(node.Content[0], options).Visit()

	return out, nil
}
