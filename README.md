# Nix data format converter

[![build](https://github.com/theobori/nix-converter/actions/workflows/build.yml/badge.svg)](https://github.com/theobori/nix-converter/actions/workflows/build.yml)

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

This GitHub repository is a toy project in the form of a CLI tool. It allows you to convert a data format, such as JSON, into Nix, and vice versa.

The project is based on various projects that provide a parser. Here are the references.
- [go-nix](https://github.com/orivej/go-nix) for the Nix language.
- [fastjson](https://github.com/valyala/fastjson) for the JSON language.
- [yaml.v3](https://gopkg.in/yaml.v3) for the YAML language.
- [go-toml](https://github.com/pelletier/go-toml) for the TOML language.

AST traversal for the Nix language remains static; Nix expressions are not evaluated.

## Supported languages

The following languages are supported.

| Language | To Nix | From Nix |
| - | - | - |
| JSON | Yes | Yes |
| YAML | Yes | Yes |
| TOML | Yes | No |

## Getting started

To start using the tool, simply run the following command.

```bash
nix-converter --help
```

By default, the program reads the standard input.

## Examples

Here are a few examples of how to use the tool.

### From Nix to JSON using the standard input.
```bash
echo -n "{a = [1 2 3];}" | nix-converter --from-nix -l json
```

### From YAML to Nix using a file named `a.yaml`.
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

### From Nix to YAML using a file named `a.nix`.
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

## Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).
