package yaml

import (
	"testing"

	"github.com/theobori/nix-converter/converter"
)

// TestMultilineToIndentedString tests that multiline YAML strings convert to Nix indented strings
func TestMultilineToIndentedString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "multiline string uses indented syntax",
			input: `message: |
  Hello world
  This is a test
  Multiple lines`,
			want: `{
  "message" = ''
    Hello world
    This is a test
    Multiple lines'';
}`,
		},
		{
			name: "nested multiline string",
			input: `config:
  notes: |
    Line one
    Line two
    Line three`,
			want: `{
  "config" = {
    "notes" = ''
      Line one
      Line two
      Line three'';
  };
}`,
		},
		{
			name:  "single line string uses quotes",
			input: `name: "just one line"`,
			want: `{
  "name" = "just one line";
}`,
		},
		{
			name: "multiline with dollar signs",
			input: `script: |
  echo $HOME
  export PATH=$PATH:/usr/bin`,
			want: `{
  "script" = ''
    echo $HOME
    export PATH=$PATH:/usr/bin'';
}`,
		},
		{
			name: "multiline with interpolation syntax",
			input: `script: |
  echo ${HOME}
  value=${VAR}`,
			want: `{
  "script" = ''
    echo ''${HOME}
    value=''${VAR}'';
}`,
		},
		{
			name: "multiline with empty lines",
			input: `text: |
  First line

  Third line after blank`,
			want: `{
  "text" = ''
    First line
    
    Third line after blank'';
}`,
		},
		{
			name: "multiline with special nix chars",
			input: `code: |
  let x = ''hello'';
  echo "test"`,
			want: `{
  "code" = ''
    let x = '''hello''';
    echo "test"'';
}`,
		},
		{
			name: "folded string becomes single line",
			input: `description: >
  This is folded
  into one line`,
			want: `{
  "description" = "This is folded into one line";
}`,
		},
		{
			name: "multiline in array",
			input: `items:
  - |
    First item
    with details
  - |
    Second item`,
			want: `{
  "items" = [
    ''
      First item
      with details
    ''
    "Second item"
  ];
}`,
		},
		{
			name: "multiline with indentation variations",
			input: `code: |
  if true; then
    echo "nested"
      echo "more nested"
  fi`,
			want: `{
  "code" = ''
    if true; then
      echo "nested"
        echo "more nested"
    fi'';
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := ToNix(tt.input, converter.NewDefaultConverterOptions())
			if err != nil {
				t.Errorf("ToNix() error = %v, input %v", err, tt.input)
				return
			}

			if result != tt.want {
				t.Errorf("ToNix() = \n%v, want \n%v", result, tt.want)
			}
		})
	}
}
