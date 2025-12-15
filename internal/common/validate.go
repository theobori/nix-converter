package common

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsCharAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func IsCharAlphaNumeric(c byte) bool {
	return IsCharAlpha(c) || IsNumeric(c)
}

func IsNumber(s string) bool {
	n := len(s)

	if n == 0 {
		return false
	}

	minus := false
	i := 0

	// The first character must be a digit or a minus
	if IsNumeric(s[i]) {
		i++
	} else if s[i] == '-' {
		i++
		minus = true
	} else {
		return false
	}

	// The second character cant be a dot if the first is a minus
	if i < n && minus && s[i] == '.' {
		return false
	}

	dot := false
	// Every remaining characters should be digits or a dot
	// There must be only one dot
	for ; i < n; i++ {
		if IsNumeric(s[i]) {
			continue
		}

		if s[i] == '.' && i < n-1 {
			if dot {
				return false
			}

			dot = true
			continue
		}

		return false
	}

	return true
}
