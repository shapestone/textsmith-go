package text

import (
	"fmt"
	"strings"
	"testing"
)

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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if !match {
		t.Fatalf("Expect diff match to be true")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if !match {
		t.Fatalf("Expect diff match to be true")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if !match {
		t.Fatalf("Expect diff match to be true")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// - Verify match
	if match {
		t.Fatalf("Expect diff match to be false")
	}

	// - Verify diff format
	if actual != expected {
		t.Fatalf("Actual vs Expected is not matching:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected Unix and Windows line endings to match after normalization")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected Unix and Classic Mac line endings to match after normalization")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// Should match because all line endings are normalized to \n
	if !match {
		t.Fatalf("Expected mixed line endings to be normalized and match")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareStrings(actual, expected))
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
		||
	`)

	// Should match because line endings are normalized
	if !match {
		t.Fatalf("Expected texts with different line endings around empty lines to match")
	}

	if actual != expected {
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareStrings(actual, expected))
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
		t.Fatalf("Actual vs Expected format mismatch:\n\n%v", compareStrings(actual, expected))
	}
}

// Test helper functions
func compareStrings(a string, b string) string {
	aArr := strings.Split(a, "\n")
	bArr := strings.Split(b, "\n")
	rows := max(len(aArr), len(bArr))

	var result strings.Builder

	for i := 0; i < rows; i++ {
		var aLine, bLine string

		// Get line from a, or empty if beyond bounds
		if i < len(aArr) {
			aLine = aArr[i]
		}

		// Get line from b, or empty if beyond bounds
		if i < len(bArr) {
			bLine = bArr[i]
		}

		// Add end-of-line markers to visualize whitespace
		aLineWithEOL := aLine + "¶"
		bLineWithEOL := bLine + "¶"

		// Check if lines match
		match := aLine == bLine

		// Format the comparison
		status := "✓"
		if !match {
			status = "✗"
		}

		result.WriteString(fmt.Sprintf("Row %d: %s\n", i+1, status))
		result.WriteString(fmt.Sprintf("  A: %s\n", aLineWithEOL))
		result.WriteString(fmt.Sprintf("  B: %s\n", bLineWithEOL))
		result.WriteString("\n")
	}

	return result.String()
}
