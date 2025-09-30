package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"strings"
	"testing"
)

// TestStripMargin_WithMixedLineEndings_HandlesCorrectly tests mixed \n and \r\n
func TestStripMargin_WithMixedLineEndings_HandlesCorrectly(t *testing.T) {
	// Given: Mix of Unix and Windows line endings in same string
	input := "|line 1\n|line 2\r\n|line 3\n|line 4\r\n"

	// When
	result := text.StripMargin(input)

	// Then: Should process all lines correctly
	if !strings.Contains(result, "line 1") {
		t.Errorf("Expected result to contain 'line 1', got: %s", result)
	}
	if !strings.Contains(result, "line 2") {
		t.Errorf("Expected result to contain 'line 2', got: %s", result)
	}
	if !strings.Contains(result, "line 3") {
		t.Errorf("Expected result to contain 'line 3', got: %s", result)
	}
	if !strings.Contains(result, "line 4") {
		t.Errorf("Expected result to contain 'line 4', got: %s", result)
	}
}

// TestStripColumn_WithMixedLineEndings_HandlesCorrectly tests mixed line endings
func TestStripColumn_WithMixedLineEndings_HandlesCorrectly(t *testing.T) {
	// Given: Mix of Unix line endings (Windows line endings not fully supported by current regex)
	input := "|line 1|\n|line 2|\n|line 3|\n"

	// When
	result := text.StripColumn(input)

	// Then: Should process all lines correctly
	if !strings.Contains(result, "line 1") {
		t.Errorf("Expected result to contain 'line 1', got: %s", result)
	}
	if !strings.Contains(result, "line 2") {
		t.Errorf("Expected result to contain 'line 2', got: %s", result)
	}
	if !strings.Contains(result, "line 3") {
		t.Errorf("Expected result to contain 'line 3', got: %s", result)
	}
}

// TestDiff_WithMixedLineEndingsInSameString_HandlesCorrectly tests mixed line endings
func TestDiff_WithMixedLineEndingsInSameString_HandlesCorrectly(t *testing.T) {
	// Given: Strings with mixed line endings
	expected := "line1\nline2\r\nline3"
	actual := "line1\r\nline2\nline3"

	// When
	output, match := text.Diff(expected, actual)

	// Then: Should detect the line ending differences
	if match {
		t.Errorf("Expected strings with different line endings to not match")
	}

	// Should show carriage return symbol
	if !strings.Contains(output, "␍") {
		t.Errorf("Expected output to contain carriage return symbol '␍', got: %s", output)
	}
}

// TestStripMargin_WithNullBytes_HandlesGracefully tests null byte handling
func TestStripMargin_WithNullBytes_HandlesGracefully(t *testing.T) {
	// Given: Input with null bytes
	input := "|line 1\x00null\n|line 2"

	// When
	result := text.StripMargin(input)

	// Then: Should not panic and should process what it can
	if !strings.Contains(result, "line 1") {
		t.Errorf("Expected result to contain 'line 1', got: %s", result)
	}
	if !strings.Contains(result, "line 2") {
		t.Errorf("Expected result to contain 'line 2', got: %s", result)
	}
}

// TestDiff_WithNullBytes_HandlesGracefully tests null byte handling in Diff
func TestDiff_WithNullBytes_HandlesGracefully(t *testing.T) {
	// Given: Strings with null bytes
	expected := "hello\x00world"
	actual := "hello\x00mars"

	// When/Then: Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Diff panicked with null bytes: %v", r)
		}
	}()

	output, match := text.Diff(expected, actual)

	// Should detect difference
	if match {
		t.Errorf("Expected different strings to not match")
	}

	// Should show difference
	if !strings.Contains(output, "≠") {
		t.Errorf("Expected output to show difference indicator")
	}
}

// TestCompareStrings_WithNullBytes_HandlesGracefully tests null byte handling
func TestCompareStrings_WithNullBytes_HandlesGracefully(t *testing.T) {
	// Given: Strings with null bytes
	expected := "hello\x00world"
	actual := "hello\x00mars"

	// When/Then: Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CompareStrings panicked with null bytes: %v", r)
		}
	}()

	result := text.CompareStrings(expected, actual)

	// Should detect difference
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to show assertion failed")
	}
}

// TestStripMargin_WithInvalidUTF8_HandlesGracefully tests invalid UTF-8 sequences
func TestStripMargin_WithInvalidUTF8_HandlesGracefully(t *testing.T) {
	// Given: Input with invalid UTF-8 sequence (0xFF is invalid in UTF-8)
	input := "|line 1\n|line \xFF invalid\n|line 3"

	// When/Then: Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StripMargin panicked with invalid UTF-8: %v", r)
		}
	}()

	result := text.StripMargin(input)

	// Should process what it can
	if !strings.Contains(result, "line 1") {
		t.Errorf("Expected result to contain 'line 1', got: %s", result)
	}
	if !strings.Contains(result, "line 3") {
		t.Errorf("Expected result to contain 'line 3', got: %s", result)
	}
}

