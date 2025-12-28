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
	dataString, err := fromNix(s, options)
	if err != nil {
		t.Fatal(err)
	}

	// Convert back to Nix
	nixString, err := toNix(dataString, options)
	if err != nil {
		t.Fatal(err)
	}

	// For exact match (preserves formatting)
	if nixString != s {
		t.Fatalf("not matching after round-trip:\nOriginal Nix: %s\nData: %s\nBack to Nix: %s",
			s, dataString, nixString)
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
	nixString, err := toNix(s, options)
	if err != nil {
		t.Fatal(err)
	}

	dataString, err := fromNix(nixString, options)
	if err != nil {
		t.Fatal(err)
	}

	if dataString != s {
		t.Fatalf("not matching after round-trip:\nOriginal data: %s\nNix: %s\nBack to data: %s",
			s, nixString, dataString)
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
