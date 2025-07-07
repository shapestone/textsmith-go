package text

import (
	"strings"
	"testing"
)

// Edge case tests for boundary conditions and complex scenarios

func TestEdgeCase_VeryLongLines_HandlesCorrectly(t *testing.T) {
	// Given
	longString := strings.Repeat("a", 1000)
	expected := longString + "x"
	actual := longString + "y"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

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

	// Verify output contains the difference symbol
	if !strings.Contains(output, "≠") {
		t.Errorf("Expected output to contain difference symbol")
	}
	if !strings.Contains(output, "△") {
		t.Errorf("Expected output to contain pointer symbol")
	}
}

func TestEdgeCase_ManyEmptyLines_HandlesCorrectly(t *testing.T) {
	// Given
	expected := "\n\n\n\n\n"
	actual := "\n\n\n\n"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	// Should have 5 lines in expected, 4 in actual, so 5 total lines in the result
	if len(result.Lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(result.Lines))
	}

	// Last line should be missing in actual
	if result.Lines[4].Status != DiffStatusMissingInActual {
		t.Errorf("Expected last line to be DiffStatusMissingInActual, got %v", result.Lines[4].Status)
	}

	// Verify output contains left arrow for missing line
	if !strings.Contains(output, "←") {
		t.Errorf("Expected output to contain left arrow")
	}
}

func TestEdgeCase_OnlyNewlineCharacters_HandlesCorrectly(t *testing.T) {
	// Given
	expected := "\n"
	actual := ""

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if !result.HasTrailingNL {
		t.Errorf("Expected HasTrailingNL to be true")
	}

	// Should have 2 lines: empty line and missing line
	if len(result.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(result.Lines))
	}

	// Second line should be missing in actual with ␤ symbol
	if result.Lines[1].Expected != "␤" {
		t.Errorf("Expected '␤' for empty line, got '%s'", result.Lines[1].Expected)
	}
}

func TestEdgeCase_AllWhitespaceTypes_VisualizesCorrectly(t *testing.T) {
	// Given - test various whitespace characters
	expected := " \t\v\f\r"
	actual := "     "

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	// The lines stored in result contain original whitespace characters
	expectedLine := result.Lines[0].Expected
	actualLine := result.Lines[0].Actual

	// Verify original whitespace is preserved in stored lines (after line ending normalization)
	if expectedLine != " \t\v\f" { // \r is normalized away
		t.Errorf("Expected original whitespace in stored line, got '%s'", expectedLine)
	}
	if actualLine != "     " {
		t.Errorf("Expected original spaces in stored line, got '%s'", actualLine)
	}

	// Verify whitespace visualization appears in rendered output
	if !strings.Contains(output, "␣") {
		t.Errorf("Expected space visualization in rendered output")
	}
	if !strings.Contains(output, "␉") {
		t.Errorf("Expected tab visualization in rendered output")
	}
	if !strings.Contains(output, "␋") {
		t.Errorf("Expected vertical tab visualization in rendered output")
	}
	if !strings.Contains(output, "␌") {
		t.Errorf("Expected form feed visualization in rendered output")
	}
}

func TestEdgeCase_UnicodeNormalization_HandlesCorrectly(t *testing.T) {
	// Given - same visual character but different Unicode representations
	expected := "café"     // é as single character (U+00E9)
	actual := "cafe\u0301" // e + combining acute accent (U+0065 + U+0301)

	// When
	result := computeDiff(expected, actual)

	// Then
	// These should be different since they're different Unicode sequences
	if result.Match {
		t.Errorf("Expected Match to be false for different Unicode normalization")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusDifferent {
		t.Errorf("Expected DiffStatusDifferent, got %v", result.Lines[0].Status)
	}
}

func TestEdgeCase_ComplexEmoji_HandlesCorrectly(t *testing.T) {
	// Given - complex emoji with ZWJ sequences
	expected := "👨‍👩‍👧‍👦" // Family emoji
	actual := "👨‍👩‍👧"     // Family without boy

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false for different complex emoji")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusDifferent {
		t.Errorf("Expected DiffStatusDifferent, got %v", result.Lines[0].Status)
	}
}

