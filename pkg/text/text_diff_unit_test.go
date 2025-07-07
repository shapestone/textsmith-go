package text

import (
	"strings"
	"testing"
)

// Unit tests for internal function behavior and data structures

func TestComputeDiff_WithIdenticalStrings_ReturnsMatchWithEqualStatus(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello world"

	// When
	result := computeDiff(expected, actual)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true for identical strings")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusEqual {
		t.Errorf("Expected DiffStatusEqual, got %v", result.Lines[0].Status)
	}
	if result.Lines[0].Expected != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result.Lines[0].Expected)
	}
	if result.Lines[0].Actual != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result.Lines[0].Actual)
	}
}

func TestComputeDiff_WithDifferentStrings_ReturnsNoMatchWithDifferentStatus(t *testing.T) {
	// Given
	expected := "hello"
	actual := "help!"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false for different strings")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusDifferent {
		t.Errorf("Expected DiffStatusDifferent, got %v", result.Lines[0].Status)
	}
	if result.Lines[0].Expected != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result.Lines[0].Expected)
	}
	if result.Lines[0].Actual != "help!" {
		t.Errorf("Expected 'help!', got '%s'", result.Lines[0].Actual)
	}
}

func TestComputeDiff_WithExpectedLonger_ReturnsMissingInActualStatus(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if len(result.Lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result.Lines))
	}

	// First two lines should be equal
	if result.Lines[0].Status != DiffStatusEqual {
		t.Errorf("Expected line 0 to be DiffStatusEqual, got %v", result.Lines[0].Status)
	}
	if result.Lines[1].Status != DiffStatusEqual {
		t.Errorf("Expected line 1 to be DiffStatusEqual, got %v", result.Lines[1].Status)
	}

	// Third line should be missing in actual
	if result.Lines[2].Status != DiffStatusMissingInActual {
		t.Errorf("Expected line 2 to be DiffStatusMissingInActual, got %v", result.Lines[2].Status)
	}
	if result.Lines[2].Expected != "line3" {
		t.Errorf("Expected 'line3', got '%s'", result.Lines[2].Expected)
	}
	if result.Lines[2].Actual != "" {
		t.Errorf("Expected empty actual, got '%s'", result.Lines[2].Actual)
	}
}

func TestComputeDiff_WithActualLonger_ReturnsMissingInExpectedStatus(t *testing.T) {
	// Given
	expected := "line1\nline2"
	actual := "line1\nline2\nline3"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if len(result.Lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result.Lines))
	}

	// First two lines should be equal
	if result.Lines[0].Status != DiffStatusEqual {
		t.Errorf("Expected line 0 to be DiffStatusEqual, got %v", result.Lines[0].Status)
	}
	if result.Lines[1].Status != DiffStatusEqual {
		t.Errorf("Expected line 1 to be DiffStatusEqual, got %v", result.Lines[1].Status)
	}

	// Third line should be missing in expected
	if result.Lines[2].Status != DiffStatusMissingInExpected {
		t.Errorf("Expected line 2 to be DiffStatusMissingInExpected, got %v", result.Lines[2].Status)
	}
	if result.Lines[2].Expected != "" {
		t.Errorf("Expected empty expected, got '%s'", result.Lines[2].Expected)
	}
	if result.Lines[2].Actual != "line3" {
		t.Errorf("Expected 'line3', got '%s'", result.Lines[2].Actual)
	}
}

func TestComputeDiff_WithEmptyStrings_ReturnsMatch(t *testing.T) {
	// Given
	expected := ""
	actual := ""

	// When
	result := computeDiff(expected, actual)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true for empty strings")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusEqual {
		t.Errorf("Expected DiffStatusEqual, got %v", result.Lines[0].Status)
	}
	if result.Lines[0].Expected != "" {
		t.Errorf("Expected empty expected, got '%s'", result.Lines[0].Expected)
	}
	if result.Lines[0].Actual != "" {
		t.Errorf("Expected empty actual, got '%s'", result.Lines[0].Actual)
	}
}

func TestComputeDiff_WithTrailingNewlines_SetsHasTrailingNLCorrectly(t *testing.T) {
	// Given
	expected := "line1\nline2\n"
	actual := "line1\nline2"

	// When
	result := computeDiff(expected, actual)

	// Then
	if !result.HasTrailingNL {
		t.Errorf("Expected HasTrailingNL to be true when expected has trailing newline")
	}
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	// Should have 3 lines: line1, line2, and empty line
	if len(result.Lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result.Lines))
	}
	if result.Lines[2].Status != DiffStatusMissingInActual {
		t.Errorf("Expected DiffStatusMissingInActual for trailing newline, got %v", result.Lines[2].Status)
	}
	if result.Lines[2].Expected != "␤" {
		t.Errorf("Expected '␤' for empty line, got '%s'", result.Lines[2].Expected)
	}
}

