package converter

type Converter interface {
	FromNix() (string, error)
	ToNix() (string, error)
	Type() string
}
