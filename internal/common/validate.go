package common

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func IsAlphaNumeric(c byte) bool {
	return IsAlpha(c) || IsNumeric(c)
}
