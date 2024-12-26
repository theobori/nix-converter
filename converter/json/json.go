package json

import (
	"strings"

	"github.com/valyala/fastjson"
)

type NixJson struct {
	data string
}

func NewNixJson(data string) *NixJson {
	return &NixJson{
		data,
	}
}

func (nj *NixJson) FromNix() (string, error) {
	return "", nil
}

func (nj *NixJson) toNix(v *fastjson.Value, out string, indent string) string {
	switch v.Type() {
	case fastjson.TypeObject:
		o, _ := v.Object()

		e := []string{}
		o.Visit(func(key []byte, v *fastjson.Value) {
			e = append(e, indent+string(key)+" = "+nj.toNix(v, out, indent+"  ")+";")
		})

		return "{\n" + strings.Join(e, "\n") + "\n" + indent[:len(indent)-2] + "}"
	case fastjson.TypeArray:
		arr, _ := v.Array()

		e := []string{}
		for _, el := range arr {
			e = append(e, indent+nj.toNix(el, out, indent+"  "))
		}

		return "[\n" + strings.Join(e, "\n") + "\n" + indent[:len(indent)-2] + "]"
	case fastjson.TypeString:
		return v.String()
	case fastjson.TypeNumber:
		return v.String()
	case fastjson.TypeFalse:
		return "false"
	case fastjson.TypeTrue:
		return "true"
	case fastjson.TypeNull:
		return "null"
	}

	return ""
}

func (nj *NixJson) ToNix() (string, error) {
	v, err := fastjson.Parse(nj.data)
	if err != nil {
		return "", err
	}

	out := nj.toNix(v, "", "  ")

	return out, nil
}
