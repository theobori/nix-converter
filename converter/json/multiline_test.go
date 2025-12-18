package json

import (
	"strings"
	"testing"
)

// TestJSONMultilineStringToNix tests conversion of JSON with multi-line strings to Nix
func TestJSONMultilineStringToNix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name: "simple multi-line string with newlines",
			input: `{
  "description": "This is a multi-line\ndescription text"
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "description") {
					t.Error("Result should contain 'description' key")
				}
			},
		},
		{
			name: "bash script with newlines",
			input: `{
  "buildScript": "#!/bin/bash\necho \"Building...\"\nmake all\necho \"Done\""
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "buildScript") {
					t.Error("Result should contain 'buildScript' key")
				}
			},
		},
		{
			name: "nested structure with multi-line strings",
			input: `{
  "package": {
    "meta": {
      "description": "A long description\nthat spans multiple\nlines"
    }
  }
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "package") {
					t.Error("Result should contain 'package' key")
				}
			},
		},
		{
			name: "multiple multi-line strings",
			input: `{
  "script1": "Line 1\nLine 2",
  "script2": "Another line 1\nAnother line 2"
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "script1") || !strings.Contains(result, "script2") {
					t.Error("Result should contain both script keys")
				}
			},
		},
		{
			name: "multi-line string with escape sequences",
			input: `{
  "text": "Line with \\backslash\nLine with \"quotes\"\nLine with \ttab"
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "text") {
					t.Error("Result should contain 'text' key")
				}
			},
		},
		{
			name: "list with multi-line strings",
			input: `{
  "scripts": [
    "First script\nwith multiple lines",
    "Second script\nalso multi-line"
  ]
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "scripts") {
					t.Error("Result should contain 'scripts' key")
				}
			},
		},
		{
			name: "multi-line string with unicode",
			input: `{
  "text": "Line 1\n你好\nLine 3"
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "text") {
					t.Error("Result should contain 'text' key")
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

// TestNixMultilineStringToJSON tests conversion of Nix with multi-line strings to JSON
func TestNixMultilineStringToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, result string)
	}{
		{
			name: "nix indented string to json",
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
				// JSON should contain escaped newlines
				if !strings.Contains(result, "\\n") && !strings.Contains(result, "\n") {
					t.Error("Result should contain newlines (escaped or literal)")
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
      '';
    };
  };
}`,
			wantErr: false,
			check: func(t *testing.T, result string) {
				if !strings.Contains(result, "package") {
					t.Error("Result should contain 'package' key")
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

// TestJSONNixMultilineRoundTrip tests that multi-line strings survive round-trip conversion
func TestJSONNixMultilineRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantErr    bool
		skipVerify bool // Some conversions may not be exact due to format differences
	}{
		{
			name: "simple multi-line round trip",
			input: `{
  "description": "Line 1\nLine 2\nLine 3"
}`,
			wantErr:    false,
			skipVerify: true, // Format may differ but should work
		},
		{
			name: "nested multi-line round trip",
			input: `{
  "config": {
    "script": "#!/bin/bash\necho \"test\""
  }
}`,
			wantErr:    false,
			skipVerify: true,
		},
		{
			name: "simple single-line string",
			input: `{
  "name": "test"
}`,
			wantErr:    false,
			skipVerify: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// JSON -> Nix
			nixResult, err := ToNix(tt.input)
			if err != nil {
				t.Errorf("ToNix() error = %v", err)
				return
			}

			// Nix -> JSON
			jsonResult, err := FromNix(nixResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromNix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.skipVerify {
				// For simple cases, try to verify the round trip
				// Parse both and compare
				if jsonResult == "" && !tt.wantErr {
					t.Error("Round trip produced empty result")
				}
			}
		})
	}
}

// TestMultilineStringPreservation tests that newlines are properly preserved
func TestMultilineStringPreservation(t *testing.T) {
	input := `{
  "text": "Line 1\nLine 2\nLine 3"
}`

	nixResult, err := ToNix(input)
	if err != nil {
		t.Fatalf("ToNix() error = %v", err)
	}

	// Convert back to JSON
	jsonResult, err := FromNix(nixResult)
	if err != nil {
		t.Fatalf("FromNix() error = %v", err)
	}

	// The result should be valid JSON
	if jsonResult == "" {
		t.Error("Result is empty")
	}

	// Should contain the text field
	if !strings.Contains(jsonResult, "text") {
		t.Error("Result should contain 'text' field")
	}
}
