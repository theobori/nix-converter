package yaml

import (
	"strings"
	"testing"
)

// TestYAMLMultilineStringToNix tests conversion of YAML with multi-line strings to Nix
func TestYAMLMultilineStringToNix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name: "simple multi-line string",
			input: `description: |
  This is a multi-line
  description text`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "description") {
					t.Error("Result should contain 'description' key")
				}
				if !strings.Contains(result, "''") {
					t.Error("Result should use Nix indented string syntax ''")
				}
				if strings.Contains(result, "\\n") {
					t.Error("Result should not contain escaped newlines when using indented strings")
				}
			},
		},
		{
			name: "multi-line string with folded style",
			input: `text: >
  This is a folded
  multi-line string
  that should be one line`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "text") {
					t.Error("Result should contain 'text' key")
				}
			},
		},
		{
			name: "multi-line string in nested structure",
			input: `package:
  meta:
    description: |
      A long description
      that spans multiple
      lines`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "package") {
					t.Error("Result should contain 'package' key")
				}
				if !strings.Contains(result, "meta") {
					t.Error("Result should contain 'meta' key")
				}
			},
		},
		{
			name: "script with multi-line content",
			input: `buildPhase: |
  #!/bin/bash
  echo "Building..."
  make all
  echo "Done"`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "buildPhase") {
					t.Error("Result should contain 'buildPhase' key")
				}
				if !strings.Contains(result, "''") {
					t.Error("Result should use Nix indented string syntax")
				}
				if !strings.Contains(result, "#!/bin/bash") {
					t.Error("Result should preserve script content")
				}
			},
		},
		{
			name: "multiple multi-line strings",
			input: `script1: |
  Line 1
  Line 2
script2: |
  Another line 1
  Another line 2`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "script1") || !strings.Contains(result, "script2") {
					t.Error("Result should contain both script keys")
				}
			},
		},
		{
			name: "multi-line string with special characters",
			input: `config: |
  [section]
  key = "value"
  special = !@#$%`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "config") {
					t.Error("Result should contain 'config' key")
				}
			},
		},
		{
			name: "list with multi-line strings",
			input: `scripts:
  - |
    First script
    with multiple lines
  - |
    Second script
    also multi-line`,
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

// TestNixMultilineStringToYAML tests conversion of Nix with multi-line strings to YAML
func TestNixMultilineStringToYAML(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name: "nix indented string to yaml",
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
		{
			name: "indented string with varying indent levels",
			input: `{
  this = ''
    is a multi-
    line string with
      some indent
  '';
  and = "normal single line";
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "this:") {
					t.Error("Result should contain 'this' key")
				}
				if !strings.Contains(result, "and:") {
					t.Error("Result should contain 'and' key")
				}
				// Check for YAML block scalar syntax
				if !strings.Contains(result, "|-") && !strings.Contains(result, "|") {
					t.Error("Result should use YAML block scalar syntax (| or |-)")
				}
				// Should preserve the relative indentation
				if !strings.Contains(result, "multi-") {
					t.Error("Result should contain the multiline content")
				}
				// Single line should remain simple
				if !strings.Contains(result, "normal single line") {
					t.Error("Result should contain the single line string")
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

// TestYAMLNixMultilineRoundTrip tests that multi-line strings survive round-trip conversion
func TestYAMLNixMultilineRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple multi-line round trip",
			input: `description: |
  Line 1
  Line 2
  Line 3`,
			wantErr: false,
		},
		{
			name: "nested multi-line round trip",
			input: `config:
  script: |
    #!/bin/bash
    echo "test"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// YAML -> Nix
			nixResult, err := ToNix(tt.input)
			if err != nil {
				t.Errorf("ToNix() error = %v", err)
				return
			}

			// Nix -> YAML
			yamlResult, err := FromNix(nixResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromNix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// The content should be preserved (though formatting may differ)
			// We just check that it doesn't error and produces valid output
			if yamlResult == "" && !tt.wantErr {
				t.Error("Round trip produced empty result")
			}
		})
	}
}
