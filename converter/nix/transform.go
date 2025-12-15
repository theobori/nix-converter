package nix

import "github.com/theobori/nix-converter/internal/common"

func MakeNameSafe(s string, forceUnsafe bool) string {
	if forceUnsafe {
		// Check some Nix edges case
		if IsNameUnsafe(s) {
			return common.MakeStringSafe(s)
		}

		return s
	}

	return common.MakeStringSafe(s)
}

func MakeElementSafe(s string) string {
	if IsElementUnsafe(s) {
		return "(" + s + ")"
	}

	return s
}
