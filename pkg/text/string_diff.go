package text

import (
	"fmt"
	"strings"
)

// visualizeString converts invisible characters to visible Unicode symbols
func visualizeString(s string) string {
	var builder strings.Builder

	for _, r := range s {
		switch r {
		case ' ':
			builder.WriteRune('\u2423') // ␣
		case '\t':
			builder.WriteRune('\u2409') // ␉
		case '\n':
			builder.WriteRune('\u240A') // ␊
		case '\r':
			builder.WriteRune('\u240D') // ␍
		case '\v':
			builder.WriteRune('\u240B') // ␋
		case '\f':
			builder.WriteRune('\u240C') // ␌
		default:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// formatStringForDisplay formats a string for display with optional visualization
func formatStringForDisplay(s string, visualize bool) string {
	if s == "" {
		return "<empty>¶"
	}

	displayStr := s
	if visualize {
		displayStr = visualizeString(s)
	}

	return fmt.Sprintf(`"%s"¶`, displayStr)
}

// findFirstDifference finds the position of the first differing character
// Returns the position and true if found, -1 and false if strings are equal
func findFirstDifference(s1, s2 string) (int, bool) {
	r1 := []rune(s1)
	r2 := []rune(s2)

	minLen := len(r1)
	if len(r2) < minLen {
		minLen = len(r2)
	}

	for i := 0; i < minLen; i++ {
		if r1[i] != r2[i] {
			return i, true
		}
	}

	// If one string is longer, the difference is at the end of the shorter one
	if len(r1) != len(r2) {
		return minLen, true
	}

	return -1, false
}

// formatCharacter formats a character for display with its Unicode code point
func formatCharacter(r rune) string {
	if r == '\n' {
		return fmt.Sprintf("'\\n' (U+%04X)", r)
	} else if r == '\r' {
		return fmt.Sprintf("'\\r' (U+%04X)", r)
	} else if r == '\t' {
		return fmt.Sprintf("'\\t' (U+%04X)", r)
	} else if r == ' ' {
		return fmt.Sprintf("' ' (U+%04X)", r)
	} else if r < 32 || r == 127 {
		return fmt.Sprintf("'\\x%02X' (U+%04X)", r, r)
	}
	return fmt.Sprintf("'%c' (U+%04X)", r, r)
}

// compareStringsInternal implements the comparison logic
func compareStringsInternal(actual, expected string, visualize bool) string {
	var result strings.Builder

	// Check if strings are equal
	if actual == expected {
		result.WriteString("CompareStrings: ✓ [MATCH]\n")
		result.WriteString(fmt.Sprintf("  Expected: %s\n", formatStringForDisplay(expected, visualize)))
		result.WriteString(fmt.Sprintf("  Actual:   %s", formatStringForDisplay(actual, visualize)))
		return result.String()
	}

	// Strings differ
	result.WriteString("CompareStrings: ✗ [ASSERTION_FAILED]\n")
	result.WriteString(fmt.Sprintf("- Expected: %s\n", formatStringForDisplay(expected, visualize)))
	result.WriteString(fmt.Sprintf("+ Actual:   %s\n", formatStringForDisplay(actual, visualize)))

	// Find and report the first difference
	pos, found := findFirstDifference(expected, actual)
	if found {
		result.WriteString("\n")
		result.WriteString(fmt.Sprintf("  Difference at position %d:\n", pos))

		expectedRunes := []rune(expected)
		actualRunes := []rune(actual)

		if pos < len(expectedRunes) {
			result.WriteString(fmt.Sprintf("      Expected character: %s\n", formatCharacter(expectedRunes[pos])))
		} else {
			result.WriteString("      Expected character: <end of string>\n")
		}

		if pos < len(actualRunes) {
			result.WriteString(fmt.Sprintf("      Actual character:   %s", formatCharacter(actualRunes[pos])))
		} else {
			result.WriteString("      Actual character:   <end of string>")
		}
	}

	return result.String()
}

// CompareStrings provides a test framework style comparison between actual and expected strings
// with detailed diff highlighting. It converts invisible characters to visible symbols for better debugging.
//
// Returns a formatted string showing:
// - Match status (✓ [MATCH] or ✗ [ASSERTION_FAILED])
// - Both strings with invisible characters visualized (␣ for space, ␉ for tab, etc.)
// - Position of first difference
// - Unicode code points of differing characters
//
// Example:
//
//	actual := "hello world"
//	expected := "hello mars"
//	result := CompareStrings(actual, expected)
//	// Returns:
//	// CompareStrings: ✗ [ASSERTION_FAILED]
//	// - Expected: "hello␣mars"¶
//	// + Actual:   "hello␣world"¶
//	//
//	//   Difference at position 6:
//	//       Expected character: 'm' (U+006D)
//	//       Actual character:   'w' (U+0077)
func CompareStrings(actual, expected string) string {
	return compareStringsInternal(actual, expected, true)
}

// CompareStringsRaw provides a test framework style comparison between actual and expected strings
// without converting invisible characters to visible symbols. Useful when strings are already formatted
// or when you want to see the raw content.
//
// Returns a formatted string showing:
// - Match status (✓ [MATCH] or ✗ [ASSERTION_FAILED])
// - Both strings without visualization
// - Position of first difference
// - Unicode code points of differing characters
//
// Example:
//
//	actual := "hello world"
//	expected := "hello mars"
//	result := CompareStringsRaw(actual, expected)
//	// Returns:
//	// CompareStrings: ✗ [ASSERTION_FAILED]
//	// - Expected: "hello mars"¶
//	// + Actual:   "hello world"¶
//	//
//	//   Difference at position 6:
//	//       Expected character: 'm' (U+006D)
//	//       Actual character:   'w' (U+0077)
func CompareStringsRaw(actual, expected string) string {
	return compareStringsInternal(actual, expected, false)
}
