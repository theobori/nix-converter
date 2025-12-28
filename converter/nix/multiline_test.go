package nix

import (
	"testing"

	"github.com/orivej/go-nix/nix/parser"
)

// TestMultilineStringParsing tests parsing of various multi-line string formats in Nix
func TestMultilineStringParsing(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name: "simple and nested indent",
			input: `''
  Line 1
  Line 2
    Indented line
''`,
			expected: "Line 1\nLine 2\n  Indented line\n",
		},
		{
			name:     "no indent and empty string",
			input:    `''Hello\nWorld''`,
			expected: "Hello\\nWorld",
		},
		{
			name: "escaped content",
			input: `''
  Has ''${escaped} and ''\ quote
''`,
			expected: "Has ${escaped} and ' quote\n",
		},
		{
			name: "blank lines and tabs",
			input: `''
	Tab line

	After blank
''`,
			expected: "Tab line\n\nAfter blank\n",
		},
		{
			name:     "empty",
			input:    `''''`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p, err := parser.ParseString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil || p.Result.Type != parser.IStringNode {
				return
			}

			result := ""
			if len(p.Result.Nodes) > 0 && len(p.Result.Nodes[0].Tokens) > 0 {
				result = ProcessIndentedString(p.TokenString(p.Result.Nodes[0].Tokens[0]))
			}
			if result != tt.expected {
				t.Errorf("got = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestMultilineStringInStructures tests multi-line strings in sets and lists
func TestMultilineStringInStructures(t *testing.T) {
	t.Parallel()
	t.Run("in attribute set", func(t *testing.T) {
		t.Parallel()
		input := `{
  script = ''
    #!/bin/bash
    echo "Hello"
  '';
  config = ''
    [section]
    key = value
  '';
}`
		result, err := GoValue(input)
		if err != nil {
			t.Fatalf("GoValue() error = %v", err)
		}
		data, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("result is not a map[string]any")
		}
		if script, ok := data["script"].(string); !ok || script != "#!/bin/bash\necho \"Hello\"\n" {
			t.Errorf("script = %q", script)
		}
		if config := data["config"].(string); config != "[section]\nkey = value\n" {
			t.Errorf("config = %q", config)
		}
	})

	t.Run("in nested set", func(t *testing.T) {
		t.Parallel()
		input := `{ package = { meta = { desc = ''Multi\nline''; }; }; }`
		result, err := GoValue(input)
		if err != nil {
			t.Fatalf("GoValue() error = %v", err)
		}
		pkg := result.(map[string]any)["package"].(map[string]any)
		desc := pkg["meta"].(map[string]any)["desc"].(string)
		if desc != "Multi\\nline" {
			t.Errorf("desc = %q, want %q", desc, "Multi\\nline")
		}
	})

	t.Run("in list", func(t *testing.T) {
		t.Parallel()
		input := `[ ''First\nitem'' ''Second\nitem'' ]`
		result, err := GoValue(input)
		if err != nil {
			t.Fatalf("GoValue() error = %v", err)
		}
		list := result.([]any)
		if list[0].(string) != "First\\nitem" || list[1].(string) != "Second\\nitem" {
			t.Errorf("list = %v", list)
		}
	})
}

// TestMultilineStringEdgeCases tests edge cases in multi-line string handling
func TestMultilineStringEdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		key   string
		want  string
	}{
		{
			name:  "special characters and quotes",
			input: `{ text = ''Special: !@#$%\nQuotes: " ' ` + "`" + `''; }`,
			key:   "text",
			want:  `Special: !@#$%\nQuotes: " ' ` + "`",
		},
		{
			name:  "unicode and emoji",
			input: `{ text = ''Unicode: 你好\nEmoji: 🚀''; }`,
			key:   "text",
			want:  `Unicode: 你好\nEmoji: 🚀`,
		},
		{
			name:  "escaped sequences",
			input: `{ text = ''Escaped: ''${var}\nQuote: ''\ ''; }`,
			key:   "text",
			want:  `Escaped: ''${var}\nQuote: ''\ `,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := GoValue(tt.input)
			if err != nil {
				t.Fatalf("GoValue() error = %v", err)
			}
			text := result.(map[string]any)[tt.key].(string)
			if text != tt.want {
				t.Errorf("got = %q, want %q", text, tt.want)
			}
		})
	}
}
