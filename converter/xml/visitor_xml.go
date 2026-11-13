package xml

import (
	"fmt"

	"github.com/clbanning/mxj/v2"
	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
)

type XMLVisitor struct {
	i       common.Indentation
	node    any
	options *converter.ConverterOptions
}

func NewXMLVisitor(node any, options *converter.ConverterOptions) *XMLVisitor {
	return &XMLVisitor{
		i:       *common.NewDefaultIndentation(),
		node:    node,
		options: options,
	}
}

func (t *XMLVisitor) visitMap(node map[string]any) (string, error) {
	if len(node) == 0 {
		return "{}", nil
	}

	e := []string{}

	for key, value := range node {
		t.i.Indent()
		left := nix.MakeNameSafe(key, t.options.UnsafeKeys)
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

func (t *XMLVisitor) visitArray(node []any) (string, error) {
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

func (t *XMLVisitor) visit(node any) (string, error) {
	switch v := node.(type) {
	case map[string]any:
		return t.visitMap(v)
	case mxj.Map:
		return t.visitMap(v)
	case []any:
		return t.visitArray(v)
	case time.Time:
		return common.MakeStringSafe(v.String()), nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return "null", nil
		}

		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case int64:
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

func (t *XMLVisitor) Visit() (string, error) {
	return t.visit(t.node)
}

func ToNix(data string, options *converter.ConverterOptions) (string, error) {
	mxj.PrependAttrWithHyphen(false)

	node, err := mxj.NewMapXml([]byte(data), true)
	if err != nil {
		return "", err
	}

	return NewXMLVisitor(node, options).Visit()
}
