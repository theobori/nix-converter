package json

import (
	"strings"

	"github.com/theobori/nix-converter/converter/nix"
	"github.com/theobori/nix-converter/internal/common"
	"github.com/valyala/fastjson"
)

type JSONVisitor struct {
	indentLevel int
	indentValue string
	value       *fastjson.Value
}

func NewJSONVisitor(value *fastjson.Value) *JSONVisitor {
	return &JSONVisitor{
		value: value,
	}
}

func (tn *JSONVisitor) indent() {
	tn.indentValue, tn.indentLevel = common.Indent(tn.indentLevel, nix.IndentSize)
}

func (tn *JSONVisitor) unIndent() {
	tn.indentValue, tn.indentLevel = common.UnIndent(tn.indentLevel, nix.IndentSize)
}

func (tn *JSONVisitor) visitObject(value *fastjson.Value) string {
	o, _ := value.Object()

	e := []string{}
	o.Visit(func(key []byte, v *fastjson.Value) {
		tn.indent()
		e = append(e, tn.indentValue+string(key)+" = "+tn.visit(v)+";")
		tn.unIndent()
	})

	return "{\n" + strings.Join(e, "\n") + "\n" + tn.indentValue + "}"
}

func (tn *JSONVisitor) visitArray(value *fastjson.Value) string {
	arr, _ := value.Array()

	e := []string{}
	for _, item := range arr {
		tn.indent()
		e = append(e, tn.indentValue+tn.visit(item))
		tn.unIndent()
	}

	return "[\n" + strings.Join(e, "\n") + "\n" + tn.indentValue + "]"
}

func (tn *JSONVisitor) visitString(value *fastjson.Value) string {
	return value.String()
}

func (tn *JSONVisitor) visitNumber(value *fastjson.Value) string {
	return value.String()
}

func (tn *JSONVisitor) visitFalse(_ *fastjson.Value) string {
	return "false"
}

func (tn *JSONVisitor) visitTrue(_ *fastjson.Value) string {
	return "true"
}

func (tn *JSONVisitor) visitNull(_ *fastjson.Value) string {
	return "null"
}

func (tn *JSONVisitor) visit(value *fastjson.Value) string {
	switch value.Type() {
	case fastjson.TypeObject:
		return tn.visitObject(value)
	case fastjson.TypeArray:
		return tn.visitArray(value)
	case fastjson.TypeString:
		return tn.visitString(value)
	case fastjson.TypeNumber:
		return tn.visitNumber(value)
	case fastjson.TypeFalse:
		return tn.visitFalse(value)
	case fastjson.TypeTrue:
		return tn.visitTrue(value)
	case fastjson.TypeNull:
		return tn.visitNull(value)
	default:
		return ""
	}
}

func (tn *JSONVisitor) Eval() string {
	return tn.visit(tn.value)
}

func ToNix(data string) (string, error) {
	v, err := fastjson.Parse(data)
	if err != nil {
		return "", err
	}

	out := NewJSONVisitor(v).Eval()

	return out, nil
}
