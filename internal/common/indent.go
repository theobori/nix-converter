package common

import "strings"

type Indentation struct {
	indentAmount int
	indentSize   int
	indentValue  string
	spaces       string
}

func NewIndentation(indentSize int) *Indentation {
	return &Indentation{
		indentAmount: 0,
		indentSize:   indentSize,
		indentValue:  "",
		spaces:       strings.Repeat(" ", indentSize),
	}
}

func NewDefaultIndentation() *Indentation {
	return NewIndentation(2)
}

func (i *Indentation) Indent() {
	i.indentAmount += 1
	i.indentValue += i.spaces
}

func (i *Indentation) UnIndent() {
	if i.indentAmount <= 0 {
		return
	}

	i.indentAmount -= 1
	i.indentValue = i.indentValue[:len(i.indentValue)-i.indentSize]
}

func (i *Indentation) IndentValue() string {
	return i.indentValue
}
