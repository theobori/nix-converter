package common

import "strings"

func MakeStringSafe(s string) string {
	// Use indented string syntax for multiline strings
	if strings.Contains(s, "\n") {
		return MakeIndentedString(s)
	}

	// Escape special characters for Nix string literals
	escaped := s
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\") // Backslash must be first
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"") // Double quotes
	escaped = strings.ReplaceAll(escaped, "\r", "\\r")  // Carriage returns
	escaped = strings.ReplaceAll(escaped, "\t", "\\t")  // Tabs
	escaped = strings.ReplaceAll(escaped, "${", "\\${") // Nix interpolation
	return "\"" + escaped + "\""
}

func MakeIndentedString(s string) string {
	// Escape special Nix indented string sequences
	escaped := s
	// Escape '' as '''
	escaped = strings.ReplaceAll(escaped, "''", "'''")
	// Escape ${ as ''${
	escaped = strings.ReplaceAll(escaped, "${", "''${")

	// Split into lines and add indentation
	lines := strings.Split(escaped, "\n")
	var result strings.Builder
	result.WriteString("''\n")

	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}

		if line != "" || len(lines) > 1 {
			result.WriteString("  ")
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	result.WriteString("''")
	return result.String()
}
