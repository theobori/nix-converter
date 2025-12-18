package nix

import (
	"testing"

	"github.com/orivej/go-nix/nix/parser"
)

// TestMultilineStringParsing tests parsing of various multi-line string formats in Nix
func TestMultilineStringParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name: "simple indented string",
			input: `''
  Hello
  World
''`,
			expected: "Hello\nWorld\n",
			wantErr:  false,
		},
		{
			name: "indented string with common indent",
			input: `''
  Line 1
  Line 2
    Indented line
  Line 3
''`,
			expected: "Line 1\nLine 2\n  Indented line\nLine 3\n",
			wantErr:  false,
		},
		{
			name: "indented string with no indent",
			input: `''
Hello
World
''`,
			expected: "Hello\nWorld\n",
			wantErr:  false,
		},
		{
			name: "indented string with escaped content",
			input: `''
  This has ''${escaped} content
  And ''\ also
''`,
			expected: "This has ${escaped} content\nAnd ' also\n",
			wantErr:  false,
		},
		{
			name: "single line indented string",
			input: `''
  Single line
''`,
			expected: "Single line\n",
			wantErr:  false,
		},
		{
			name:     "empty indented string",
			input:    `''''`,
			expected: "",
			wantErr:  false,
		},
		{
			name: "indented string with blank lines",
			input: `''
  First line

  Second line
''`,
			expected: "First line\n\nSecond line\n",
			wantErr:  false,
		},
		{
			name: "indented string with tabs",
			input: `''
	Tab indented
	Line 2
''`,
			expected: "Tab indented\nLine 2\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := parser.ParseString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if p.Result.Type != parser.IStringNode {
				t.Errorf("Expected IStringNode, got %v", p.Result.Type)
				return
			}

			// Handle empty indented strings
			if len(p.Result.Nodes) == 0 || len(p.Result.Nodes[0].Tokens) == 0 {
				result := ""
				if result != tt.expected {
					t.Errorf("ParseString() got = %q, want %q", result, tt.expected)
				}
				return
			}

			raw := p.TokenString(p.Result.Nodes[0].Tokens[0])
			result := ProcessIndentedString(raw)
			if result != tt.expected {
				t.Errorf("ParseString() got = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestMultilineStringInSet tests multi-line strings within Nix attribute sets
func TestMultilineStringInSet(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		key     string
		want    string
		wantErr bool
	}{
		{
			name: "simple set with multiline string",
			input: `{
  description = ''
    This is a multi-line
    description for the package
  '';
}`,
			key:     "description",
			want:    "This is a multi-line\ndescription for the package\n",
			wantErr: false,
		},
		{
			name: "set with multiple multiline strings",
			input: `{
  script = ''
    #!/bin/bash
    echo "Hello"
    echo "World"
  '';
  config = ''
    [section]
    key = value
  '';
}`,
			key:     "script",
			want:    "#!/bin/bash\necho \"Hello\"\necho \"World\"\n",
			wantErr: false,
		},
		{
			name: "nested set with multiline string",
			input: `{
  package = {
    meta = {
      longDescription = ''
        This package provides functionality
        for converting between Nix and other formats.
      '';
    };
  };
}`,
			key:     "longDescription",
			want:    "This package provides functionality\nfor converting between Nix and other formats.\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GoValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Navigate to the key and check the value
			data, ok := result.(map[string]any)
			if !ok {
				t.Errorf("Expected map[string]any, got %T", result)
				return
			}

			var value string
			var found bool

			// Simple key lookup
			if v, ok := data[tt.key]; ok {
				value, found = v.(string)
			} else {
				// Try nested lookup
				for _, v := range data {
					if nested, ok := v.(map[string]any); ok {
						if findInNested(nested, tt.key, &value) {
							found = true
							break
						}
					}
				}
			}

			if !found {
				t.Errorf("Key %q not found in result", tt.key)
				return
			}

			if value != tt.want {
				t.Errorf("Value for key %q = %q, want %q", tt.key, value, tt.want)
			}
		})
	}
}

// Helper function to find a key in nested maps
func findInNested(m map[string]any, key string, result *string) bool {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			*result = s
			return true
		}
	}
	for _, v := range m {
		if nested, ok := v.(map[string]any); ok {
			if findInNested(nested, key, result) {
				return true
			}
		}
	}
	return false
}

// TestMultilineStringInList tests multi-line strings within Nix lists
func TestMultilineStringInList(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		index   int
		want    string
		wantErr bool
	}{
		{
			name: "list with multiline strings",
			input: `[
  ''
    First item
    with multiple lines
  ''
  ''
    Second item
  ''
]`,
			index:   0,
			want:    "First item\nwith multiple lines\n",
			wantErr: false,
		},
		{
			name: "mixed list with multiline and regular strings",
			input: `[
  "simple string"
  ''
    Multi-line
    string
  ''
]`,
			index:   1,
			want:    "Multi-line\nstring\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GoValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			list, ok := result.([]any)
			if !ok {
				t.Errorf("Expected []any, got %T", result)
				return
			}

			if tt.index >= len(list) {
				t.Errorf("Index %d out of bounds for list of length %d", tt.index, len(list))
				return
			}

			value, ok := list[tt.index].(string)
			if !ok {
				t.Errorf("Expected string at index %d, got %T", tt.index, list[tt.index])
				return
			}

			if value != tt.want {
				t.Errorf("Value at index %d = %q, want %q", tt.index, value, tt.want)
			}
		})
	}
}

// TestMultilineStringEdgeCases tests edge cases in multi-line string handling
func TestMultilineStringEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "multiline string with special characters",
			input: `{
  text = ''
    Special chars: !@#$%^&*()
    Quotes: " ' ` + "`" + `
  '';
}`,
			wantErr: false,
		},
		{
			name: "multiline string with unicode",
			input: `{
  text = ''
    Unicode: 你好世界
    Emoji: 🚀 🎉
  '';
}`,
			wantErr: false,
		},
		{
			name: "multiline string with escape sequences",
			input: `{
  text = ''
    Escaped: ''${var}
    Another: ''\ 
  '';
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GoValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
