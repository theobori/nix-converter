package yaml

import (
	"strings"
	"testing"
)

// TestMultilineToIndentedString tests that multiline YAML strings convert to Nix indented strings
func TestMultilineToIndentedString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checks  []func(t *testing.T, result string)
	}{
		{
			name: "multiline string uses indented syntax",
			input: `message: |
  Hello world
  This is a test
  Multiple lines`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use Nix indented string syntax ''")
					}
				},
				func(t *testing.T, result string) {
					if strings.Contains(result, "\\n") {
						t.Error("Should not contain escaped newlines")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "Hello world") {
						t.Error("Should preserve content")
					}
				},
			},
		},
		{
			name: "nested multiline string",
			input: `config:
  notes: |
    Line one
    Line two
    Line three`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if strings.Contains(result, "\\n") {
						t.Error("Should not escape newlines")
					}
				},
			},
		},
		{
			name: "single line string uses quotes",
			input: `name: "just one line"`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if strings.Contains(result, "''") {
						t.Error("Single line should not use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "\"just one line\"") {
						t.Error("Should use regular quoted string")
					}
				},
			},
		},
		{
			name: "multiline with dollar signs",
			input: `script: |
  echo $HOME
  export PATH=$PATH:/usr/bin`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "$HOME") {
						t.Error("Should preserve $HOME")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "$PATH") {
						t.Error("Should preserve $PATH")
					}
				},
			},
		},
		{
			name: "multiline with interpolation syntax",
			input: `script: |
  echo ${HOME}
  value=${VAR}`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''${") {
						t.Error("Should escape ${ as ''${")
					}
				},
			},
		},
		{
			name: "multiline with empty lines",
			input: `text: |
  First line

  Third line after blank`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "First line") {
						t.Error("Should preserve first line")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "Third line") {
						t.Error("Should preserve line after blank")
					}
				},
			},
		},
		{
			name: "multiline with special nix chars",
			input: `code: |
  let x = ''hello'';
  echo "test"`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "'''") {
						t.Error("Should escape '' as '''")
					}
				},
			},
		},
		{
			name: "folded string becomes single line",
			input: `description: >
  This is folded
  into one line`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if strings.Contains(result, "''") {
						t.Error("Folded strings collapse to single line, should not use indented syntax")
					}
				},
			},
		},
		{
			name: "multiline in array",
			input: `items:
  - |
    First item
    with details
  - |
    Second item`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax for array items")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "First item") {
						t.Error("Should preserve array item content")
					}
				},
			},
		},
		{
			name: "multiline with indentation variations",
			input: `code: |
  if true; then
    echo "nested"
      echo "more nested"
  fi`,
			wantErr: false,
			checks: []func(t *testing.T, result string){
				func(t *testing.T, result string) {
					if !strings.Contains(result, "''") {
						t.Error("Should use indented string syntax")
					}
				},
				func(t *testing.T, result string) {
					if !strings.Contains(result, "if true") {
						t.Error("Should preserve code structure")
					}
				},
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
			
			for _, check := range tt.checks {
				check(t, result)
			}
		})
	}
}
