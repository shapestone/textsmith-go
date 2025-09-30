package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"strings"
	"testing"
)

// TestCompareStrings_WithMatchingStrings_ReturnsMatch tests identical strings
func TestCompareStrings_WithMatchingStrings_ReturnsMatch(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello world"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected result to contain '✓ [MATCH]', got: %s", result)
	}

	if !strings.Contains(result, "hello␣world") {
		t.Errorf("Expected result to contain visualized string 'hello␣world', got: %s", result)
	}
}

// TestCompareStrings_WithDifferentStrings_ReturnsAssertionFailed tests different strings
func TestCompareStrings_WithDifferentStrings_ReturnsAssertionFailed(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello mars"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "hello␣mars") {
		t.Errorf("Expected result to contain visualized expected 'hello␣mars', got: %s", result)
	}

	if !strings.Contains(result, "hello␣world") {
		t.Errorf("Expected result to contain visualized actual 'hello␣world', got: %s", result)
	}

	if !strings.Contains(result, "Difference at position 6") {
		t.Errorf("Expected result to show difference at position 6, got: %s", result)
	}

	if !strings.Contains(result, "'m' (U+006D)") {
		t.Errorf("Expected result to show expected character 'm' with Unicode, got: %s", result)
	}

	if !strings.Contains(result, "'w' (U+0077)") {
		t.Errorf("Expected result to show actual character 'w' with Unicode, got: %s", result)
	}
}

// TestCompareStrings_WithWhitespaceDifference_ShowsVisualSymbols tests whitespace detection
func TestCompareStrings_WithWhitespaceDifference_ShowsVisualSymbols(t *testing.T) {
	// Given
	actual := "hello\tworld"
	expected := "hello world"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "␉") {
		t.Errorf("Expected result to contain tab symbol '␉', got: %s", result)
	}

	if !strings.Contains(result, "␣") {
		t.Errorf("Expected result to contain space symbol '␣', got: %s", result)
	}

	if !strings.Contains(result, "'\\t' (U+0009)") {
		t.Errorf("Expected result to show tab character with Unicode, got: %s", result)
	}

	if !strings.Contains(result, "' ' (U+0020)") {
		t.Errorf("Expected result to show space character with Unicode, got: %s", result)
	}
}

// TestCompareStrings_WithEmptyStrings_ShowsEmptyIndicator tests empty string handling
func TestCompareStrings_WithEmptyStrings_ShowsEmptyIndicator(t *testing.T) {
	// Given
	actual := ""
	expected := ""

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected result to contain '✓ [MATCH]', got: %s", result)
	}

	if !strings.Contains(result, "<empty>") {
		t.Errorf("Expected result to contain '<empty>' indicator, got: %s", result)
	}
}

// TestCompareStrings_WithEmptyVsNonEmpty_ShowsDifference tests empty vs non-empty
func TestCompareStrings_WithEmptyVsNonEmpty_ShowsDifference(t *testing.T) {
	// Given
	actual := "hello"
	expected := ""

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "<empty>") {
		t.Errorf("Expected result to contain '<empty>' for expected, got: %s", result)
	}

	if !strings.Contains(result, "hello") {
		t.Errorf("Expected result to contain 'hello' for actual, got: %s", result)
	}

	if !strings.Contains(result, "Difference at position 0") {
		t.Errorf("Expected result to show difference at position 0, got: %s", result)
	}
}

// TestCompareStrings_WithUnicodeContent_HandlesCorrectly tests Unicode support
func TestCompareStrings_WithUnicodeContent_HandlesCorrectly(t *testing.T) {
	// Given
	actual := "Hello 世界! 🌍"
	expected := "Hello 世界! 🌎"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "世界") {
		t.Errorf("Expected result to contain Unicode characters '世界', got: %s", result)
	}

	if !strings.Contains(result, "Difference at position") {
		t.Errorf("Expected result to show difference position, got: %s", result)
	}
}

// TestCompareStrings_WithNewlines_VisualizesCorrectly tests newline visualization
func TestCompareStrings_WithNewlines_VisualizesCorrectly(t *testing.T) {
	// Given
	actual := "line1\nline2"
	expected := "line1\rline2"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "␊") {
		t.Errorf("Expected result to contain newline symbol '␊', got: %s", result)
	}

	if !strings.Contains(result, "␍") {
		t.Errorf("Expected result to contain carriage return symbol '␍', got: %s", result)
	}
}

// TestCompareStrings_WithMultipleWhitespaceTypes_VisualizesAll tests various whitespace
func TestCompareStrings_WithMultipleWhitespaceTypes_VisualizesAll(t *testing.T) {
	// Given
	actual := "test \t\n\r\v\f content"
	expected := "test content"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	// Check for various whitespace symbols
	if !strings.Contains(result, "␣") {
		t.Errorf("Expected result to contain space symbol '␣', got: %s", result)
	}

	if !strings.Contains(result, "␉") {
		t.Errorf("Expected result to contain tab symbol '␉', got: %s", result)
	}

	if !strings.Contains(result, "␊") {
		t.Errorf("Expected result to contain newline symbol '␊', got: %s", result)
	}

	if !strings.Contains(result, "␍") {
		t.Errorf("Expected result to contain carriage return symbol '␍', got: %s", result)
	}

	if !strings.Contains(result, "␋") {
		t.Errorf("Expected result to contain vertical tab symbol '␋', got: %s", result)
	}

	if !strings.Contains(result, "␌") {
		t.Errorf("Expected result to contain form feed symbol '␌', got: %s", result)
	}
}

