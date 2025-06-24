package text

import (
	"strings"
	"unicode/utf8"
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
	// Replace Windows CRLF (\r\n) with LF (\n) first to avoid double conversion
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
	// Normalize line endings first, then apply whitespace visualization
	normalized := normalizeLineEndings(text)
	withVisibleWS := showWhitespaces(normalized)
	return strings.Split(withVisibleWS, "\n")
}

// hasTrailingNewline checks if the original text ends with a newline character
func hasTrailingNewline(text string) bool {
	if len(text) == 0 {
		return false
	}
	normalized := normalizeLineEndings(text)
	return strings.HasSuffix(normalized, "\n")
}

// Diff compares two strings and outputs a diff format and a boolean value to indicate if the two strings matched
func Diff(expected string, actual string) (string, bool) {
	expectedArr := splitLines(expected)
	// find the longest string in an array of strings
	expectedWidth := utf8.RuneCountInString("Expected")
	for _, s := range expectedArr {
		if expectedWidth < utf8.RuneCountInString(s) {
			expectedWidth = utf8.RuneCountInString(s)
		}
	}

	actualArr := splitLines(actual)
	// find the longest string in an array of strings
	actualWidth := utf8.RuneCountInString("Actual")
	for _, s := range actualArr {
		if actualWidth < utf8.RuneCountInString(s) {
			actualWidth = utf8.RuneCountInString(s)
		}
	}
	width := max(expectedWidth, actualWidth)

	minVal := min(len(expectedArr), len(actualArr))
	var sb strings.Builder
	status := true

	// Determine if we should add trailing newlines based on input
	expectedHasTrailing := hasTrailingNewline(expected)
	actualHasTrailing := hasTrailingNewline(actual)
	shouldAddTrailingNewline := expectedHasTrailing || actualHasTrailing

	sb.WriteString(rpad("Expected", width) + ` | ` + rpad("Actual", width) + "\n")
	sb.WriteString(strings.Repeat(`-`, width) + ` | ` + strings.Repeat(`-`, width) + "\n")

	for i := 0; i < minVal; i++ {
		if expectedArr[i] == actualArr[i] {
			sb.WriteString(rpad(expectedArr[i], width) + ` | ` + rpad(actualArr[i], width) + "\n")
		} else if expectedArr[i] != actualArr[i] {
			expected, actual := expectedArr[i], actualArr[i]
			sb.WriteString(rpad(expected, width) + " \u2260 " + rpad(actual, width) + "\n")
			sd := stringDiff(expectedArr[i], actualArr[i])
			line := rpad(sd, width) + `   ` + rpad(sd, width)
			if shouldAddTrailingNewline {
				line += "\n"
			}
			sb.WriteString(line)
			return sb.String(), false
		}
	}

	if len(expectedArr) > len(actualArr) {
		expectedStr := expectedArr[minVal]
		actualStr := ``
		if utf8.RuneCountInString(expectedStr) == 0 {
			expectedStr = `␤`
		}
		line := rpad(expectedStr, width) + " \u2190 " + rpad(actualStr, width)
		if shouldAddTrailingNewline {
			line += "\n"
		}
		sb.WriteString(line)
		status = false
	} else if len(expectedArr) < len(actualArr) {
		expectedStr := ``
		actualStr := actualArr[minVal]
		if utf8.RuneCountInString(actualStr) == 0 {
			actualStr = `␤`
		}
		line := rpad(expectedStr, width) + " \u2192 " + rpad(actualStr, width)
		if shouldAddTrailingNewline {
			line += "\n"
		}
		sb.WriteString(line)
		status = false
	}

	// For identical strings, we need to remove the final trailing newline
	// if the input strings don't have trailing newlines
	result := sb.String()
	if status && !shouldAddTrailingNewline && strings.HasSuffix(result, "\n") {
		result = strings.TrimSuffix(result, "\n")
	}

	return result, status
}