func TestComputeDiff_WithWhitespaceContent_StoresOriginalStrings(t *testing.T) {
	// Given
	expected := "hello\tworld"
	actual := "hello world"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusDifferent {
		t.Errorf("Expected DiffStatusDifferent, got %v", result.Lines[0].Status)
	}
	// computeDiff should store original strings, not visualized ones
	if result.Lines[0].Expected != "hello\tworld" {
		t.Errorf("Expected 'hello\\tworld', got '%s'", result.Lines[0].Expected)
	}
	if result.Lines[0].Actual != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result.Lines[0].Actual)
	}
}

func TestComputeDiff_WithCrossplatformLineEndings_NormalizesCorrectly(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\r\nline2\r\nline3"

	// When
	result := computeDiff(expected, actual)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true after line ending normalization")
	}
	if len(result.Lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result.Lines))
	}
	for i, line := range result.Lines {
		if line.Status != DiffStatusEqual {
			t.Errorf("Expected line %d to be DiffStatusEqual, got %v", i, line.Status)
		}
	}
}

func TestComputeDiff_CalculatesWidthsCorrectly(t *testing.T) {
	// Given
	expected := "short"
	actual := "very long string indeed"

	// When
	result := computeDiff(expected, actual)

	// Then
	// Width should be calculated based on visible characters (after showWhitespaces)
	expectedMinWidth := len("Expected")
	maxContentWidth := max(len("short"), len("very long string indeed"))
	expectedWidth := max(expectedMinWidth, maxContentWidth)

	if result.ExpectedWidth != expectedWidth {
		t.Errorf("Expected ExpectedWidth to be %d, got %d", expectedWidth, result.ExpectedWidth)
	}
	if result.ActualWidth != expectedWidth {
		t.Errorf("Expected ActualWidth to be %d, got %d", expectedWidth, result.ActualWidth)
	}
}

func TestComputeDiff_WithUnicodeContent_HandlesRuneCountCorrectly(t *testing.T) {
	// Given
	expected := "Hello 🌍"
	actual := "Hello 🚀"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	// computeDiff should store original strings
	if result.Lines[0].Expected != "Hello 🌍" {
		t.Errorf("Expected 'Hello 🌍', got '%s'", result.Lines[0].Expected)
	}
	if result.Lines[0].Actual != "Hello 🚀" {
		t.Errorf("Expected 'Hello 🚀', got '%s'", result.Lines[0].Actual)
	}
}

func TestComputeDiff_WithMixedLineEndingsInSameString_NormalizesAll(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2\r\nline3\rline4"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if len(result.Lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(result.Lines))
	}

	// First three lines should match
	for i := 0; i < 3; i++ {
		if result.Lines[i].Status != DiffStatusEqual {
			t.Errorf("Expected line %d to be DiffStatusEqual, got %v", i, result.Lines[i].Status)
		}
	}

	// Fourth line should be missing in expected
	if result.Lines[3].Status != DiffStatusMissingInExpected {
		t.Errorf("Expected line 3 to be DiffStatusMissingInExpected, got %v", result.Lines[3].Status)
	}
	if result.Lines[3].Actual != "line4" {
		t.Errorf("Expected 'line4', got '%s'", result.Lines[3].Actual)
	}
}

// Tests for renderDiff function

func TestRenderDiff_WithEqualLines_ReturnsHeaderAndContent(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "hello world", Actual: "hello world", Status: DiffStatusEqual},
		},
		ExpectedWidth: 11,
		ActualWidth:   11,
		Match:         true,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	expectedHeader := "Expected    | Actual     "
	if lines[0] != expectedHeader {
		t.Errorf("Expected header '%s', got '%s'", expectedHeader, lines[0])
	}

	expectedSeparator := "----------- | -----------"
	if lines[1] != expectedSeparator {
		t.Errorf("Expected separator '%s', got '%s'", expectedSeparator, lines[1])
	}

	// renderDiff should apply whitespace visualization
	expectedContent := "hello␣world | hello␣world"
	if lines[2] != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, lines[2])
	}
}

func TestRenderDiff_WithDifferentLines_ReturnsHeaderContentAndDiffPointer(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "hello", Actual: "help!", Status: DiffStatusDifferent},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         false,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	expectedHeader := "Expected | Actual  "
	if lines[0] != expectedHeader {
		t.Errorf("Expected header '%s', got '%s'", expectedHeader, lines[0])
	}

	expectedSeparator := "-------- | --------"
	if lines[1] != expectedSeparator {
		t.Errorf("Expected separator '%s', got '%s'", expectedSeparator, lines[1])
	}

	expectedContent := "hello    ≠ help!   "
	if lines[2] != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, lines[2])
	}

	expectedPointer := "   △          △    "
	if lines[3] != expectedPointer {
		t.Errorf("Expected pointer '%s', got '%s'", expectedPointer, lines[3])
	}
}

