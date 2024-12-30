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
	"github.com/theobori/nix-converter/converter/toml"
	"github.com/theobori/nix-converter/converter/yaml"
)

func ConverterFromLanguage(language string, data string) (*converter.Converter, error) {
	var c converter.Converter

	switch language {
	case "json":
		c = json.NewJSONConverter(data)
	case "yaml":
		c = yaml.NewYAMLConverter(data)
	case "toml":
		c = toml.NewTOMLConverter(data)
	default:
		return nil, fmt.Errorf("this configuration language is not implemented")
	}

	return &c, nil
}

func main() {
	var (
		err      error
		language string
		filename string
		fromNix  bool
	)

	flag.StringVar(&language, "language", "json", "Configuration language name")
	flag.StringVar(&language, "l", "json", "Configuration language name (shorthand)")
	language = strings.ToLower(language)

	flag.StringVar(&filename, "filename", "", "Read input from a file")
	flag.StringVar(&filename, "f", "", "Read input from a file (shorthand)")

	flag.BoolVar(&fromNix, "from-nix", false, "Convert Nix to a data format, instead of data format to Nix")

	flag.Parse()

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

	c, err := ConverterFromLanguage(language, data)
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
