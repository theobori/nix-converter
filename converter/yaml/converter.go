package yaml

import "github.com/theobori/nix-converter/converter"

type YAMLConverter struct {
	data    string
	options *converter.ConverterOptions
}

func NewYAMLConverter(data string, options *converter.ConverterOptions) *YAMLConverter {
	return &YAMLConverter{
		data,
		options,
	}
}

func (y *YAMLConverter) FromNix() (string, error) {
	return FromNix(y.data, y.options)
}

func (y *YAMLConverter) ToNix() (string, error) {
	return ToNix(y.data, y.options)
}

func (y *YAMLConverter) Type() string {
	return "yaml"
}
