package yaml

type YAMLConverter struct {
	data string
}

func NewYAMLConverter(data string) *YAMLConverter {
	return &YAMLConverter{
		data,
	}
}

func (y *YAMLConverter) FromNix() (string, error) {
	return FromNix(y.data)
}

func (y *YAMLConverter) ToNix() (string, error) {
	return ToNix(y.data)
}

func (y *YAMLConverter) Type() string {
	return "yaml"
}
