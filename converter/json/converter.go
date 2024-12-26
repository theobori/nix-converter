package json

type JSONConverter struct {
	data string
}

func NewJSONConverter(data string) *JSONConverter {
	return &JSONConverter{
		data,
	}
}

func (n *JSONConverter) FromNix() (string, error) {
	return FromNix(n.data)
}

func (n *JSONConverter) ToNix() (string, error) {
	return ToNix(n.data)
}

func (y *JSONConverter) Type() string {
	return "json"
}
