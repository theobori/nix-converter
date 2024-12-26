package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/theobori/data2nix/converter"
	"github.com/theobori/data2nix/converter/json"
)

func ConverterFromMode(mode string, data string) (*converter.Converter, error) {
	var c converter.Converter

	switch mode {
	case "json":
		c = json.NewNixJson(data)
	default:
		return nil, fmt.Errorf("this configuration language is not implemented")
	}

	return &c, nil
}

func main() {
	var (
		err      error
		mode     string
		filename string
	)

	flag.StringVar(&mode, "mode", "json", "Configuration language name")
	flag.StringVar(&mode, "m", "json", "Configuration language name")
	mode = strings.ToLower(mode)

	flag.StringVar(&filename, "filename", "", "Read input from a file")
	flag.StringVar(&filename, "f", "", "Read input from a file")

	flag.Parse()

	var bytes []byte
	if filename == "" {
		bytes, err = io.ReadAll(os.Stdin)
	} else {
		bytes, err = os.ReadFile(filename)
	}

	if err != nil {
		panic(err)
	}

	data := string(bytes)

	c, err := ConverterFromMode(mode, data)
	if err != nil {
		log.Fatalln(err)
	}

	s, err := (*c).ToNix()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(s)
}
