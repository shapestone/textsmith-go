package text

import (
	"strings"
	"unicode/utf8"
)

// DiffResult represents the result of comparing two strings
type DiffResult struct {
	Lines         []DiffLine
	ExpectedWidth int
	ActualWidth   int
	Match         bool
	HasTrailingNL bool
}

// DiffLine represents a single line comparison
type DiffLine struct {
	Expected string
	Actual   string
	Status   DiffStatus
}

// DiffStatus indicates the type of difference in a line
type DiffStatus int

const (
	DiffStatusEqual DiffStatus = iota
	DiffStatusDifferent
	DiffStatusMissingInActual
	DiffStatusMissingInExpected
)

// rpad is a right space padding function
func rpad(str string, length int) string {
	rc := length - utf8.RuneCountInString(str)
	if rc <= 0 {
		return str
	}
	return str + strings.Repeat(` `, rc)
}

// stringDiff compares two strings and returns a visual indicator of where they differ
func stringDiff(a, b string) string {
	ar := []rune(a)
	br := []rune(b)
	ml := min(len(ar), len(br))
	i := 0
	for i < ml {
		if ar[i] != br[i] {
			break
		}
		i++
	}
	return strings.Repeat(" ", i) + "\u25B3"
}

// normalizeLineEndings converts all line endings to \n for consistent processing
func normalizeLineEndings(text string) string {
	// Replace Windows CRLF (\r\n) with LF (\n) first
	// This must be done before replacing standalone \r to avoid double conversion
	text = strings.ReplaceAll(text, "\r\n", "\n")
	// Replace remaining CR (\r) with LF (\n) for classic Mac compatibility
	text = strings.ReplaceAll(text, "\r", "\n")

	return text
}