func TestRenderDiff_WithMissingInActual_ReturnsHeaderContentAndLeftArrow(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "line1", Actual: "line1", Status: DiffStatusEqual},
			{Expected: "line2", Actual: "", Status: DiffStatusMissingInActual},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         false,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	expectedContent1 := "line1    | line1   "
	if lines[2] != expectedContent1 {
		t.Errorf("Expected content line 1 '%s', got '%s'", expectedContent1, lines[2])
	}

	expectedContent2 := "line2    ←         "
	if lines[3] != expectedContent2 {
		t.Errorf("Expected content line 2 '%s', got '%s'", expectedContent2, lines[3])
	}
}

func TestRenderDiff_WithMissingInExpected_ReturnsHeaderContentAndRightArrow(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "line1", Actual: "line1", Status: DiffStatusEqual},
			{Expected: "", Actual: "line2", Status: DiffStatusMissingInExpected},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         false,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	expectedContent1 := "line1    | line1   "
	if lines[2] != expectedContent1 {
		t.Errorf("Expected content line 1 '%s', got '%s'", expectedContent1, lines[2])
	}

	expectedContent2 := "         → line2   "
	if lines[3] != expectedContent2 {
		t.Errorf("Expected content line 2 '%s', got '%s'", expectedContent2, lines[3])
	}
}

func TestRenderDiff_WithTrailingNewline_AddsNewlineToOutput(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "line1", Actual: "line1", Status: DiffStatusEqual},
			{Expected: "␤", Actual: "", Status: DiffStatusMissingInActual},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         false,
		HasTrailingNL: true,
	}

	// When
	output := renderDiff(result)

	// Then
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Expected output to end with newline when HasTrailingNL is true")
	}

	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (excluding final newline), got %d", len(lines))
	}

	expectedContent2 := "␤        ←         "
	if lines[3] != expectedContent2 {
		t.Errorf("Expected content line 2 '%s', got '%s'", expectedContent2, lines[3])
	}
}

func TestRenderDiff_WithMatchAndNoTrailingNewline_RemovesTrailingNewline(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "hello", Actual: "hello", Status: DiffStatusEqual},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         true,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	if strings.HasSuffix(output, "\n") {
		t.Errorf("Expected output to NOT end with newline when Match is true and HasTrailingNL is false")
	}
}

func TestRenderDiff_WithMatchAndTrailingNewline_KeepsTrailingNewline(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "hello", Actual: "hello", Status: DiffStatusEqual},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         true,
		HasTrailingNL: true,
	}

	// When
	output := renderDiff(result)

	// Then
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("Expected output to end with newline when HasTrailingNL is true")
	}
}

func TestRenderDiff_WithDifferentStatus_StopsAtFirstDifference(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "line1", Actual: "line1", Status: DiffStatusEqual},
			{Expected: "line2", Actual: "lineX", Status: DiffStatusDifferent},
			{Expected: "line3", Actual: "line3", Status: DiffStatusEqual}, // This should not appear in output
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         false,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	// Should only have header, separator, first equal line, different line, and pointer line
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(lines))
	}

	// Should not contain "line3"
	if strings.Contains(output, "line3") {
		t.Errorf("Expected output to stop at first difference, but found 'line3'")
	}
}

func TestRenderDiff_WithEmptyLines_HandlesCorrectly(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "", Actual: "", Status: DiffStatusEqual},
		},
		ExpectedWidth: 8,
		ActualWidth:   8,
		Match:         true,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	expectedContent := "         |         "
	if lines[2] != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, lines[2])
	}
}

func TestRenderDiff_WithUnicodeContent_HandlesRuneWidthCorrectly(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "Hello 🌍", Actual: "Hello 🚀", Status: DiffStatusDifferent},
		},
		ExpectedWidth: 12,
		ActualWidth:   12,
		Match:         false,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// Check that Unicode characters are handled correctly in alignment
	// renderDiff should apply whitespace visualization
	expectedContent := "Hello␣🌍      ≠ Hello␣🚀     "
	if lines[2] != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, lines[2])
	}
}

func TestRenderDiff_WithWhitespaceVisualization_ShowsInvisibleChars(t *testing.T) {
	// Given
	result := DiffResult{
		Lines: []DiffLine{
			{Expected: "hello\tworld", Actual: "hello world", Status: DiffStatusDifferent},
		},
		ExpectedWidth: 11,
		ActualWidth:   11,
		Match:         false,
		HasTrailingNL: false,
	}

	// When
	output := renderDiff(result)

	// Then
	lines := strings.Split(output, "\n")
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(lines))
	}

	// renderDiff should apply whitespace visualization
	expectedContent := "hello␉world ≠ hello␣world"
	if lines[2] != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, lines[2])
	}
}
