package json

import (
	"slices"
	"strings"

	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
	"github.com/valyala/fastjson"
)

type JSONVisitor struct {
	i       common.Indentation
	value   *fastjson.Value
	options *converter.ConverterOptions
}

func NewJSONVisitor(value *fastjson.Value, options *converter.ConverterOptions) *JSONVisitor {
	return &JSONVisitor{
		i:       *common.NewDefaultIndentation(),
		value:   value,
		options: options,
	}
}

func (j *JSONVisitor) visitObject(value *fastjson.Value) string {
	o, _ := value.Object()

	e := []string{}
	o.Visit(func(key []byte, v *fastjson.Value) {
		j.i.Indent()
		left := nix.SafeName(string(key))
		right := j.visit(v)

		e = append(e, j.i.IndentValue()+left+" = "+right+";")
		j.i.UnIndent()
	})

	if j.options.SortIterators.SortHashmap {
		slices.Sort(e)
	}

	return "{\n" + strings.Join(e, "\n") + "\n" + j.i.IndentValue() + "}"
}

func (j *JSONVisitor) visitArray(value *fastjson.Value) string {
	arr, _ := value.Array()

	e := []string{}
	for _, item := range arr {
		j.i.Indent()
		element := nix.SafeElement(j.visit(item))

		e = append(e, j.i.IndentValue()+element)
		j.i.UnIndent()
	}

	if j.options.SortIterators.SortList {
		slices.Sort(e)
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + j.i.IndentValue() + "]"
}

func (j *JSONVisitor) visitString(value *fastjson.Value) string {
	// v := value.String()
	return value.String()
}

func (j *JSONVisitor) visitNumber(value *fastjson.Value) string {
	return value.String()
}

func (j *JSONVisitor) visitFalse(_ *fastjson.Value) string {
	return "false"
}

func (j *JSONVisitor) visitTrue(_ *fastjson.Value) string {
	return "true"
}

func (j *JSONVisitor) visitNull(_ *fastjson.Value) string {
	return "null"
}

func (j *JSONVisitor) visit(value *fastjson.Value) string {
	switch value.Type() {
	case fastjson.TypeObject:
		return j.visitObject(value)
	case fastjson.TypeArray:
		return j.visitArray(value)
	case fastjson.TypeString:
		return j.visitString(value)
	case fastjson.TypeNumber:
		return j.visitNumber(value)
	case fastjson.TypeFalse:
		return j.visitFalse(value)
	case fastjson.TypeTrue:
		return j.visitTrue(value)
	case fastjson.TypeNull:
		return j.visitNull(value)
	default:
		return ""
	}
}

func (j *JSONVisitor) Visit() string {
	return j.visit(j.value)
}

func ToNix(data string, options *converter.ConverterOptions) (string, error) {
	v, err := fastjson.Parse(data)
	if err != nil {
		return "", err
	}

	out := NewJSONVisitor(v, options).Visit()

	return out, nil
}