// TestCompareStrings_WithDifferentLengths_ShowsEndOfString tests length differences
func TestCompareStrings_WithDifferentLengths_ShowsEndOfString(t *testing.T) {
	// Given
	actual := "hello"
	expected := "hello world"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "<end of string>") {
		t.Errorf("Expected result to show '<end of string>', got: %s", result)
	}

	if !strings.Contains(result, "Difference at position 5") {
		t.Errorf("Expected result to show difference at position 5, got: %s", result)
	}
}

// TestCompareStringsRaw_WithMatchingStrings_ReturnsMatch tests raw comparison with match
func TestCompareStringsRaw_WithMatchingStrings_ReturnsMatch(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello world"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected result to contain '✓ [MATCH]', got: %s", result)
	}

	// Should NOT contain visualization symbols
	if strings.Contains(result, "␣") {
		t.Errorf("Expected raw result to NOT contain space symbol '␣', got: %s", result)
	}

	// Should contain the actual space
	if !strings.Contains(result, "hello world") {
		t.Errorf("Expected result to contain raw string 'hello world', got: %s", result)
	}
}

// TestCompareStringsRaw_WithDifferentStrings_ShowsRawContent tests raw comparison without visualization
func TestCompareStringsRaw_WithDifferentStrings_ShowsRawContent(t *testing.T) {
	// Given
	actual := "hello world"
	expected := "hello mars"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	// Should NOT contain visualization symbols
	if strings.Contains(result, "␣") {
		t.Errorf("Expected raw result to NOT contain space symbol '␣', got: %s", result)
	}

	if !strings.Contains(result, "hello mars") {
		t.Errorf("Expected result to contain raw expected 'hello mars', got: %s", result)
	}

	if !strings.Contains(result, "hello world") {
		t.Errorf("Expected result to contain raw actual 'hello world', got: %s", result)
	}

	if !strings.Contains(result, "Difference at position 6") {
		t.Errorf("Expected result to show difference at position 6, got: %s", result)
	}
}

// TestCompareStringsRaw_WithWhitespace_ShowsRawWhitespace tests raw whitespace handling
func TestCompareStringsRaw_WithWhitespace_ShowsRawWhitespace(t *testing.T) {
	// Given
	actual := "hello\tworld"
	expected := "hello world"

	// When
	result := text.CompareStringsRaw(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	// Should NOT contain visualization symbols for space
	if strings.Contains(result, "␣") {
		t.Errorf("Expected raw result to NOT contain space symbol '␣', got: %s", result)
	}

	// Should NOT contain visualization symbols for tab in displayed string
	if strings.Contains(result, "␉") {
		t.Errorf("Expected raw result to NOT contain tab symbol '␉' in displayed string, got: %s", result)
	}

	// But should show the character codes in the difference section
	if !strings.Contains(result, "'\\t' (U+0009)") {
		t.Errorf("Expected result to show tab character code, got: %s", result)
	}

	if !strings.Contains(result, "' ' (U+0020)") {
		t.Errorf("Expected result to show space character code, got: %s", result)
	}
}

// TestCompareStrings_WithComplexEmoji_HandlesCorrectly tests multi-codepoint emoji
func TestCompareStrings_WithComplexEmoji_HandlesCorrectly(t *testing.T) {
	// Given (family emoji is multiple codepoints)
	actual := "Hello 👨‍👩‍👧‍👦"
	expected := "Hello 👨‍👩‍👧"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "Difference at position") {
		t.Errorf("Expected result to show difference position, got: %s", result)
	}
}

// TestCompareStrings_WithCJKCharacters_HandlesCorrectly tests CJK character handling
func TestCompareStrings_WithCJKCharacters_HandlesCorrectly(t *testing.T) {
	// Given
	actual := "你好世界"
	expected := "你好世界！"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to contain '✗ [ASSERTION_FAILED]', got: %s", result)
	}

	if !strings.Contains(result, "你好世界") {
		t.Errorf("Expected result to contain CJK characters, got: %s", result)
	}

	if !strings.Contains(result, "Difference at position 4") {
		t.Errorf("Expected result to show difference at position 4, got: %s", result)
	}
}

// TestCompareStrings_WithEndMarker_ShowsEndMarker tests the ¶ end marker
func TestCompareStrings_WithEndMarker_ShowsEndMarker(t *testing.T) {
	// Given
	actual := "test"
	expected := "test"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "¶") {
		t.Errorf("Expected result to contain end marker '¶', got: %s", result)
	}
}

// TestCompareStrings_WithExpectedPrefixSuffix_ShowsCorrectFormat tests output format
func TestCompareStrings_WithExpectedPrefixSuffix_ShowsCorrectFormat(t *testing.T) {
	// Given
	actual := "hello"
	expected := "world"

	// When
	result := text.CompareStrings(actual, expected)

	// Then
	if !strings.Contains(result, "- Expected:") {
		t.Errorf("Expected result to contain '- Expected:', got: %s", result)
	}

	if !strings.Contains(result, "+ Actual:") {
		t.Errorf("Expected result to contain '+ Actual:', got: %s", result)
	}
}
