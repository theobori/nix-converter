package json

import "github.com/theobori/nix-converter/converter"

type JSONConverter struct {
	data    string
	options *converter.ConverterOptions
}

func NewJSONConverter(data string, options *converter.ConverterOptions) *JSONConverter {
	return &JSONConverter{
		data,
		options,
	}
}

func (j *JSONConverter) FromNix() (string, error) {
	return FromNix(j.data, j.options)
}

func (j *JSONConverter) ToNix() (string, error) {
	return ToNix(j.data, j.options)
}

func (j *JSONConverter) Type() string {
	return "json"
}
