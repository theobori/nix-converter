package json

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
)

// TestJSONMultilineStringToNix tests conversion of JSON with multi-line strings to Nix
func TestJSONMultilineStringToNix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple and nested multiline",
			input: `{"description": "Multi\nline", "package": {"meta": {"desc": "Long\ntext"}}}`,
			want: `{
  "description" = "Multi\nline";
  "package" = {
    "meta" = {
      "desc" = "Long\ntext";
    };
  };
}`,
		},
		{
			name:  "bash script",
			input: `{"buildScript": "#!/bin/bash\necho \"Build\"\nmake all"}`,
			want: `{
  "buildScript" = "#!/bin/bash\necho \"Build\"\nmake all";
}`,
		},
		{
			name:  "multiple strings and lists",
			input: `{"s1": "A\nB", "s2": "C\nD", "items": ["E\nF", "G\nH"]}`,
			want: `{
  "s1" = "A\nB";
  "s2" = "C\nD";
  "items" = [
    "E\nF"
    "G\nH"
  ];
}`,
		},
		{
			name:  "escape sequences and unicode",
			input: `{"text": "Back\\slash\nQuote\"\nTab\t\n你好\n🚀"}`,
			want: `{
  "text" = "Back\\slash\nQuote\"\nTab\t\n你好\n🚀";
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("ToNix() error = %v", err)
			}
			if result != tt.want {
				t.Errorf("ToNix() = \n%q\nwant \n%q", result, tt.want)
			}
		})
	}
}

// TestNixMultilineStringToJSON tests conversion of Nix with multi-line strings to JSON
func TestNixMultilineStringToJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple and nested indented strings",
			input: `{ description = ''Multi\nline''; package.meta.desc = ''Long\ntext''; }`,
			want: `{
  "description": "Multi\\nline",
  "package": "Long\\ntext"
}`,
		},
		{
			name:  "script and multiple strings",
			input: `{ buildScript = ''#!/bin/bash\necho "Build"''; s1 = ''A\nB''; s2 = ''C\nD''; }`,
			want: `{
  "buildScript": "#!/bin/bash\\necho \"Build\"",
  "s1": "A\\nB",
  "s2": "C\\nD"
}`,
		},
		{
			name:  "list with indented strings",
			input: `{ items = [ ''First\nitem'' ''Second\nitem'' ]; }`,
			want: `{
  "items": [
    "First\\nitem",
    "Second\\nitem"
  ]
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("FromNix() error = %v", err)
			}
			if result != tt.want {
				t.Errorf("FromNix() = \n%q\nwant \n%q", result, tt.want)
			}
		})
	}
}

// TestJSONNixMultilineRoundTrip tests that multi-line strings survive round-trip conversion
func TestJSONNixMultilineRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "multiline round trip",
			input: `{"description": "Line 1\nLine 2\nLine 3"}`,
			want: `{
  "description": "Line 1\\nLine 2\\nLine 3"
}`,
		},
		{
			name:  "nested multiline",
			input: `{"config": {"script": "#!/bin/bash\necho \"test\""}}`,
			want: `{
  "config": {
    "script": "#!/bin/bash\\necho \\\"test\\\""
  }
}`,
		},
		{
			name:  "single line",
			input: `{"name": "test"}`,
			want: `{
  "name": "test"
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nixResult, err := ToNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("ToNix() error = %v", err)
			}

			jsonResult, err := FromNix(nixResult, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("FromNix() error = %v", err)
			}

			if jsonResult == "" {
				t.Fatal("Round trip produced empty result")
			}

			if jsonResult != tt.want {
				t.Errorf("RoundTrip() = \n%q\nwant \n%q", jsonResult, tt.want)
			}
		})
	}
}
