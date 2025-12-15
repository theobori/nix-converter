package yaml

func MakeNameSafe(s string, forceSafe bool) string {
	if !forceSafe {
		// Check some YAML edges case
		if IsStringUnsafe(s) {
			return MakeStringSafe(s)
		}

		return s
	}

	return MakeStringSafe(s)
}

func MakeStringSafe(s string) string {
	return "\"" + s + "\""
}
