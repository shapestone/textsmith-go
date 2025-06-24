package text_test

import (
	"strings"
	"testing"

	"github.com/shapestone/textsmith/pkg/text"
)

func TestCompareStrings_WithMatchingStrings_ReturnsMatch(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello world"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match indicator, got: %s", result)
	}
	if !strings.Contains(result, `"hello␣world"¶`) {
		t.Errorf("Expected visualized string with visible space, got: %s", result)
	}
}

func TestCompareStrings_WithDifferentStrings_ReturnsAssertionFailed(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello mars"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed indicator, got: %s", result)
	}
	if !strings.Contains(result, "- Expected") {
		t.Errorf("Expected 'Expected' line with minus prefix, got: %s", result)
	}
	if !strings.Contains(result, "+ Actual") {
		t.Errorf("Expected 'Actual' line with plus prefix, got: %s", result)
	}
}

func TestCompareStrings_WithEmptyStrings_ReturnsMatch(t *testing.T) {
	// Given
	actual := ""
	expected := ""

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match indicator for empty strings, got: %s", result)
	}
	if !strings.Contains(result, "<empty>¶") {
		t.Errorf("Expected empty string visualization, got: %s", result)
	}
}

func TestCompareStrings_WithUnicodeCharacters_HandlesCorrectly(t *testing.T) {
	// Given
	actual := "Hello 世界! 🌍"
	expected := "Hello 世界! 🌍"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match indicator for Unicode strings, got: %s", result)
	}
	// Check for the visualized version with visible spaces
	if !strings.Contains(result, "Hello␣世界!␣🌍") {
		t.Errorf("Expected Unicode characters with visualized spaces to be preserved, got: %s", result)
	}
}

func TestCompareStrings_WithDifferentUnicodeCharacters_ShowsDifference(t *testing.T) {
	// Given
	actual := "Hello 世界"
	expected := "Hello 地球"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed for different Unicode, got: %s", result)
	}
	if !strings.Contains(result, "Difference at position 6") {
		t.Errorf("Expected difference position, got: %s", result)
	}
}

func TestCompareStrings_WithWhitespaceCharacters_VisualizesCorrectly(t *testing.T) {
	// Given
	actual := "hello\tworld\n"
	expected := "hello world "

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed, got: %s", result)
	}
	if !strings.Contains(result, "␉") { // Tab visualization
		t.Errorf("Expected tab visualization, got: %s", result)
	}
	if !strings.Contains(result, "␊") { // Newline visualization
		t.Errorf("Expected newline visualization, got: %s", result)
	}
	if !strings.Contains(result, "␣") { // Space visualization
		t.Errorf("Expected space visualization, got: %s", result)
	}
}

func TestVisualizeLineFunction_WithEmptyString_ReturnsEmptyVisualization(t *testing.T) {
	// This test assumes visualizeLine is exported or we're testing via CompareStrings
	// Given
	actual := ""
	expected := "not empty"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "<empty>¶") {
		t.Errorf("Expected empty string visualization, got: %s", result)
	}
}

func TestVisualizeLineFunction_WithSpecialCharacters_ShowsAllVisualizations(t *testing.T) {
	// Given
	actual := " \t\r\n\v\f"
	expected := "normal"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	specialChars := []string{"␣", "␉", "␍", "␊", "␋", "␌"}
	for _, char := range specialChars {
		if !strings.Contains(result, char) {
			t.Errorf("Expected special character visualization %s, got: %s", char, result)
		}
	}
}

func TestCharacterDiff_WithFirstCharacterDifferent_ShowsPositionZero(t *testing.T) {
	// Given
	actual := "abc"
	expected := "xyz"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "Difference at position 0") {
		t.Errorf("Expected difference at position 0, got: %s", result)
	}
	if !strings.Contains(result, "Expected character: 'x'") {
		t.Errorf("Expected character information, got: %s", result)
	}
	if !strings.Contains(result, "Actual character:   'a'") {
		t.Errorf("Expected actual character information, got: %s", result)
	}
}

