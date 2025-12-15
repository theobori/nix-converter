package nix

func MakeNameSafe(s string) string {
	if IsNameUnsafe(s) {
		return "\"" + s + "\""
	}

	return s
}

func MakeElementSafe(s string) string {
	if IsElementUnsafe(s) {
		return "(" + s + ")"
	}

	return s
}
