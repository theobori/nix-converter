package toml

import (
	"github.com/BurntSushi/toml"
	"github.com/theobori/nix-converter/converter/nix"
)

func FromNix(data string) (string, error) {
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
