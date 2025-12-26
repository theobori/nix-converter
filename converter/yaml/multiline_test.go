package yaml

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
)

// TestYAMLMultilineStringToNix tests conversion of YAML with multi-line strings to Nix
func TestYAMLMultilineStringToNix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "literal and folded styles",
			input: "description: |\n  Multi\n  line\ntext: >\n  Folded\n  text",
			want: `{
  "description" = ''
    Multi
    line
  '';
  "text" = "Folded text";
}`,
		},
		{
			name:  "nested and script",
			input: "package:\n  meta:\n    desc: |\n      Long text\nbuildPhase: |\n  #!/bin/bash\n  make all",
			want: `{
  "package" = {
    "meta" = {
      "desc" = ''
        Long text
      '';
    };
  };
  "buildPhase" = ''
    #!/bin/bash
    make all'';
}`,
		},
		{
			name:  "multiple strings and lists",
			input: "s1: |\n  A\n  B\ns2: |\n  C\nscripts:\n  - |\n    First\n  - |\n    Second",
			want: `{
  "s1" = ''
    A
    B
  '';
  "s2" = ''
    C
  '';
  "scripts" = [
    ''
      First
    ''
    "Second"
  ];
}`,
		},
		{
			name:  "special characters",
			input: "config: |\n  [section]\n  key = \"value\"\n  special = !@#$%",
			want: `{
  "config" = ''
    [section]
    key = "value"
    special = !@#$%'';
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

// TestNixMultilineStringToYAML tests conversion of Nix with multi-line strings to YAML
func TestNixMultilineStringToYAML(t *testing.T) {
	options := converter.ConverterOptions{
		SortIterators: converter.NewDefaultConverterOptions().SortIterators,
		UnsafeKeys:    true,
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple and nested indented strings",
			input: "{ description = ''Multi\nline''; package.meta.desc = ''Long\ntext''; }",
			want:  "description: |-\n  Multi\n  line\npackage:\n  meta:\n    desc: |-\n      Long\n      text",
		},
		{
			name:  "script and multiple strings",
			input: "{ buildScript = ''#!/bin/bash\necho \"Build\"''; s1 = ''A\nB''; s2 = ''C\nD''; }",
			want:  "buildScript: |-\n  #!/bin/bash\n  echo \"Build\"\ns1: |-\n  A\n  B\ns2: |-\n  C\n  D",
		},
		{
			name:  "list with indented strings",
			input: "{ items = [ ''First\nitem'' ''Second\nitem'' ]; }",
			want:  "items:\n  - |-\n      First\n      item\n  - |-\n      Second\n      item",
		},
		{
			name:  "varying indent with block scalar",
			input: "{ this = ''is a multi-\nline with\n  indent''; and = \"single line\"; }",
			want:  "this: |-\n  is a multi-\n  line with\n    indent\nand: \"single line\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FromNix(tt.input, &options)
			if err != nil {
				t.Fatalf("FromNix() error = %v", err)
			}
			if result != tt.want {
				t.Errorf("FromNix() = \n%q\nwant \n%q", result, tt.want)
			}
		})
	}
}

// TestYAMLNixMultilineRoundTrip tests that multi-line strings survive round-trip conversion
func TestYAMLNixMultilineRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "multiline round trip",
			input: "description: |\n  Line 1\n  Line 2\n  Line 3",
			want:  "\"description\": |-\n  Line 1\n  Line 2\n  Line 3",
		},
		{
			name:  "nested multiline",
			input: "config:\n  script: |\n    #!/bin/bash\n    echo \"test\"",
			want:  "\"config\":\n  \"script\": |-\n    #!/bin/bash\n    echo \"test\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nixResult, err := ToNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("ToNix() error = %v", err)
			}

			yamlResult, err := FromNix(nixResult, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Fatalf("FromNix() error = %v", err)
			}

			if yamlResult == "" {
				t.Fatal("Round trip produced empty result")
			}

			if yamlResult != tt.want {
				t.Errorf("RoundTrip() = \n%q\nwant \n%q", yamlResult, tt.want)
			}
		})
	}
}