func showWhitespaces(orig string) string {
	var builder strings.Builder

	for _, r := range orig {
		switch r {
		case ' ':
			builder.WriteRune('\u2423')
		case '\t':
			builder.WriteRune('\u2409')
		case '\v':
			builder.WriteRune('\u240B')
		case '\f':
			builder.WriteRune('\u240C')
		case '\r':
			// Show carriage return as a visible symbol
			builder.WriteRune('\u240D') // ␍ symbol for CR
		default:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// splitLines splits text into lines while handling cross-platform line endings
func splitLines(text string) []string {
	// Only normalize line endings for consistent processing
	normalized := normalizeLineEndings(text)

	// Handle empty string case
	if normalized == "" {
		return []string{""}
	}

	lines := strings.Split(normalized, "\n")

	// Remove the final empty string only if it results from a trailing newline
	// This preserves the correct line count for strings that are all newlines
	if len(lines) > 1 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

// hasTrailingNewline checks if the original text ends with a newline character
func hasTrailingNewline(text string) bool {
	if len(text) == 0 {
		return false
	}
	normalized := normalizeLineEndings(text)
	return strings.HasSuffix(normalized, "\n")
}

// runesEqual compares two rune slices for equality
func runesEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// computeDiff performs the diff computation and returns a structured result
func computeDiff(expected string, actual string) DiffResult {
	expectedHasTrailing := hasTrailingNewline(expected)
	actualHasTrailing := hasTrailingNewline(actual)

	expectedArr := splitLines(expected)
	actualArr := splitLines(actual)

	// Convert lines to runes for proper Unicode comparison
	expectedRunes := make([][]rune, len(expectedArr))
	actualRunes := make([][]rune, len(actualArr))

	for i, line := range expectedArr {
		expectedRunes[i] = []rune(line)
	}
	for i, line := range actualArr {
		actualRunes[i] = []rune(line)
	}

	// Calculate maximum width for both columns based on visible characters
	expectedWidth := utf8.RuneCountInString("Expected")
	for _, s := range expectedArr {
		visible := showWhitespaces(s)
		if w := utf8.RuneCountInString(visible); w > expectedWidth {
			expectedWidth = w
		}
	}

	actualWidth := utf8.RuneCountInString("Actual")
	for _, s := range actualArr {
		visible := showWhitespaces(s)
		if w := utf8.RuneCountInString(visible); w > actualWidth {
			actualWidth = w
		}
	}

	// Use the same width for both columns (maximum of both)
	maxWidth := max(expectedWidth, actualWidth)

	minVal := min(len(expectedRunes), len(actualRunes))
	var lines []DiffLine
	match := true
	foundDifference := false

	// Determine if we should add trailing newlines based on input
	shouldAddTrailingNewline := expectedHasTrailing || actualHasTrailing

	// Compare common lines using rune comparison
	for i := 0; i < minVal; i++ {
		if runesEqual(expectedRunes[i], actualRunes[i]) {
			lines = append(lines, DiffLine{
				Expected: expectedArr[i],
				Actual:   actualArr[i],
				Status:   DiffStatusEqual,
			})
		} else {
			lines = append(lines, DiffLine{
				Expected: expectedArr[i],
				Actual:   actualArr[i],
				Status:   DiffStatusDifferent,
			})
			match = false
			foundDifference = true
			break // Stop at first difference for large texts
		}
	}

	// Only process remaining lines if we haven't found a difference yet OR if the strings have different lengths
	if !foundDifference || len(expectedRunes) != len(actualRunes) {
		// Handle remaining lines in expected
		for i := minVal; i < len(expectedRunes); i++ {
			expectedStr := expectedArr[i]
			lines = append(lines, DiffLine{
				Expected: expectedStr,
				Actual:   "",
				Status:   DiffStatusMissingInActual,
			})
			match = false
		}

		// Handle remaining lines in actual
		for i := minVal; i < len(actualRunes); i++ {
			actualStr := actualArr[i]
			lines = append(lines, DiffLine{
				Expected: "",
				Actual:   actualStr,
				Status:   DiffStatusMissingInExpected,
			})
			match = false
		}
	}

	// Handle trailing newlines - this section needs to be outside the previous conditional
	// to ensure it's always processed when there are trailing newline differences
	if expectedHasTrailing != actualHasTrailing {
		if expectedHasTrailing && !actualHasTrailing {
			lines = append(lines, DiffLine{
				Expected: `␤`,
				Actual:   "",
				Status:   DiffStatusMissingInActual,
			})
			match = false
		} else if !expectedHasTrailing && actualHasTrailing {
			lines = append(lines, DiffLine{
				Expected: "",
				Actual:   `␤`,
				Status:   DiffStatusMissingInExpected,
			})
			match = false
		}
	} else if expectedHasTrailing && actualHasTrailing && !foundDifference && len(expectedRunes) == len(actualRunes) {
		// Both have trailing newlines, no content differences found, and same number of lines
		// Add an empty line to represent the trailing newline effect
		lines = append(lines, DiffLine{
			Expected: "",
			Actual:   "",
			Status:   DiffStatusEqual,
		})
	}

	return DiffResult{
		Lines:         lines,
		ExpectedWidth: maxWidth,
		ActualWidth:   maxWidth,
		Match:         match,
		HasTrailingNL: shouldAddTrailingNewline,
	}
}

// renderDiff converts a DiffResult into a visual string representation
func renderDiff(result DiffResult) string {
	var sb strings.Builder
	width := result.ExpectedWidth

	// Header
	sb.WriteString(rpad("Expected", width) + ` | ` + rpad("Actual", width) + "\n")
	sb.WriteString(strings.Repeat(`-`, width) + ` | ` + strings.Repeat(`-`, width) + "\n")

	// Content lines
	for _, line := range result.Lines {
		// Apply whitespace visualization only during rendering
		expectedVisible := showWhitespaces(line.Expected)
		actualVisible := showWhitespaces(line.Actual)

		switch line.Status {
		case DiffStatusEqual:
			sb.WriteString(rpad(expectedVisible, width) + ` | ` + rpad(actualVisible, width) + "\n")
		case DiffStatusDifferent:
			sb.WriteString(rpad(expectedVisible, width) + " \u2260 " + rpad(actualVisible, width) + "\n")
			sd := stringDiff(expectedVisible, actualVisible)
			line := rpad(sd, width) + `   ` + rpad(sd, width)
			if result.HasTrailingNL {
				line += "\n"
			}
			sb.WriteString(line)
			return sb.String()
		case DiffStatusMissingInActual:
			line := rpad(expectedVisible, width) + " \u2190 " + rpad(actualVisible, width)
			if result.HasTrailingNL {
				line += "\n"
			}
			sb.WriteString(line)
		case DiffStatusMissingInExpected:
			line := rpad(expectedVisible, width) + " \u2192 " + rpad(actualVisible, width)
			if result.HasTrailingNL {
				line += "\n"
			}
			sb.WriteString(line)
		}
	}

	// For identical strings, remove final trailing newline if input doesn't have trailing newlines
	output := sb.String()
	if result.Match && !result.HasTrailingNL && strings.HasSuffix(output, "\n") {
		output = strings.TrimSuffix(output, "\n")
	}

	return output
}

// Diff compares two strings and outputs a diff format and a boolean value to indicate if the two strings matched
func Diff(expected string, actual string) (string, bool) {
	result := computeDiff(expected, actual)
	return renderDiff(result), result.Match
}
