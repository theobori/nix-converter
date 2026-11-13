package xml

import "github.com/theobori/nix-converter/converter"

type XMLConverter struct {
	data    string
	options *converter.ConverterOptions
}

func NewXMLConverter(data string, options *converter.ConverterOptions) *XMLConverter {
	return &XMLConverter{
		data,
		options,
	}
}

func (y *XMLConverter) FromNix() (string, error) {
	return FromNix(y.data, y.options)
}

func (y *XMLConverter) ToNix() (string, error) {
	return ToNix(y.data, y.options)
}

func (y *XMLConverter) Type() string {
	return "xml"
}
