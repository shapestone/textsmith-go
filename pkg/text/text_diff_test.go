package text

import (
	"fmt"
	"strings"
	"testing"
)

// Acceptance tests for public API behavior and exact output formatting

func TestDiff_OutputFormat_WithTabSpaceDifference_ShowsCorrectFormat(t *testing.T) {
	// Given
	a := "hello\tworld"
	b := "hello world"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␉world ≠ hello␣world|
		|     △             △     |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithIdenticalStrings_ShowsHeaderOnly(t *testing.T) {
	// Given
	a := "hello world"
	b := "hello world"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␣world | hello␣world|
	`)

	// - Verify match
	if !match {
		t.Fatalf("Expect diff match to be true")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithSimpleDifference_ShowsExactFormat(t *testing.T) {
	// Given
	a := "hello"
	b := "help!"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|hello    ≠ help!   |
		|   △          △    |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithExpectedLongerThanActual_ShowsLeftArrow(t *testing.T) {
	// Given
	a := "line1\nline2\nline3"
	b := "line1\nline2"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    ←         |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithActualLongerThanExpected_ShowsRightArrow(t *testing.T) {
	// Given
	a := "line1\nline2"
	b := "line1\nline2\nline3"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|         → line3   |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithEmptyLineInExpected_ShowsEmptyLineSymbol(t *testing.T) {
	// Given
	a := "line1\nline2\n"
	b := "line1\nline2"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|␤        ←         |
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithComplexWhitespace_ShowsAllWhitespaceSymbols(t *testing.T) {
	// Given
	a := "hello\t world"
	b := "hello  world"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected     | Actual      |
		|------------ | ------------|
		|hello␉␣world ≠ hello␣␣world|
		|     △              △      |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithDifferentLengthStrings_ShowsProperAlignment(t *testing.T) {
	// Given
	a := "short"
	b := "very long string"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected         | Actual          |
		|---------------- | ----------------|
		|short            ≠ very␣long␣string|
		|△                  △               |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithMultipleSpacesAndTabs_ShowsAllInvisibleChars(t *testing.T) {
	// Given
	a := "hello\t\t world"
	b := "hello   world"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected      | Actual       |
		|------------- | -------------|
		|hello␉␉␣world ≠ hello␣␣␣world|
		|     △               △       |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithOnlyWhitespace_ShowsOnlyWhitespaceSymbols(t *testing.T) {
	// Given
	a := "   "
	b := "\t\t\t"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|␣␣␣      ≠ ␉␉␉     |
		|△          △       |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithEmptyStrings_ShowsEmptyComparison(t *testing.T) {
	// Given
	a := ""
	b := ""

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|         |         |
	`)

	// - Verify match
	if !match {
		t.Fatalf("Expect diff match to be true")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithMultilineStrings_ShowsCorrectStructure(t *testing.T) {
	// Given
	a := "first line\nsecond line"
	b := "first line\nsecond line"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|first␣line  | first␣line |
		|second␣line | second␣line|
	`)

	// - Verify match
	if !match {
		t.Fatalf("Expect diff match to be true")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	a := "Hello 🌍 World"
	b := "Hello 🚀 World"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected      | Actual       |
		|------------- | -------------|
		|Hello␣🌍␣World ≠ Hello␣🚀␣World|
		|      △               △      |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_OutputFormat_WithMixedWhitespaceTypes_ShowsAllCorrectly(t *testing.T) {
	// Given
	a := "line1\n\tline2\n   line3"
	b := "line1\n line2\n\tline3"

	// When
	actual, match := Diff(a, b)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|␉line2   ≠ ␣line2  |
		|△          △       |
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_CrossPlatform_WindowsVsUnixLineEndings_ShouldMatch(t *testing.T) {
	// Given
	unixText := "line1\nline2\nline3"
	windowsText := "line1\r\nline2\r\nline3"

	// When
	actual, match := Diff(unixText, windowsText)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    | line3   |
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected Unix and Windows line endings to match after normalization")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_CrossPlatform_ClassicMacVsUnixLineEndings_ShouldMatch(t *testing.T) {
	// Given
	unixText := "line1\nline2\nline3"
	macText := "line1\rline2\rline3"

	// When
	actual, match := Diff(unixText, macText)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    | line3   |
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected Unix and Classic Mac line endings to match after normalization")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_CrossPlatform_MixedLineEndings_ShouldNormalize(t *testing.T) {
	// Given - text with mixed line endings
	mixedText := "line1\nline2\r\nline3\rline4"
	normalText := "line1\nline2\nline3\nline4"

	// When
	actual, match := Diff(mixedText, normalText)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    | line3   |
		|line4    | line4   |
	`)

	// Should match because all line endings are normalized to \n
	if !match {
		t.Fatalf("Expected mixed line endings to be normalized and match")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_CrossPlatform_EmptyLinesWithDifferentEndings_ShouldMatch(t *testing.T) {
	// Given
	unixText := "line1\n\nline3"
	windowsText := "line1\r\n\r\nline3"

	// When
	actual, match := Diff(unixText, windowsText)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|         |         |
		|line3    | line3   |
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected texts with different line endings around empty lines to match")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

func TestDiff_CrossPlatform_TrailingLineEndingsDifferent_ShouldMatch(t *testing.T) {
	// Given
	unixText := "line1\nline2\n"
	windowsText := "line1\r\nline2\r\n"

	// When
	actual, match := Diff(unixText, windowsText)

	// Then
	expected := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|         |         |
		||
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected texts with different trailing line endings to match")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareMultilineStrings(actual, expected))
	}
}

// compareMultilineStrings compares two multi-line strings with detailed analysis
// and provides helpful debugging information when strings don't match.
func compareMultilineStrings(actual, expected string) string {
	var result strings.Builder
	result.WriteString("=== Multi-line String Comparison ===\n\n")

	// Overall match status
	if actual == expected {
		result.WriteString("✓ Strings match completely\n\n")
		return result.String()
	}

	result.WriteString("✗ Strings differ\n\n")

	// First, show the raw string comparison for overall differences
	result.WriteString("Overall string comparison:\n")
	overallComparison, _ := Diff(actual, expected)
	indentedOverall := strings.ReplaceAll(overallComparison, "\n", "\n  ")
	result.WriteString("  " + indentedOverall)
	result.WriteString("\n")

	// Then do line-by-line analysis
	actualLines := strings.Split(actual, "\n")
	expectedLines := strings.Split(expected, "\n")
	maxLines := max(len(actualLines), len(expectedLines))

	result.WriteString("Line-by-line analysis:\n")

	hasLineDifferences := false

	for i := 0; i < maxLines; i++ {
		var actualLine, expectedLine string
		var actualExists, expectedExists bool

		// Get lines or empty string if beyond bounds
		if i < len(actualLines) {
			actualLine = actualLines[i]
			actualExists = true
		}
		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
			expectedExists = true
		}

		// Check various difference conditions
		lineMatches := actualLine == expectedLine
		isExtraLine := !actualExists || !expectedExists

		if !lineMatches || isExtraLine {
			hasLineDifferences = true

			if isExtraLine {
				if !actualExists {
					result.WriteString(fmt.Sprintf("Line %d: Missing in actual (expected: %q)\n", i+1, expectedLine))
				} else {
					result.WriteString(fmt.Sprintf("Line %d: Extra in actual (actual: %q)\n", i+1, actualLine))
				}
			} else {
				result.WriteString(fmt.Sprintf("Line %d differs:\n", i+1))
				lineComparison, _ := Diff(actualLine, expectedLine)
				indentedComparison := strings.ReplaceAll(lineComparison, "\n", "\n    ")
				result.WriteString("    " + indentedComparison)
				result.WriteString("\n")
			}
		}
	}

	// Handle the case where strings differ but all lines match (trailing newline differences)
	if !hasLineDifferences && actual != expected {
		result.WriteString("No line content differences found.\n")
		result.WriteString("Difference is likely in trailing characters (newlines, spaces, etc.)\n")

		// Show character-by-character comparison
		result.WriteString("\nCharacter-by-character analysis:\n")
		result.WriteString(fmt.Sprintf("Actual length: %d, Expected length: %d\n", len(actual), len(expected)))

		minLen := min(len(actual), len(expected))
		for i := 0; i < minLen; i++ {
			if actual[i] != expected[i] {
				result.WriteString(fmt.Sprintf("First difference at position %d:\n", i))
				result.WriteString(fmt.Sprintf("  Actual: %q (0x%02X)\n", actual[i], actual[i]))
				result.WriteString(fmt.Sprintf("  Expected: %q (0x%02X)\n", expected[i], expected[i]))
				break
			}
		}

		if len(actual) != len(expected) {
			result.WriteString("Length difference detected:\n")
			if len(actual) > len(expected) {
				result.WriteString(fmt.Sprintf("  Actual has %d extra characters: %q\n",
					len(actual)-len(expected), actual[len(expected):]))
			} else {
				result.WriteString(fmt.Sprintf("  Expected has %d extra characters: %q\n",
					len(expected)-len(actual), expected[len(actual):]))
			}
		}
	}

	// Summary
	result.WriteString(fmt.Sprintf("\nSummary: Total lines compared: %d\n", maxLines))
	result.WriteString(fmt.Sprintf("Actual lines: %d, Expected lines: %d\n", len(actualLines), len(expectedLines)))

	return result.String()
}
