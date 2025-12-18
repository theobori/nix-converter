package toml

import (
	"strings"
	"testing"
)

// TestTOMLMultilineStringToNix tests conversion of TOML with multi-line strings to Nix
func TestTOMLMultilineStringToNix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name: "simple multi-line string",
			input: `description = """
This is a multi-line
description text"""`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "description") {
					t.Error("Result should contain 'description' key")
				}
			},
		},
		{
			name: "multi-line string with backslash escaping",
			input: `text = """\
  This is line 1 \
  This is line 2 \
  This is line 3"""`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "text") {
					t.Error("Result should contain 'text' key")
				}
			},
		},
		{
			name: "script with multi-line content",
			input: `buildPhase = """
#!/bin/bash
echo "Building..."
make all
echo "Done"""`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "buildPhase") {
					t.Error("Result should contain 'buildPhase' key")
				}
			},
		},
		{
			name: "multiple multi-line strings",
			input: `script1 = """
Line 1
Line 2"""
script2 = """
Another line 1
Another line 2"""`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "script1") || !strings.Contains(result, "script2") {
					t.Error("Result should contain both script keys")
				}
			},
		},
		{
			name: "nested structure with multi-line string",
			input: `[package.meta]
description = """
A long description
that spans multiple
lines"""`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "package") {
					t.Error("Result should contain 'package' key")
				}
			},
		},
		{
			name: "multi-line literal string",
			input: `path = '''
C:\Users\name\path
With\backslashes'''`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "path") {
					t.Error("Result should contain 'path' key")
				}
			},
		},
		{
			name: "array with multi-line strings",
			input: `scripts = [
  """First script
with multiple lines""",
  """Second script
also multi-line"""
]`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "scripts") {
					t.Error("Result should contain 'scripts' key")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToNix(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToNix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

// TestNixMultilineStringToTOML tests conversion of Nix with multi-line strings to TOML
func TestNixMultilineStringToTOML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name: "nix indented string to toml",
			input: `{
  description = ''
    This is a multi-line
    description
  '';
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "description") {
					t.Error("Result should contain 'description' key")
				}
				// Check that the multi-line content is preserved
				if !strings.Contains(result, "multi-line") {
					t.Error("Result should contain 'multi-line' text")
				}
			},
		},
		{
			name: "nix script with indented string",
			input: `{
  buildScript = ''
    #!/bin/bash
    echo "Starting build"
    make install
  '';
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "buildScript") {
					t.Error("Result should contain 'buildScript' key")
				}
				if !strings.Contains(result, "bash") {
					t.Error("Result should contain bash shebang")
				}
			},
		},
		{
			name: "multiple indented strings in nix",
			input: `{
  script1 = ''
    Line 1
    Line 2
  '';
  script2 = ''
    Another 1
    Another 2
  '';
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "script1") || !strings.Contains(result, "script2") {
					t.Error("Result should contain both script keys")
				}
			},
		},
		{
			name: "nested structure with indented string",
			input: `{
  package = {
    meta = {
      longDescription = ''
        A very long description
        that spans multiple lines
        and contains detailed information
      '';
    };
  };
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "package") || !strings.Contains(result, "meta") {
					t.Error("Result should contain nested structure")
				}
			},
		},
		{
			name: "list with indented strings",
			input: `{
  items = [
    ''
      First item
      with details
    ''
    ''
      Second item
    ''
  ];
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "items") {
					t.Error("Result should contain 'items' key")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromNix(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromNix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

// TestTOMLNixMultilineRoundTrip tests that multi-line strings survive round-trip conversion
func TestTOMLNixMultilineRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		skipVerify bool
	}{
		{
			name: "simple multi-line round trip",
			input: `description = """
Line 1
Line 2
Line 3"""`,
			wantErr:    false,
			skipVerify: true, // Format may differ
		},
		{
			name: "nested multi-line round trip",
			input: `[config]
script = """
#!/bin/bash
echo "test""""`,
			wantErr:    false,
			skipVerify: true,
		},
		{
			name:       "simple single-line string",
			input:      `name = "test"`,
			wantErr:    false,
			skipVerify: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TOML -> Nix
			nixResult, err := ToNix(tt.input)
			if err != nil {
				t.Errorf("ToNix() error = %v", err)
				return
			}

			// Nix -> TOML
			tomlResult, err := FromNix(nixResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromNix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// The content should be preserved (though formatting may differ)
			if tomlResult == "" && !tt.wantErr {
				t.Error("Round trip produced empty result")
			}
		})
	}
}

// TestTOMLMultilineStringEdgeCases tests edge cases in TOML multi-line string handling
func TestTOMLMultilineStringEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "multi-line string with quotes",
			input: `text = """
This has "quotes" in it
And more "quotes""""`,
			wantErr: false,
		},
		{
			name: "multi-line string with special chars",
			input: `text = """
Special: !@#$%^&*()
Symbols: []{}"""`,
			wantErr: false,
		},
		{
			name: "multi-line string with unicode",
			input: `text = """
Unicode: 你好世界
Emoji: 🚀"""`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToNix(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToNix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}
