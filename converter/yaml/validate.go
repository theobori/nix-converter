package yaml

import "github.com/theobori/nix-converter/internal/common"

func IsYAMLString(s string) bool {
	for i := range s {
		if !common.IsAlphaNumeric(s[i]) && s[i] != ' ' {
			return false
		}
	}

	return true
}
