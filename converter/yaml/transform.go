package yaml

import "github.com/theobori/nix-converter/internal/common"

func MakeNameSafe(s string, forceUnsafe bool) string {
	if forceUnsafe {
		// Check some YAML edges case
		if IsStringUnsafe(s) {
			return common.MakeStringSafe(s)
		}

		return s
	}

	return common.MakeStringSafe(s)
}
