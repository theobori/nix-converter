package common

import "strings"

func IndentBase(indentLevel int, indentSize int, symbol string) (string, int) {
	level := indentLevel + indentSize

	return strings.Repeat(symbol, level), level
}

func Indent(indentLevel int, indentSize int) (string, int) {
	return IndentBase(indentLevel, indentSize, " ")
}

func UnIndent(indentLevel int, indentSize int) (string, int) {
	return IndentBase(indentLevel, -indentSize, " ")
}

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || IsNumeric(c)
}
