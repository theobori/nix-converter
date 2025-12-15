package yaml

func IsCharUnsafe(c byte) bool {
	return c == '#' || c == '&' || c == '*' || c == '>' || c == '!' || c == ','
}

func IsStringUnsafe(s string) bool {
	n := len(s)

	if n == 0 || IsCharUnsafe(s[0]) {
		return true
	}

	for i := 1; i < n; i++ {
		if s[i] == '#' && s[i-1] == ' ' {
			return true
		}
	}

	return false
}
