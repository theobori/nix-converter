package toml

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/nix"
)

func FromNix(data string, options *converter.ConverterOptions) (string, error) {
	if options.SortIterators.SortList || options.SortIterators.SortHashmap {
		return "", fmt.Errorf("sorting options are not supported from 'Nix' to 'TOML'")
	}

	// Get a Go value
	v, err := nix.GoValue(data)
	if err != nil {
		return "", err
	}

	// Get a TOML representation of the Go value
	tomlBytes, err := toml.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(tomlBytes), nil
}