func TestEdgeCase_RightToLeftText_HandlesCorrectly(t *testing.T) {
	// Given - Arabic text (right-to-left)
	expected := "مرحبا بالعالم"
	actual := "مرحبا بالكون"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false for different RTL text")
	}
	if len(result.Lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result.Lines))
	}
	if result.Lines[0].Status != DiffStatusDifferent {
		t.Errorf("Expected DiffStatusDifferent, got %v", result.Lines[0].Status)
	}
}

func TestEdgeCase_MixedDirectionality_HandlesCorrectly(t *testing.T) {
	// Given - mixed LTR and RTL text
	expected := "Hello مرحبا World"
	actual := "Hello مرحبا Earth"

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
}

func TestEdgeCase_ControlCharacters_HandlesCorrectly(t *testing.T) {
	// Given - various control characters
	expected := "hello\x00\x01\x02world"
	actual := "hello\x03\x04\x05world"

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

	// Control characters should be preserved in the lines
	expectedLine := result.Lines[0].Expected
	actualLine := result.Lines[0].Actual

	if !strings.Contains(expectedLine, "hello") {
		t.Errorf("Expected 'hello' in expected line")
	}
	if !strings.Contains(actualLine, "hello") {
		t.Errorf("Expected 'hello' in actual line")
	}
}

func TestEdgeCase_ZeroWidthCharacters_HandlesCorrectly(t *testing.T) {
	// Given - zero-width characters
	expected := "hello\u200Bworld" // Zero-width space
	actual := "helloworld"

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
}

func TestEdgeCase_SurrogatePairs_HandlesCorrectly(t *testing.T) {
	// Given - characters requiring surrogate pairs in UTF-16
	expected := "𝔘𝔫𝔦𝔠𝔬𝔡𝔢" // Mathematical script characters
	actual := "Unicode"

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
}

func TestEdgeCase_LargeMultilineText_HandlesCorrectly(t *testing.T) {
	// Given - large number of lines
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = "line " + string(rune('0'+i%10))
	}
	expected := strings.Join(lines, "\n")

	// Change one line in the middle
	lines[50] = "changed line"
	actual := strings.Join(lines, "\n")

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	// Should have 51 lines: 50 equal lines + 1 different line
	if len(result.Lines) != 51 {
		t.Errorf("Expected 51 lines, got %d", len(result.Lines))
	}

	// First 50 lines should be equal
	for i := 0; i < 50; i++ {
		if result.Lines[i].Status != DiffStatusEqual {
			t.Errorf("Expected line %d to be DiffStatusEqual, got %v", i, result.Lines[i].Status)
		}
	}

	// 51st line should be different
	if result.Lines[50].Status != DiffStatusDifferent {
		t.Errorf("Expected line 50 to be DiffStatusDifferent, got %v", result.Lines[50].Status)
	}

	// Output should stop at the first difference
	outputLines := strings.Split(output, "\n")
	// Should have header (2 lines) + 50 equal lines + 1 different line + 1 pointer line = 54 lines
	if len(outputLines) != 54 {
		t.Errorf("Expected 54 output lines, got %d", len(outputLines))
	}
}

func TestEdgeCase_WidthCalculationWithUnicode_HandlesCorrectly(t *testing.T) {
	// Given - strings with different visual widths
	expected := "A" // 1 column
	actual := "🌍"   // 2 columns (wide character)

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	// Width should be calculated based on rune count, not visual width
	minWidth := len("Expected")
	if result.ExpectedWidth < minWidth {
		t.Errorf("Expected ExpectedWidth to be at least %d, got %d", minWidth, result.ExpectedWidth)
	}
	if result.ActualWidth < minWidth {
		t.Errorf("Expected ActualWidth to be at least %d, got %d", minWidth, result.ActualWidth)
	}
}

