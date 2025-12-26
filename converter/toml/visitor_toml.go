package toml

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
)

type TOMLVisitor struct {
	i       common.Indentation
	node    any
	options *converter.ConverterOptions
}

const (
	MaxNixNumber = 9223372036854775807 // 64 bits signed - 1
	MinNixNumber = -9223372036854775807
)

func NewTOMLVisitor(node any, options *converter.ConverterOptions) *TOMLVisitor {
	return &TOMLVisitor{
		i:       *common.NewDefaultIndentation(),
		node:    node,
		options: options,
	}
}

func (t *TOMLVisitor) visitMap(node map[string]any) (string, error) {
	if len(node) == 0 {
		return "{}", nil
	}

	e := []string{}

	for key, value := range node {
		t.i.Indent()
		left := nix.MakeNameSafe(string(key), t.options.UnsafeKeys)
		valueResult, err := t.visit(value)
		if err != nil {
			return "", err
		}

		e = append(e, t.i.IndentValue()+left+" = "+valueResult+";")
		t.i.UnIndent()
	}

	if t.options.SortIterators.SortHashmap {
		slices.Sort(e)
	}

	return "{\n" + strings.Join(e, "\n") + "\n" + t.i.IndentValue() + "}", nil
}

func (t *TOMLVisitor) visitArray(node []any) (string, error) {
	if len(node) == 0 {
		return "[]", nil
	}

	e := []string{}
	for _, item := range node {
		t.i.Indent()
		itemResult, err := t.visit(item)
		if err != nil {
			return "", err
		}

		e = append(e, t.i.IndentValue()+nix.MakeElementSafe(itemResult))
		t.i.UnIndent()
	}

	if t.options.SortIterators.SortList {
		slices.Sort(e)
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + t.i.IndentValue() + "]", nil
}

func (t *TOMLVisitor) visit(node any) (string, error) {
	switch v := node.(type) {
	case map[string]interface{}:
		return t.visitMap(v)
	case []any:
		return t.visitArray(v)
	case time.Time:
		return common.MakeStringSafe(v.String()), nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "null", nil
		}

		if v < MinNixNumber || v > MaxNixNumber {
			return "", fmt.Errorf("number out of range")
		}

		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int64:
		if v < MinNixNumber || v > MaxNixNumber {
			return "", fmt.Errorf("number out of range")
		}

		return fmt.Sprintf("%d", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		s := v.(string)
		if strings.Contains(s, "\n") {
			return common.MakeIndentedString(s, t.i.IndentValue()), nil
		}
		return common.MakeStringSafe(s), nil
	}
}

func (t *TOMLVisitor) Visit() (string, error) {
	return t.visit(t.node)
}

func ToNix(data string, options *converter.ConverterOptions) (string, error) {
	var node map[string]any

	err := toml.Unmarshal([]byte(data), &node)
	if err != nil {
		return "", err
	}

	return NewTOMLVisitor(node, options).Visit()
}
