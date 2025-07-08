package text

import (
	"strings"
	"testing"
)

// compareMultilineStrings provides a visual diff for debugging test failures
func compareMultilineStrings(diffOutput, expected string) string {
	debugDiff, _ := Diff(expected, diffOutput)
	return debugDiff
}

// Tests focused on renderDiff logic through the public Diff function
func TestDiff_RenderLogic_WithEqualLines_ShowsCorrectFormat(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello world"

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if !isMatch {
		t.Fatalf("Expected isMatch to be true for identical strings")
	}

	expectedOutput := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␣world | hello␣world|
	`)

	if diffOutput != expectedOutput {
		t.Fatalf("Rendered output does not match expected:\n\n%s", compareMultilineStrings(diffOutput, expectedOutput))
	}
}

func TestDiff_RenderLogic_WithDifferentLines_ShowsCorrectFormat(t *testing.T) {
	// Given
	expected := "hello"
	actual := "help!"

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if isMatch {
		t.Fatalf("Expected isMatch to be false for different strings")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|hello    ≠ help!   |
		|   △          △    |
	`)

	if diffOutput != expectedOutput {
		t.Fatalf("Rendered output does not match expected:\n\n%s", compareMultilineStrings(diffOutput, expectedOutput))
	}
}

func TestDiff_RenderLogic_WithTabSpaceDifference_ShowsCorrectFormat(t *testing.T) {
	// Given
	expected := "hello\tworld"
	actual := "hello world"

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if isMatch {
		t.Fatalf("Expected isMatch to be false for tab vs space difference")
	}

	expectedOutput := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␉world ≠ hello␣world|
		|     △             △     |
	`)

	if diffOutput != expectedOutput {
		t.Fatalf("Rendered output does not match expected:\n\n%s", compareMultilineStrings(diffOutput, expectedOutput))
	}
}

func TestDiff_RenderLogic_WithMissingInActual_ShowsLeftArrow(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2"

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if isMatch {
		t.Fatalf("Expected isMatch to be false when lines are missing in actual")
	}

	// Should contain left arrow for missing in actual
	if !strings.Contains(diffOutput, "←") {
		t.Fatalf("Expected output to contain '←' for missing in actual, got: %s", diffOutput)
	}

	// Should show the missing line
	if !strings.Contains(diffOutput, "line3") {
		t.Fatalf("Expected output to contain 'line3', got: %s", diffOutput)
	}

	// Should have proper header format
	if !strings.Contains(diffOutput, "Expected") || !strings.Contains(diffOutput, "Actual") {
		t.Fatalf("Expected output to contain proper headers, got: %s", diffOutput)
	}
}

func TestDiff_RenderLogic_WithMissingInExpected_ShowsRightArrow(t *testing.T) {
	// Given
	expected := "line1\nline2"
	actual := "line1\nline2\nline3"

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if isMatch {
		t.Fatalf("Expected isMatch to be false when lines are missing in expected")
	}

	// Should contain right arrow for missing in expected
	if !strings.Contains(diffOutput, "→") {
		t.Fatalf("Expected output to contain '→' for missing in expected, got: %s", diffOutput)
	}

	// Should show the extra line
	if !strings.Contains(diffOutput, "line3") {
		t.Fatalf("Expected output to contain 'line3', got: %s", diffOutput)
	}
}

func TestDiff_RenderLogic_WithWhitespaceSymbols_ShowsCorrectSymbols(t *testing.T) {
	// Given
	expected := "hello\t\r\vworld"
	actual := "hello     world"

	// When
	diffOutput, _ := Diff(expected, actual)

	// Then
	// Should show tab symbol
	if !strings.Contains(diffOutput, "␉") {
		t.Errorf("Expected output to contain tab symbol '␉', got: %s", diffOutput)
	}

	// Should show carriage return symbol
	if !strings.Contains(diffOutput, "␍") {
		t.Errorf("Expected output to contain carriage return symbol '␍', got: %s", diffOutput)
	}

	// Should show vertical tab symbol
	if !strings.Contains(diffOutput, "␋") {
		t.Errorf("Expected output to contain vertical tab symbol '␋', got: %s", diffOutput)
	}

	// Should show space symbols
	if !strings.Contains(diffOutput, "␣") {
		t.Errorf("Expected output to contain space symbol '␣', got: %s", diffOutput)
	}
}

func TestDiff_RenderLogic_WithHeaderAlignment_ShowsCorrectPadding(t *testing.T) {
	// Given
	expected := "short"
	actual := "this is a much longer string"

	// When
	diffOutput, _ := Diff(expected, actual)

	// Then
	lines := strings.Split(diffOutput, "\n")
	if len(lines) < 2 {
		t.Fatalf("Expected at least 2 lines for header, got %d", len(lines))
	}

	headerLine := lines[0]
	dashLine := lines[1]

	// Header should contain "Expected" and "Actual"
	if !strings.Contains(headerLine, "Expected") || !strings.Contains(headerLine, "Actual") {
		t.Fatalf("Expected header line to contain 'Expected' and 'Actual', got: %s", headerLine)
	}

	// Dash line should contain dashes and separator
	if !strings.Contains(dashLine, "-") || !strings.Contains(dashLine, "|") {
		t.Fatalf("Expected dash line to contain dashes and separator, got: %s", dashLine)
	}

	// Both lines should have similar length (proper alignment)
	headerLen := len(headerLine)
	dashLen := len(dashLine)
	if abs(headerLen-dashLen) > 2 { // Allow small difference for rounding
		t.Fatalf("Expected header and dash lines to have similar length, got header: %d, dash: %d", headerLen, dashLen)
	}
}

func TestDiff_RenderLogic_WithTrailingNewlines_ShowsCorrectSymbol(t *testing.T) {
	// Given
	expected := "hello\nworld"
	actual := "hello\nworld\n"

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if isMatch {
		t.Fatalf("Expected isMatch to be false when trailing newlines differ")
	}

	// Should show trailing newline symbol
	if !strings.Contains(diffOutput, "␤") {
		t.Fatalf("Expected output to contain trailing newline symbol '␤', got: %s", diffOutput)
	}
}

func TestDiff_RenderLogic_WithMacVsWindowsLineEndings_ShowsCorrectSymbols(t *testing.T) {
	// Given
	macLineEnding := "Hello, World!\n"       // macOS/Unix line ending
	windowsLineEnding := "Hello, World!\r\n" // Windows line ending

	// When
	diffOutput, isMatch := Diff(macLineEnding, windowsLineEnding)

	// Then
	if isMatch {
		t.Fatalf("Expected isMatch to be false when line endings differ (macOS vs Windows)")
	}

	// Should show carriage return symbol for Windows line ending
	if !strings.Contains(diffOutput, "␍") {
		t.Fatalf("Expected output to contain carriage return symbol '␍' for Windows line ending, got: %s", diffOutput)
	}

	// Should show the difference indicator
	if !strings.Contains(diffOutput, "≠") {
		t.Fatalf("Expected output to contain difference indicator '≠', got: %s", diffOutput)
	}

	// Should show triangle pointer indicating where the difference starts
	if !strings.Contains(diffOutput, "△") {
		t.Fatalf("Expected output to contain triangle pointer '△', got: %s", diffOutput)
	}

	// Should contain the base text "Hello, World!"
	if !strings.Contains(diffOutput, "Hello,␣World!") {
		t.Fatalf("Expected output to contain 'Hello,␣World!' with visible space, got: %s", diffOutput)
	}
}

func TestDiff_RenderLogic_WithUnicodeContent_MaintainsAlignment(t *testing.T) {
	// Given
	expected := "hello 世界"
	actual := "hello world"

	// When
	diffOutput, _ := Diff(expected, actual)

	// Then
	lines := strings.Split(diffOutput, "\n")
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines, got %d", len(lines))
	}

	// Check that Unicode characters are preserved
	if !strings.Contains(diffOutput, "世界") {
		t.Fatalf("Expected output to contain Unicode characters '世界', got: %s", diffOutput)
	}

	// Check that alignment is maintained (content line should have proper separators)
	contentLine := lines[2] // Skip header and dash line
	if !strings.Contains(contentLine, "≠") {
		t.Fatalf("Expected content line to contain difference indicator '≠', got: %s", contentLine)
	}
}

func TestDiff_RenderLogic_WithEmptyStrings_ShowsHeaderOnly(t *testing.T) {
	// Given
	expected := ""
	actual := ""

	// When
	diffOutput, isMatch := Diff(expected, actual)

	// Then
	if !isMatch {
		t.Fatalf("Expected isMatch to be true for empty strings")
	}

	lines := strings.Split(diffOutput, "\n")

	// Should have at least header and dash line, plus one content line for empty strings
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines for empty string diff, got %d", len(lines))
	}

	// Should contain proper headers
	if !strings.Contains(lines[0], "Expected") || !strings.Contains(lines[0], "Actual") {
		t.Fatalf("Expected proper headers in first line, got: %s", lines[0])
	}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