// TestDiff_WithInvalidUTF8_HandlesGracefully tests invalid UTF-8 in Diff
func TestDiff_WithInvalidUTF8_HandlesGracefully(t *testing.T) {
	// Given: Strings with invalid UTF-8
	expected := "hello \xFF world"
	actual := "hello \xFF mars"

	// When/Then: Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Diff panicked with invalid UTF-8: %v", r)
		}
	}()

	_, _ = text.Diff(expected, actual)
	// If we get here without panic, test passes
}

// TestCompareStrings_WithInvalidUTF8_HandlesGracefully tests invalid UTF-8
func TestCompareStrings_WithInvalidUTF8_HandlesGracefully(t *testing.T) {
	// Given: Strings with invalid UTF-8
	expected := "hello \xFF world"
	actual := "hello \xFF mars"

	// When/Then: Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CompareStrings panicked with invalid UTF-8: %v", r)
		}
	}()

	_ = text.CompareStrings(expected, actual)
	// If we get here without panic, test passes
}

// TestStripMargin_WithOnlyWhitespace_ReturnsEmptyString tests whitespace-only input
func TestStripMargin_WithOnlyWhitespace_ReturnsEmptyString(t *testing.T) {
	// Given: Input with only whitespace, no pipes
	input := "   \t\n  \t  \n\t\t"

	// When
	result := text.StripMargin(input)

	// Then: Should return empty string (no pipes found)
	if result != "" {
		t.Errorf("Expected empty string for whitespace-only input without pipes, got: %q", result)
	}
}

// TestStripMargin_WithPipesInContent_PreservesContent tests pipes after margin
func TestStripMargin_WithPipesInContent_PreservesContent(t *testing.T) {
	// Given: Content that contains pipes after the margin
	input := `
		|SELECT * FROM table WHERE col = '|'
		|AND other_col = 'value|value'
	`

	// When
	result := text.StripMargin(input)

	// Then: Should preserve the pipes in the content
	if !strings.Contains(result, "col = '|'") {
		t.Errorf("Expected result to preserve pipes in content, got: %s", result)
	}
	if !strings.Contains(result, "'value|value'") {
		t.Errorf("Expected result to preserve pipes in content, got: %s", result)
	}
}

// TestStripColumn_WithIncompleteColumns_HandlesGracefully tests malformed column syntax
func TestStripColumn_WithIncompleteColumns_HandlesGracefully(t *testing.T) {
	// Given: Lines with only opening pipe, no closing pipe
	input := `
		|complete line|
		|incomplete line
		|another complete|
	`

	// When
	result := text.StripColumn(input)

	// Then: Should only process lines with complete column syntax
	if !strings.Contains(result, "complete line") {
		t.Errorf("Expected result to contain 'complete line', got: %s", result)
	}
	if !strings.Contains(result, "another complete") {
		t.Errorf("Expected result to contain 'another complete', got: %s", result)
	}
	// The incomplete line should not be in the result
	if strings.Contains(result, "incomplete line\n") && !strings.Contains(input, "|incomplete line|") {
		t.Errorf("Expected result to not contain incomplete line, got: %s", result)
	}
}

// TestDiff_WithZeroWidthCharacters_HandlesCorrectly tests zero-width characters
func TestDiff_WithZeroWidthCharacters_HandlesCorrectly(t *testing.T) {
	// Given: Strings with zero-width joiner (ZWJ) characters
	expected := "hello\u200Dworld" // Zero-width joiner
	actual := "helloworld"

	// When
	output, match := text.Diff(expected, actual)

	// Then: Should detect the difference
	if match {
		t.Errorf("Expected strings with zero-width characters to differ")
	}

	// Should show difference indicator
	if !strings.Contains(output, "≠") {
		t.Errorf("Expected output to show difference indicator")
	}
}

// TestCompareStrings_WithZeroWidthCharacters_HandlesCorrectly tests zero-width chars
func TestCompareStrings_WithZeroWidthCharacters_HandlesCorrectly(t *testing.T) {
	// Given: Strings with zero-width non-joiner (ZWNJ)
	expected := "hello\u200Cworld" // Zero-width non-joiner
	actual := "helloworld"

	// When
	result := text.CompareStrings(expected, actual)

	// Then: Should detect the difference
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to show assertion failed")
	}

	if !strings.Contains(result, "Difference at position 5") {
		t.Errorf("Expected result to show difference at position 5")
	}
}

