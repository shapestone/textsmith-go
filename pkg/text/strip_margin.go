package text

import (
	"regexp"
	"strings"
)

// Regex for strip margin functionality
var stripMarginGroup = regexp.MustCompile(`(?m)^[ \t]*\|(.*)(?:\r?\n|$)`)

// The StripMargin function lets you define multiline strings where each line is prepended with optional whitespace
// and a pipeline symbol
//
// Code example:
//
//	text.StripMargin(`
//	|<content line 1>
//	|<content line 2>
//	`)
func StripMargin(s string) string {
	// Handle an empty string case
	if s == "" {
		return ""
	}

	// Use Unicode-safe string operations
	lines := strings.Split(s, "\n")
	var result []string

	for _, line := range lines {
		// Check if line matches the margin pattern
		if match := stripMarginGroup.FindStringSubmatch(line + "\n"); match != nil {
			result = append(result, match[1])
		}
	}

	// If no matches found, return empty string
	if len(result) == 0 {
		return ""
	}

	return strings.Join(result, "\n")
}

// Regex for strip column functionality
var stripColumnGroup = regexp.MustCompile(`(?m)^[ \t]*\|(.*)(?:\|[ \t]*\n|\|[ \t]*$)`)

// The StripColumn function lets you define multiline strings where each line is prepended with optional whitespace
// and pipeline symbols
//
// Code example:
//
// text.StripColumn(`
//
//	|<content line 1>|
//	|<content line 2>|
//
// `)
func StripColumn(s string) string {
	ms := stripColumnGroup.FindAllStringSubmatch(s, -1)
	if ms == nil {
		return ``
	}

	lines := ``
	for idx, m := range ms {
		if idx > 0 {
			lines += "\n"
		}
		lines += m[1]
	}

	return lines
}
