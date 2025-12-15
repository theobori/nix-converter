package common

func MakeStringSafe(s string) string {
	return "\"" + s + "\""
}
