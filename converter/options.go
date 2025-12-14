package converter

import (
	"github.com/theobori/nix-converter/converter/options"
)

type ConverterOptions struct {
	SortIterators options.SortIterators
}

func NewDefaultConverterOptions() *ConverterOptions {
	return &ConverterOptions{
		SortIterators: *options.NewDefaultSortIterators(),
	}
}
