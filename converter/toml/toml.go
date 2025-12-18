package toml

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/theobori/nix-converter/internal/common"
)

type TOMLVisitor struct {
	i    common.Indentation
	node any
}

const MaxNixNumber = 9223372036854775807 // 64 bits signed - 1
const MinNixNumber = -9223372036854775807

func NewTOMLVisitor(node any) *TOMLVisitor {
	return &TOMLVisitor{
		i:    *common.NewDefaultIndentation(),
		node: node,
	}
}

func (t *TOMLVisitor) visitMap(node map[string]any) (string, error) {
	e := []string{}

	for key, value := range node {
		t.i.Indent()
		valueResult, err := t.visit(value)
		if err != nil {
			return "", err
		}

		e = append(e, t.i.IndentValue()+string(key)+" = "+valueResult+";")
		t.i.UnIndent()
	}

	return "{\n" + strings.Join(e, "\n") + "\n" + t.i.IndentValue() + "}", nil
}

func (t *TOMLVisitor) visitArray(node []any) (string, error) {
	e := []string{}
	for _, item := range node {
		t.i.Indent()
		itemResult, err := t.visit(item)
		if err != nil {
			return "", err
		}

		e = append(e, t.i.IndentValue()+itemResult)
		t.i.UnIndent()
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
		return "\"" + v.String() + "\"", nil
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
		return "\"" + common.EscapeNixString(v.(string)) + "\"", nil
	}
}

func (t *TOMLVisitor) Visit() (string, error) {
	return t.visit(t.node)
}

func ToNix(data string) (string, error) {
	var node map[string]any

	err := toml.Unmarshal([]byte(data), &node)
	if err != nil {
		return "", err
	}

	return NewTOMLVisitor(node).Visit()
}
