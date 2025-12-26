package common

import "strings"

func MakeStringSafe(s string) string {
	// Use indented string syntax for multiline strings
	if strings.Contains(s, "\n") {
		return MakeIndentedString(s, "")
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

func MakeIndentedString(s string, indent string) string {
	// Escape special Nix indented string sequences
	escaped := s
	// Escape '' as '''
	escaped = strings.ReplaceAll(escaped, "''", "'''")
	// Escape ${ as ''${
	escaped = strings.ReplaceAll(escaped, "${", "''${")

	hasTrailingNewline := strings.HasSuffix(escaped, "\n")
	trimmed := strings.TrimSuffix(escaped, "\n")
	lines := strings.Split(trimmed, "\n")

	var result strings.Builder
	result.WriteString("''\n")

	contentIndent := indent + "  "

	for i, line := range lines {
		if line != "" || len(lines) > 1 {
			result.WriteString(contentIndent)
			result.WriteString(line)
		}

		if i < len(lines)-1 || hasTrailingNewline {
			result.WriteString("\n")
		}
	}

	if hasTrailingNewline {
		result.WriteString(indent + "''")
	} else {
		result.WriteString("''")
	}
	return result.String()
}
