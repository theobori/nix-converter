package nix

import "github.com/theobori/nix-converter/internal/common"

func IsCharSafe(c byte) bool {
	return common.IsCharAlphaNumeric(c) || c == '-' || c == '_'
}

func IsNameUnsafe(s string) bool {
	n := len(s)

	if n == 0 {
		return true
	}

	if !common.IsCharAlpha(s[0]) {
		return true
	}

	for i := 1; i < n; i++ {
		if !IsCharSafe(s[i]) {
			return true
		}
	}

	return false
}

func IsElementUnsafe(s string) bool {
	return len(s) > 1 && s[0] == '-'
}