func TestEdgeCase_EmptyVsNewlineOnly_HandlesCorrectly(t *testing.T) {
	// Given
	expected := ""
	actual := "\n"

	// When
	result := computeDiff(expected, actual)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if len(result.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(result.Lines))
	}

	// First line should be equal (both empty)
	if result.Lines[0].Status != DiffStatusEqual {
		t.Errorf("Expected first line to be DiffStatusEqual, got %v", result.Lines[0].Status)
	}

	// Second line should be missing in expected
	if result.Lines[1].Status != DiffStatusMissingInExpected {
		t.Errorf("Expected second line to be DiffStatusMissingInExpected, got %v", result.Lines[1].Status)
	}
}

// FIXED: This test had an incorrect assumption about line ending normalization
func TestEdgeCase_LineEndingNormalization_WorksCorrectly(t *testing.T) {
	// Given - test that normalization works for equivalent content
	expected := "line1\r\nline2\nline3"
	actual := "line1\nline2\nline3"

	// When
	result := computeDiff(expected, actual)

	// Then
	// After normalization, both should have the same content
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

// NEW: Test consecutive line endings that create empty lines
func TestEdgeCase_ConsecutiveLineEndingsCreateEmptyLines_HandlesCorrectly(t *testing.T) {
	// Given - consecutive line endings that create empty lines
	expected := "line1\n\r\nline2\r\n\nline3" // Creates empty lines between content
	actual := "line1\n\nline2\n\nline3"       // Equivalent after normalization

	// When
	result := computeDiff(expected, actual)

	// Then
	// After normalization, both should have the same content (including empty lines)
	if !result.Match {
		t.Errorf("Expected Match to be true after line ending normalization")
	}
	if len(result.Lines) != 5 {
		t.Errorf("Expected 5 lines, got %d", len(result.Lines))
	}

	expectedLines := []string{"line1", "", "line2", "", "line3"}
	for i, line := range result.Lines {
		if line.Status != DiffStatusEqual {
			t.Errorf("Expected line %d to be DiffStatusEqual, got %v", i, line.Status)
		}
		if line.Expected != expectedLines[i] {
			t.Errorf("Expected line %d to be '%s', got '%s'", i, expectedLines[i], line.Expected)
		}
	}
}

func TestEdgeCase_StringDiffFunction_HandlesEdgeCases(t *testing.T) {
	// Given - test the internal stringDiff function with edge cases
	testCases := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{
			name:     "identical strings",
			a:        "hello",
			b:        "hello",
			expected: "     \u25B3",
		},
		{
			name:     "completely different",
			a:        "abc",
			b:        "xyz",
			expected: "\u25B3",
		},
		{
			name:     "different lengths",
			a:        "short",
			b:        "sh",
			expected: "  \u25B3",
		},
		{
			name:     "empty strings",
			a:        "",
			b:        "",
			expected: "\u25B3",
		},
		{
			name:     "one empty",
			a:        "test",
			b:        "",
			expected: "\u25B3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := stringDiff(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestEdgeCase_HasTrailingNewlineFunction_HandlesEdgeCases(t *testing.T) {
	// Given - test the internal hasTrailingNewline function
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "unix newline",
			input:    "text\n",
			expected: true,
		},
		{
			name:     "windows newline",
			input:    "text\r\n",
			expected: true,
		},
		{
			name:     "mac newline",
			input:    "text\r",
			expected: true,
		},
		{
			name:     "no newline",
			input:    "text",
			expected: false,
		},
		{
			name:     "only newline",
			input:    "\n",
			expected: true,
		},
		{
			name:     "multiple newlines",
			input:    "text\n\n",
			expected: true,
		},
		{
			name:     "mixed line endings",
			input:    "text\r\n\r",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := hasTrailingNewline(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
