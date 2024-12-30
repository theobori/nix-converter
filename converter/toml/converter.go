package toml

import "fmt"

type TOMLConverter struct {
	data string
}

func NewTOMLConverter(data string) *TOMLConverter {
	return &TOMLConverter{
		data,
	}
}

func (t *TOMLConverter) FromNix() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (t *TOMLConverter) ToNix() (string, error) {
	return ToNix(t.data)
}

func (t *TOMLConverter) Type() string {
	return "toml"
}
