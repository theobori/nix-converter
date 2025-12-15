package yaml

func MakeNameSafe(s string) string {
	return MakeStringSafe(s)
}

func MakeStringSafe(s string) string {
	return "\"" + s + "\""
}