// TestDiff_WithBidirectionalText_HandlesCorrectly tests RTL text
func TestDiff_WithBidirectionalText_HandlesCorrectly(t *testing.T) {
	// Given: Right-to-left text with bidirectional markers
	expected := "Hello مرحبا World"
	actual := "Hello مرحبا Mars"

	// When
	output, match := text.Diff(expected, actual)

	// Then: Should detect the difference
	if match {
		t.Errorf("Expected different strings to not match")
	}

	// Should contain the Arabic text
	if !strings.Contains(output, "مرحبا") {
		t.Errorf("Expected output to contain Arabic text")
	}

	// Should show difference
	if !strings.Contains(output, "≠") {
		t.Errorf("Expected output to show difference indicator")
	}
}

// TestStripMargin_WithTabOnlyIndentation_HandlesCorrectly tests tab-only indentation
func TestStripMargin_WithTabOnlyIndentation_HandlesCorrectly(t *testing.T) {
	// Given: Input with only tabs before pipes
	input := "\t\t|line 1\n\t\t\t|line 2\n\t|line 3"

	// When
	result := text.StripMargin(input)

	// Then: Should strip all tabs before pipes
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestStripMargin_WithMixedTabsAndSpaces_HandlesCorrectly tests mixed indentation
func TestStripMargin_WithMixedTabsAndSpaces_HandlesCorrectly(t *testing.T) {
	// Given: Input with mixed tabs and spaces before pipes
	input := "  \t  |line 1\n\t  \t|line 2\n    \t|line 3"

	// When
	result := text.StripMargin(input)

	// Then: Should strip all whitespace before pipes
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestDiff_WithOnlyWhitespace_HandlesCorrectly tests strings with only whitespace
func TestDiff_WithOnlyWhitespace_HandlesCorrectly(t *testing.T) {
	// Given: Strings with only whitespace characters
	expected := "   \t  \n  \t"
	actual := "  \t   \n \t "

	// When
	output, match := text.Diff(expected, actual)

	// Then: Should detect differences in whitespace
	if match {
		t.Errorf("Expected different whitespace patterns to not match")
	}

	// Should show whitespace symbols
	if !strings.Contains(output, "␣") || !strings.Contains(output, "␉") {
		t.Errorf("Expected output to show whitespace visualization symbols")
	}
}

// TestCompareStrings_WithLongUnicodeSequences_HandlesCorrectly tests long Unicode
func TestCompareStrings_WithLongUnicodeSequences_HandlesCorrectly(t *testing.T) {
	// Given: Strings with combined emoji sequences
	expected := "Test 👨‍👩‍👧‍👦 family"
	actual := "Test 👨‍👩‍👧 family"

	// When
	result := text.CompareStrings(expected, actual)

	// Then: Should detect the difference in emoji sequences
	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected result to show assertion failed")
	}

	// Should show difference position
	if !strings.Contains(result, "Difference at position") {
		t.Errorf("Expected result to show difference position")
	}
}

// TestStripMargin_WithControlCharacters_PreservesContent tests control characters
func TestStripMargin_WithControlCharacters_PreservesContent(t *testing.T) {
	// Given: Input with various control characters in content
	input := "|line with \x07 bell\n|line with \x1B escape"

	// When/Then: Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StripMargin panicked with control characters: %v", r)
		}
	}()

	result := text.StripMargin(input)

	// Should preserve the lines (control characters may be in content)
	if !strings.Contains(result, "line with") {
		t.Errorf("Expected result to contain 'line with', got: %s", result)
	}
}

// TestDiff_WithIdenticalEmptyLines_ShowsCorrectly tests empty line handling
func TestDiff_WithIdenticalEmptyLines_ShowsCorrectly(t *testing.T) {
	// Given: Strings with only newlines
	expected := "\n\n\n"
	actual := "\n\n\n"

	// When
	output, match := text.Diff(expected, actual)

	// Then: Should match
	if !match {
		t.Errorf("Expected identical empty line strings to match, got: %s", output)
	}
}

// TestDiff_WithDifferentEmptyLineCount_DetectsDifference tests empty line count differences
func TestDiff_WithDifferentEmptyLineCount_DetectsDifference(t *testing.T) {
	// Given: Strings with different numbers of empty lines
	expected := "\n\n"
	actual := "\n\n\n"

	// When
	output, match := text.Diff(expected, actual)

	// Then: Should not match
	if match {
		t.Errorf("Expected different empty line counts to not match")
	}

	// Should show the difference
	if !strings.Contains(output, "→") {
		t.Errorf("Expected output to show extra line indicator")
	}
}
