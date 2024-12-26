package common

import (
	"testing"
)

type ConvertFn func(string) (string, error)

func TestFromNix(t *testing.T, s string, fromNix ConvertFn, toNix ConvertFn) {
	// Convert to JSON
	nixString, err := fromNix(s)
	if err != nil {
		t.Fatal(err)
	}

	// Convert back to Nix
	jsonString, err := toNix(nixString)
	if err != nil {
		t.Fatal(err)
	}

	if jsonString != s {
		t.Fatal("not matching the original Nix string")
	}
}

func TestFromNixStrings(t *testing.T, nixStrings []string, fromNix ConvertFn, toNix ConvertFn) {
	for _, s := range nixStrings {
		TestFromNix(t, s, fromNix, toNix)
	}
}

func TestToNix(t *testing.T, s string, fromNix ConvertFn, toNix ConvertFn) {
	// Convert to Nix
	nixString, err := toNix(s)
	if err != nil {
		t.Fatal(err)
	}

	// Convert back to JSON
	jsonString, err := fromNix(nixString)
	if err != nil {
		t.Fatal(err)
	}

	if jsonString != s {
		t.Fatal("not matching the original JSON string")
	}
}

func TestToNixStrings(t *testing.T, jsonStrings []string, fromNix ConvertFn, toNix ConvertFn) {
	for _, s := range jsonStrings {
		TestToNix(t, s, fromNix, toNix)
	}
}
