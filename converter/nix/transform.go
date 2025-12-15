package nix

func SafeName(s string) string {
	if IsNameUnsafe(s) {
		return "\"" + s + "\""
	}

	return s
}

func SafeElement(s string) string {
	if IsElementUnsafe(s) {
		return "(" + s + ")"
	}

	return s
}
