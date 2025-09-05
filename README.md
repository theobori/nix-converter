# Nix data format converter

[![build](https://github.com/theobori/nix-converter/actions/workflows/build.yml/badge.svg)](https://github.com/theobori/nix-converter/actions/workflows/build.yml)

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

This GitHub repository is a toy project in the form of a CLI tool. It allows you to convert a data format, such as JSON, into Nix, and vice versa.

The project is based on various projects that provide a parser. Here are the references.
- [go-nix](https://github.com/orivej/go-nix) for parsing the Nix language.
- [fastjson](https://github.com/valyala/fastjson) for parsing the JSON language.
- [yaml.v3](https://gopkg.in/yaml.v3) for parsing the YAML language.
- [pelletier/go-toml](https://github.com/pelletier/go-toml) for parsing the TOML language.
- [BurntSushi/toml](https://github.com/BurntSushi/toml) providing a TOML marshaller.

AST traversal for the Nix language remains static; Nix expressions are not evaluated.

## Supported languages

The following languages are supported.

| Language | To Nix | From Nix |
| - | - | - |
| **JSON** | Yes | Yes |
| **YAML** | Yes | Yes |
| **TOML** | Yes (unstable output) | Yes |

The YAML evaluation support anchors. They are handled during the YAML to Nix conversion.

## Getting started

To start using the tool, simply run the following command.

```bash
nix-converter --help
```

By default, the program reads the standard input.

## Examples

Here are a few examples of how to use the tool.

### From Nix to JSON using the standard input
```bash
echo -n "{a = [1 2 3];}" | nix-converter --from-nix -l json
```

### From YAML to Nix using a file named `a.yaml`
```yaml
# a.yaml
- 0
- 1
- - 2
  - - 3
    - - 4
      - - 5
        - 6
        - 7
        - 8
```

```bash
nix-converter -f a.yaml -l yaml
```

It is also possible to use multiple UNIX pipe.
```bash
nix-converter -f a.yaml -l yaml | nix-converter --from-nix -l json
```

### From YAML to Nix with anchor using a file named `anchor.yaml`
```yaml
# anchor.yaml
definitions:
  steps:
    - step: &build-test
        name: Build and test
        script:
          - mvn package
        artifacts:
          - target/**
pipelines:
  branches:
    develop:
      - step: *build-test
    main:
      - step: *build-test
```

```bash
nix-converter -f anchor.yaml -l yaml
```

### From Nix to YAML using a file named `a.nix`
```nix
# a.nix
{
  id = "1c7d8e9f0";
  users = [
    {
      name = "Alice and her cats";
      age = 2.8;
      "pets" = [
        {
          type = "321a";
          name = "Luna";
          toys = [ ];
          hello = { };
        }
        {
          type = "dog";
          name = "Max";
        }
      ];
    }
    {
      name = "Bob";
      age = 34;
      pets = null;
    }
  ];
  settings = {
    theme = {
      dark = {
        primary = "#1a1a1a";
        accent = "#4287f5";
      };
      light = {
        primary = "#ffffff";
        accent = "#2196f3";
      };
    };
    notifications = true;
  };
  meta = {
    created = "2024-01-01";
    modified = {
      by = "system";
      timestamp = "2024-02-15T14:30:00Z";
    };
  };
}
```

```bash
nix-converter --from-nix -f a.nix -l yaml
```

### Get a Go value from Nix

```go
import (
	"fmt"
	"log"

	"github.com/theobori/nix-converter/converter/nix"
)

func main() {
	nixString := `
{
  meta = {
    created = "2024-01-01";
    modified = {
      by = "system";
      timestamp = "2024-02-15T14:30:00Z";
    };
  };
  id = "1c7d8e9f0";
}`
	// Here we should get a Go map, map[string]any
	v, err := nix.GoValue(nixString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%t\n", v)
}
```

## Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).
