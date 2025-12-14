package toml

import "github.com/theobori/nix-converter/converter"

type TOMLConverter struct {
	data    string
	options *converter.ConverterOptions
}

func NewTOMLConverter(data string, options *converter.ConverterOptions) *TOMLConverter {
	return &TOMLConverter{
		data,
		options,
	}
}

func (t *TOMLConverter) FromNix() (string, error) {
	return FromNix(t.data, t.options)
}

func (t *TOMLConverter) ToNix() (string, error) {
	return ToNix(t.data, t.options)
}

func (t *TOMLConverter) Type() string {
	return "toml"
}
