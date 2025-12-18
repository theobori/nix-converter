package toml

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
)

// TestTOMLMultilineStringToNix tests conversion of TOML with multi-line strings to Nix
func TestTOMLMultilineStringToNix(t *testing.T) {
	options := converter.NewDefaultConverterOptions()
	options.SortIterators.SortHashmap = true

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "simple multiline",
			input: `description = """
Multi
line"""`,
			want: `{
  "description" = ''
  Multi
  line
'';
}`,
		},
		{
			name: "nested multiline",
			input: `[package.meta]
desc = """
Long
text"""`,
			want: `{
  "package" = {
    "meta" = {
      "desc" = ''
  Long
  text
'';
    };
  };
}`,
		},
		{
			name: "script",
			input: `buildPhase = """
#!/bin/bash
echo "Build"
make all"""`,
			want: `{
  "buildPhase" = ''
  #!/bin/bash
  echo "Build"
  make all
'';
}`,
		},
		{
			name: "multiple strings and arrays",
			input: `s1 = """
A
B"""
s2 = """
C
D"""
scripts = [
  """
E
F""",
  """
G
H"""
]`,
			want: `{
  "s1" = ''
  A
  B
'';
  "s2" = ''
  C
  D
'';
  "scripts" = [
    ''
  E
  F
''
    ''
  G
  H
''
  ];
}`,
		},
		{
			name:  "literal string with backslashes",
			input: `path = '''C:\Users\name'''`,
			want: `{
  "path" = "C:\\Users\\name";
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToNix(tt.input, options)
			if err != nil {
				t.Fatalf("ToNix() error = %v", err)
			}
			if result != tt.want {
				t.Errorf("ToNix() = \n%q\nwant \n%q", result, tt.want)
			}
		})
	}
}

// TestNixMultilineStringToTOML tests conversion of Nix with multi-line strings to TOML
func TestNixMultilineStringToTOML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple and nested indented strings",
			input: `{ description = ''Multi\nline''; package = { meta = { desc = ''Long\ntext''; }; }; }`,
			want: `description = "Multi\\nline"

[package]
  [package.meta]
    desc = "Long\\ntext"
`,
		},
		{
			name:  "script and multiple strings",
			input: `{ buildScript = ''#!/bin/bash\necho "Build"''; s1 = ''A\nB''; s2 = ''C\nD''; }`,
			want: `buildScript = "#!/bin/bash\\necho \"Build\""
s1 = "A\\nB"
s2 = "C\\nD"
`,
		},
		{
			name:  "list with indented strings",
			input: `{ items = [ ''First\nitem'' ''Second\nitem'' ]; }`,
			want: `items = ["First\\nitem", "Second\\nitem"]
`,
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

// TestTOMLNixMultilineRoundTrip tests that multi-line strings survive round-trip conversion
func TestTOMLNixMultilineRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "multiline round trip",
			input: `description = """Line 1\nLine 2"""`,
			want: `description = "Line 1\nLine 2\n"
`,
		},
		{
			name: "nested multiline",
			input: `[config]
script = """
#!/bin/bash
echo "test"""`,
			want: `[config]
  script = "#!/bin/bash\necho \"test\n"
`,
		},
		{
			name:  "single line",
			input: `name = "test"`,
			want: `name = "test"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nixResult, err := ToNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("ToNix() error = %v", err)
			}

			tomlResult, err := FromNix(nixResult, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("FromNix() error = %v", err)
			}

			if tomlResult == "" {
				t.Fatal("Round trip produced empty result")
			}

			if tomlResult != tt.want {
				t.Errorf("RoundTrip() = \n%q\nwant \n%q", tomlResult, tt.want)
			}
		})
	}
}

// TestTOMLMultilineStringEdgeCases tests edge cases in TOML multi-line string handling
func TestTOMLMultilineStringEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "quotes and special chars",
			input: `text = """
Has "quotes"
Special: !@#$%
Symbols: []{}"""`,
			want: `{
  "text" = ''
  Has "quotes"
  Special: !@#$%
  Symbols: []{}
'';
}`,
		},
		{
			name: "unicode and emoji",
			input: `text = """
Unicode: 你好世界
Emoji: 🚀"""`,
			want: `{
  "text" = ''
  Unicode: 你好世界
  Emoji: 🚀
'';
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("ToNix() error = %v", err)
			}
			if result == "" {
				t.Fatal("Expected non-empty result")
			}
			if result != tt.want {
				t.Errorf("ToNix() = \n%q\nwant \n%q", result, tt.want)
			}
		})
	}
}