func TestCharacterDiff_WithMiddleCharacterDifferent_ShowsCorrectPosition(t *testing.T) {
	// Given
	actual := "abcdef"
	expected := "abXdef"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "Difference at position 2") {
		t.Errorf("Expected difference at position 2, got: %s", result)
	}
	if !strings.Contains(result, "Expected character: 'X'") {
		t.Errorf("Expected character X, got: %s", result)
	}
	if !strings.Contains(result, "Actual character:   'c'") {
		t.Errorf("Expected actual character c, got: %s", result)
	}
}

func TestCharacterDiff_WithUnicodeCharacters_ShowsUnicodeInfo(t *testing.T) {
	// Given
	actual := "a🌍b"
	expected := "a🌎b"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "Difference at position 1") {
		t.Errorf("Expected difference at position 1, got: %s", result)
	}
	if !strings.Contains(result, "U+") {
		t.Errorf("Expected Unicode code point information, got: %s", result)
	}
}

func TestCharacterDiff_WithDifferentLengths_ShowsLengthDifference(t *testing.T) {
	// Given
	actual := "short"
	expected := "much longer string"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "Difference at position") {
		t.Errorf("Expected position difference for different lengths, got: %s", result)
	}
}

func TestCharacterDiff_WithOnlyLengthDifference_ShowsLengthInfo(t *testing.T) {
	// Given
	actual := "test"
	expected := "test123"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "Length difference:") {
		t.Errorf("Expected length difference message, got: %s", result)
	}
	if !strings.Contains(result, "Expected length: 7") {
		t.Errorf("Expected length 7 message, got: %s", result)
	}
	if !strings.Contains(result, "Actual length:   4") {
		t.Errorf("Expected length 4 message, got: %s", result)
	}
}

func TestCharacterDiff_WithEmptyVsNonEmpty_ShowsLengthDifference(t *testing.T) {
	// Given
	actual := ""
	expected := "not empty"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "Length difference:") {
		t.Errorf("Expected length difference message for empty vs non-empty, got: %s", result)
	}
	if !strings.Contains(result, "Expected length: 9") {
		t.Errorf("Expected length 9 message, got: %s", result)
	}
	if !strings.Contains(result, "Actual length:   0") {
		t.Errorf("Expected length 0 message, got: %s", result)
	}
}

func TestCompareStrings_WithComplexUnicodeEmoji_HandlesCorrectly(t *testing.T) {
	// Given
	actual := "Hello 👨‍👩‍👧‍👦 family!"
	expected := "Hello 👨‍👩‍👧‍👦 family!"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match for complex emoji, got: %s", result)
	}
}

func TestCompareStrings_WithDifferentComplexEmoji_ShowsDifference(t *testing.T) {
	// Given
	actual := "Hello 👨‍👩‍👧‍👦"
	expected := "Hello 👨‍👩‍👧"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed for different complex emoji, got: %s", result)
	}
}

func TestCompareStrings_WithCombiningCharacters_HandlesCorrectly(t *testing.T) {
	// Given
	actual := "café"         // é as single character
	expected := "cafe\u0301" // e + combining acute accent

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	// These should be different as they're different Unicode representations
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed for different Unicode normalization, got: %s", result)
	}
}

func TestCompareStrings_WithRightToLeftText_HandlesCorrectly(t *testing.T) {
	// Given
	actual := "Hello עברית" // Hebrew text
	expected := "Hello עברית"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match for RTL text, got: %s", result)
	}
}

func TestCompareStrings_WithControlCharacters_VisualizesCorrectly(t *testing.T) {
	// Given
	actual := "test\u0000\u0001\u0002"
	expected := "test123"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed for control characters, got: %s", result)
	}
	if !strings.Contains(result, "U+0000") {
		t.Errorf("Expected Unicode code point for null character, got: %s", result)
	}
}

func TestCompareStringsRaw_WithMatchingStrings_ReturnsMatchWithoutVisualization(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello world"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match indicator, got: %s", result)
	}
	// Should NOT contain visualized spaces (␣) but should contain raw spaces
	if strings.Contains(result, "␣") {
		t.Errorf("Expected no visualization symbols in raw mode, got: %s", result)
	}
	if !strings.Contains(result, `"hello world"¶`) {
		t.Errorf("Expected raw string with actual spaces, got: %s", result)
	}
}

