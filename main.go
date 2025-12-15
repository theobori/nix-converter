package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/theobori/nix-converter/converter"
	"github.com/theobori/nix-converter/converter/json"
	"github.com/theobori/nix-converter/converter/options"
	"github.com/theobori/nix-converter/converter/toml"
	"github.com/theobori/nix-converter/converter/yaml"
)

func ConverterFromLanguage(language string, data string, options *converter.ConverterOptions) (*converter.Converter, error) {
	var c converter.Converter

	switch language {
	case "json":
		c = json.NewJSONConverter(data, options)
	case "yaml":
		c = yaml.NewYAMLConverter(data, options)
	case "toml":
		c = toml.NewTOMLConverter(data, options)
	default:
		return nil, fmt.Errorf("this configuration language is not implemented")
	}

	return &c, nil
}

func main() {
	var (
		err               error
		language          string
		filename          string
		fromNix           bool
		sortIteratorsLine string
		unsafeKeys        bool
	)

	flag.StringVar(&language, "language", "json", "Configuration language name")
	flag.StringVar(&language, "l", "json", "Configuration language name (shorthand)")
	language = strings.ToLower(language)

	flag.StringVar(&filename, "filename", "", "Read input from a file")
	flag.StringVar(&filename, "f", "", "Read input from a file (shorthand)")

	flag.BoolVar(&fromNix, "from-nix", false, "Convert Nix to a data format, instead of data format to Nix")
	flag.StringVar(&sortIteratorsLine, "sort-iterators", "", "If possible, it sorts iterators, specify them separated by ',' like 'list,hashmap'")
	flag.BoolVar(&unsafeKeys, "unsafe-keys", false, "If possible, it skips double quotes around hashmaps keys")

	flag.Parse()

	var sortIterators *options.SortIterators
	if sortIteratorsLine != "" {
		sortIterators, err = options.NewSortIteratorsFromLine(sortIteratorsLine)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		sortIterators = options.NewDefaultSortIterators()
	}

	converterOptions := converter.ConverterOptions{
		SortIterators: *sortIterators,
		UnsafeKeys:    unsafeKeys,
	}

	var bytes []byte
	if filename == "" {
		bytes, err = io.ReadAll(os.Stdin)
	} else {
		bytes, err = os.ReadFile(filename)
	}

	if err != nil {
		log.Fatalln(err)
	}

	data := string(bytes)

	c, err := ConverterFromLanguage(language, data, &converterOptions)
	if err != nil {
		log.Fatalln(err)
	}

	var s string
	if fromNix {
		s, err = (*c).FromNix()
	} else {
		s, err = (*c).ToNix()
	}

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(s)
}
