package text

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// CompareStrings provides a test framework style comparison between actual and expected strings
// with detailed diff highlighting for testing purposes. Converts invisible characters to visible symbols.
func CompareStrings(actual, expected string) string {
	return compareStringsInternal(actual, expected, true)
}

// CompareStringsRaw provides a test framework style comparison between actual and expected strings
// without converting invisible characters to visible symbols. Use when strings are already visualized.
func CompareStringsRaw(actual, expected string) string {
	return compareStringsInternal(actual, expected, false)
}

// compareStringsInternal is the internal implementation that handles both visualized and raw comparisons
func compareStringsInternal(actual, expected string, visualize bool) string {
	var result strings.Builder

	// Determine if strings match
	match := actual == expected

	if match {
		result.WriteString("CompareStrings: ✓ [MATCH]\n")
		if visualize {
			result.WriteString(fmt.Sprintf("  Expected: %s\n", visualizeLine(expected)))
			result.WriteString(fmt.Sprintf("  Actual:   %s\n", visualizeLine(actual)))
		} else {
			result.WriteString(fmt.Sprintf("  Expected: \"%s\"¶\n", expected))
			result.WriteString(fmt.Sprintf("  Actual:   \"%s\"¶\n", actual))
		}
	} else {
		result.WriteString("CompareStrings: ✗ [ASSERTION_FAILED]\n")
		if visualize {
			result.WriteString(fmt.Sprintf("- Expected: %s\n", visualizeLine(expected)))
			result.WriteString(fmt.Sprintf("+ Actual:   %s\n", visualizeLine(actual)))
		} else {
			result.WriteString(fmt.Sprintf("- Expected: \"%s\"¶\n", expected))
			result.WriteString(fmt.Sprintf("+ Actual:   \"%s\"¶\n", actual))
		}
		result.WriteString("\n")

		// Find and show differences
		diffDetails := characterDiff(actual, expected)
		if diffDetails != "no character differences found" {
			result.WriteString("  " + diffDetails + "\n")
		}
	}

	return result.String()
}

// visualizeLine converts a string to show invisible characters and provides better formatting
func visualizeLine(line string) string {
	if line == "" {
		return "<empty>¶"
	}

	var builder strings.Builder
	builder.WriteRune('"') // Start with quote for clarity

	for _, r := range line {
		switch r {
		case ' ':
			builder.WriteRune('␣') // Visible space
		case '\t':
			builder.WriteRune('␉') // Visible tab
		case '\r':
			builder.WriteRune('␍') // Visible carriage return
		case '\n':
			builder.WriteRune('␊') // Visible line feed
		case '\v':
			builder.WriteRune('␋') // Visible vertical tab
		case '\f':
			builder.WriteRune('␌') // Visible form feed
		default:
			builder.WriteRune(r)
		}
	}

	builder.WriteRune('"') // End with quote
	builder.WriteRune('¶') // End of line marker
	return builder.String()
}

// characterDiff shows detailed information about where characters differ between actual and expected strings
func characterDiff(actual, expected string) string {
	actualRunes := []rune(actual)
	expectedRunes := []rune(expected)
	minLen := min(len(actualRunes), len(expectedRunes))

	var result strings.Builder

	// Find first difference
	firstDiff := -1
	for i := 0; i < minLen; i++ {
		if actualRunes[i] != expectedRunes[i] {
			firstDiff = i
			break
		}
	}

	if firstDiff != -1 {
		expectedChar := expectedRunes[firstDiff]
		actualChar := actualRunes[firstDiff]
		result.WriteString(fmt.Sprintf("Difference at position %d:\n", firstDiff))
		result.WriteString(fmt.Sprintf("    Expected character: '%c' (U+%04X)\n", expectedChar, expectedChar))
		result.WriteString(fmt.Sprintf("    Actual character:   '%c' (U+%04X)", actualChar, actualChar))
	} else if len(actualRunes) != len(expectedRunes) {
		result.WriteString(fmt.Sprintf("Length difference:\n"))
		result.WriteString(fmt.Sprintf("    Expected length: %d\n", len(expectedRunes)))
		result.WriteString(fmt.Sprintf("    Actual length:   %d", len(actualRunes)))
	} else {
		return "no character differences found"
	}

	return result.String()
}

// maxLineLength finds the maximum line length in an array of strings
func maxLineLength(lines []string) int {
	maxLen := 0
	for _, line := range lines {
		if length := utf8.RuneCountInString(line); length > maxLen {
			maxLen = length
		}
	}
	return maxLen
}

// compareStringsSideBySide is kept for backwards compatibility but not used in the main function
func compareStringsSideBySide(a string, b string) string {
	aArr := strings.Split(a, "\n")
	bArr := strings.Split(b, "\n")
	maxRows := max(len(aArr), len(bArr))

	// Calculate column widths
	maxAWidth := max(utf8.RuneCountInString("A: Content"), maxLineLength(aArr)+10) // +10 for quotes and symbols
	maxBWidth := max(utf8.RuneCountInString("B: Content"), maxLineLength(bArr)+10)

	var result strings.Builder

	// Header
	result.WriteString("=== Side-by-Side Comparison ===\n")
	headerA := "A: Content"
	headerB := "B: Content"
	result.WriteString(fmt.Sprintf("%-*s | %s\n", maxAWidth, headerA, headerB))
	result.WriteString(fmt.Sprintf("%s-+-%s\n",
		strings.Repeat("-", maxAWidth),
		strings.Repeat("-", maxBWidth)))

	for i := 0; i < maxRows; i++ {
		var aLine, bLine string
		var aContent, bContent string

		if i < len(aArr) {
			aLine = aArr[i]
			aContent = visualizeLine(aLine)
		} else {
			aContent = "<missing>"
		}

		if i < len(bArr) {
			bLine = bArr[i]
			bContent = visualizeLine(bLine)
		} else {
			bContent = "<missing>"
		}

		// Status indicator
		var indicator string
		if i >= len(aArr) {
			indicator = "➕" // Extra in B
		} else if i >= len(bArr) {
			indicator = "➖" // Missing in B
		} else if aLine == bLine {
			indicator = "✓" // Match
		} else {
			indicator = "✗" // Diff
		}

		result.WriteString(fmt.Sprintf("%-*s %s %s\n",
			maxAWidth, aContent, indicator, bContent))
	}

	return result.String()
}