func TestCompareStringsRaw_WithDifferentStrings_ReturnsAssertionFailedWithoutVisualization(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello mars"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed indicator, got: %s", result)
	}
	if !strings.Contains(result, "- Expected") {
		t.Errorf("Expected 'Expected' line with minus prefix, got: %s", result)
	}
	if !strings.Contains(result, "+ Actual") {
		t.Errorf("Expected 'Actual' line with plus prefix, got: %s", result)
	}
	// Should NOT contain visualization symbols
	if strings.Contains(result, "␣") {
		t.Errorf("Expected no space visualization in raw mode, got: %s", result)
	}
}

func TestCompareStringsRaw_WithWhitespaceCharacters_DoesNotVisualize(t *testing.T) {
	// Given
	actual := "hello\tworld\n"
	expected := "hello world "

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed, got: %s", result)
	}

	// Should NOT contain any visualization symbols
	visualizationSymbols := []string{"␣", "␉", "␊", "␍", "␋", "␌"}
	for _, symbol := range visualizationSymbols {
		if strings.Contains(result, symbol) {
			t.Errorf("Expected no visualization symbol %s in raw mode, got: %s", symbol, result)
		}
	}

	// Should contain actual whitespace characters in quotes
	if !strings.Contains(result, "\"hello\tworld\n\"¶") {
		t.Errorf("Expected raw whitespace characters, got: %s", result)
	}
}

func TestCompareStringsRaw_WithEmptyString_ShowsEmptyWithoutVisualization(t *testing.T) {
	// Given
	actual := ""
	expected := ""

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match indicator for empty strings, got: %s", result)
	}
	// Raw mode should show empty quotes, not the <empty> visualization
	if !strings.Contains(result, `""¶`) {
		t.Errorf("Expected raw empty string format, got: %s", result)
	}
	// Should NOT contain the <empty> visualization used in visualized mode
	if strings.Contains(result, "<empty>") {
		t.Errorf("Expected no <empty> visualization in raw mode, got: %s", result)
	}
}

func TestCompareStringsRaw_WithSpecialCharacters_ShowsRawFormat(t *testing.T) {
	// Given
	actual := " \t\r\n\v\f"
	expected := "normal"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed, got: %s", result)
	}

	// Should contain the raw special characters within quotes
	if !strings.Contains(result, "\" \t\r\n\v\f\"¶") {
		t.Errorf("Expected raw special characters, got: %s", result)
	}

	// Should NOT contain visualization symbols
	specialSymbols := []string{"␣", "␉", "␍", "␊", "␋", "␌"}
	for _, symbol := range specialSymbols {
		if strings.Contains(result, symbol) {
			t.Errorf("Expected no visualization symbol %s in raw mode, got: %s", symbol, result)
		}
	}
}

func TestCompareStringsRaw_WithUnicodeCharacters_HandlesCorrectlyWithoutVisualization(t *testing.T) {
	// Given
	actual := "Hello 世界! 🌍"
	expected := "Hello 世界! 🌍"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected match indicator for Unicode strings, got: %s", result)
	}
	// Should contain Unicode characters with actual spaces, not visualized spaces
	if !strings.Contains(result, "Hello 世界! 🌍") {
		t.Errorf("Expected Unicode characters with actual spaces, got: %s", result)
	}
	if strings.Contains(result, "␣") {
		t.Errorf("Expected no space visualization in raw mode, got: %s", result)
	}
}

func TestCompareStringsRaw_StillShowsCharacterDifferences(t *testing.T) {
	// Given
	actual := "abc"
	expected := "axc"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion failed, got: %s", result)
	}
	// Character difference detection should still work
	if !strings.Contains(result, "Difference at position 1") {
		t.Errorf("Expected difference position, got: %s", result)
	}
	if !strings.Contains(result, "Expected character: 'x'") {
		t.Errorf("Expected character information, got: %s", result)
	}
	if !strings.Contains(result, "Actual character:   'b'") {
		t.Errorf("Expected actual character information, got: %s", result)
	}
}
