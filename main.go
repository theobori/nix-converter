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
	"github.com/theobori/nix-converter/converter/yaml"
)

func ConverterFromMode(mode string, data string) (*converter.Converter, error) {
	var c converter.Converter

	switch mode {
	case "json":
		c = json.NewJSONConverter(data)
	case "yaml":
		c = yaml.NewYAMLConverter(data)
	default:
		return nil, fmt.Errorf("this data format is not implemented")
	}

	return &c, nil
}

func main() {
	var (
		err      error
		mode     string
		filename string
		fromNix  bool
	)

	flag.StringVar(&mode, "mode", "json", "Data format name")
	flag.StringVar(&mode, "m", "json", "Data format name (shorthand)")
	mode = strings.ToLower(mode)

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

	c, err := ConverterFromMode(mode, data)
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
