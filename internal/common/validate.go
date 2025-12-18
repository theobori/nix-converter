package common

import "strings"

func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || IsNumeric(c)
}

// EscapeNixString escapes special characters in a string for Nix double-quoted strings
func EscapeNixString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "$", "\\$")
	return s
}

// FormatNixIndentedString formats a string using Nix indented string syntax ''...''
func FormatNixIndentedString(s string, baseIndent string) string {
	// Handle empty strings
	if s == "" {
		return "''\n" + baseIndent + "''"
	}

	lines := strings.Split(s, "\n")
	
	// Remove trailing empty line if present (common in YAML block scalars)
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	if len(lines) == 0 {
		return "''\n" + baseIndent + "''"
	}
	
	// Build indented string content
	result := "''\n"
	for _, line := range lines {
		// Escape special sequences in indented strings
		// ${ needs to become ''${
		line = strings.ReplaceAll(line, "${", "''${")
		// '' needs to become '''
		line = strings.ReplaceAll(line, "''", "'''")
		
		result += baseIndent + "  " + line + "\n"
	}
	result += baseIndent + "''"
	
	return result
}

// ShouldUseIndentedString determines if a string should use Nix indented string syntax
func ShouldUseIndentedString(s string) bool {
	return strings.Contains(s, "\n")
}
