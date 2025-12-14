package common

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
)

type ConvertFn func(string, *converter.ConverterOptions) (string, error)

func TestHelperFromNix(
	t *testing.T,
	s string,
	fromNix ConvertFn,
	toNix ConvertFn,
	options *converter.ConverterOptions,
) {
	// Convert to data
	nixString, err := fromNix(s, options)
	if err != nil {
		t.Fatal(err)
	}

	// Convert back to Nix
	dataString, err := toNix(nixString, options)
	if err != nil {
		t.Fatal(err)
	}

	if dataString != s {
		t.Fatal("not matching the original Nix string")
	}
}

func TestHelperFromNixStrings(
	t *testing.T,
	nixStrings []string,
	fromNix ConvertFn,
	toNix ConvertFn,
	options *converter.ConverterOptions,
) {
	for _, s := range nixStrings {
		TestHelperFromNix(t, s, fromNix, toNix, options)
	}
}

func TestHelperToNix(
	t *testing.T,
	s string,
	fromNix ConvertFn,
	toNix ConvertFn,
	options *converter.ConverterOptions,
) {
	// Convert to Nix
	nixString, err := toNix(s, options)
	if err != nil {
		t.Fatal(err)
	}

	// Convert back to data
	dataString, err := fromNix(nixString, options)
	if err != nil {
		t.Fatal(err)
	}

	if dataString != s {
		t.Fatal("not matching the original data string")
	}
}

func TestHelperToNixStrings(
	t *testing.T,
	dataStrings []string,
	fromNix ConvertFn,
	toNix ConvertFn,
	options *converter.ConverterOptions,
) {
	for _, s := range dataStrings {
		TestHelperToNix(t, s, fromNix, toNix, options)
	}
}
