package json

type JSONConverter struct {
	data string
}

func NewJSONConverter(data string) *JSONConverter {
	return &JSONConverter{
		data,
	}
}

func (j *JSONConverter) FromNix() (string, error) {
	return FromNix(j.data)
}

func (j *JSONConverter) ToNix() (string, error) {
	return ToNix(j.data)
}

func (j *JSONConverter) Type() string {
	return "json"
}
