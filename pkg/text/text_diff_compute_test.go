package text

import (
	"strings"
	"testing"
)

// Tests focused on computeDiff logic through the public Diff function
func TestDiff_ComputeLogic_WithIdenticalStrings_ReturnsMatch(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello world"

	// When
	_, match := Diff(expected, actual)

	// Then
	if !match {
		t.Fatalf("Expected Diff to return true for identical strings, got %t", match)
	}
}

func TestDiff_ComputeLogic_WithDifferentStrings_ReturnsNoMatch(t *testing.T) {
	// Given
	expected := "hello"
	actual := "help!"

	// When
	_, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false for different strings, got %t", match)
	}
}

func TestDiff_ComputeLogic_WithMultilineStrings_DetectsDifference(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2\nline4"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false for different multiline strings, got %t", match)
	}

	// Should show the difference in line 3
	if !strings.Contains(output, "line3") {
		t.Errorf("Expected output to contain 'line3' from expected, got: %s", output)
	}

	if !strings.Contains(output, "line4") {
		t.Errorf("Expected output to contain 'line4' from actual, got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithMissingLinesInActual_ShowsMissingIndicator(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false when lines are missing in actual, got %t", match)
	}

	// Should show left arrow for missing in actual
	if !strings.Contains(output, "←") {
		t.Errorf("Expected output to contain '←' for missing in actual, got: %s", output)
	}

	if !strings.Contains(output, "line3") {
		t.Errorf("Expected output to contain missing 'line3', got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithMissingLinesInExpected_ShowsMissingIndicator(t *testing.T) {
	// Given
	expected := "line1\nline2"
	actual := "line1\nline2\nline3"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false when lines are missing in expected, got %t", match)
	}

	// Should show right arrow for missing in expected
	if !strings.Contains(output, "→") {
		t.Errorf("Expected output to contain '→' for missing in expected, got: %s", output)
	}

	if !strings.Contains(output, "line3") {
		t.Errorf("Expected output to contain extra 'line3', got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithTrailingNewlineDifference_DetectsDifference(t *testing.T) {
	// Given
	expected := "hello\nworld"
	actual := "hello\nworld\n"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false when trailing newlines differ, got %t", match)
	}

	// Should show trailing newline symbol
	if !strings.Contains(output, "␤") {
		t.Errorf("Expected output to contain trailing newline symbol '␤', got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithCarriageReturnLineFeed_DetectsDifference(t *testing.T) {
	// Given
	expected := "line1\nline2"
	actual := "line1\r\nline2"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false when line endings differ (\\n vs \\r\\n), got %t", match)
	}

	// Should show carriage return symbol
	if !strings.Contains(output, "␍") {
		t.Errorf("Expected output to contain carriage return symbol '␍', got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithTabsAndSpaces_DetectsDifference(t *testing.T) {
	// Given
	expected := "hello\tworld"
	actual := "hello world"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false when tabs vs spaces differ, got %t", match)
	}

	// Should show tab symbol
	if !strings.Contains(output, "␉") {
		t.Errorf("Expected output to contain tab symbol '␉', got: %s", output)
	}

	// Should show space symbol
	if !strings.Contains(output, "␣") {
		t.Errorf("Expected output to contain space symbol '␣', got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithUnicodeCharacters_HandlesCorrectly(t *testing.T) {
	// Given
	expected := "hello 世界"
	actual := "hello 世界!"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false for different Unicode strings, got %t", match)
	}

	// Should contain the Unicode characters
	if !strings.Contains(output, "世界") {
		t.Errorf("Expected output to contain Unicode characters '世界', got: %s", output)
	}
}

func TestDiff_ComputeLogic_WithEmptyStrings_ReturnsMatch(t *testing.T) {
	// Given
	expected := ""
	actual := ""

	// When
	_, match := Diff(expected, actual)

	// Then
	if !match {
		t.Fatalf("Expected Diff to return true for empty strings, got %t", match)
	}
}

func TestDiff_ComputeLogic_WithEmptyVsNonEmpty_DetectsDifference(t *testing.T) {
	// Given
	expected := ""
	actual := "hello"

	// When
	output, match := Diff(expected, actual)

	// Then
	if match {
		t.Fatalf("Expected Diff to return false for empty vs non-empty strings, got %t", match)
	}

	if !strings.Contains(output, "hello") {
		t.Errorf("Expected output to contain 'hello', got: %s", output)
	}
}
