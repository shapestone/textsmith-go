package text

import (
	"testing"
)

// Tests for normalizeWhitespace function (per changes.md specification)

func TestNormalizeWhitespace_WithLeadingWhitespace_TrimsIt(t *testing.T) {
	// Given
	input := "   hello"

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithTrailingWhitespace_TrimsIt(t *testing.T) {
	// Given
	input := "hello   "

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithLeadingAndTrailingWhitespace_TrimsBoth(t *testing.T) {
	// Given
	input := "  hello world  "

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithMultipleInternalSpaces_CollapsesToSingle(t *testing.T) {
	// Given
	input := "hello    world"

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithTabsInside_CollapsesToSingleSpace(t *testing.T) {
	// Given
	input := "hello\t\tworld"

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithMixedWhitespace_CollapsesToSingleSpace(t *testing.T) {
	// Given
	input := "hello \t  \t world"

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithWhitespaceOnly_ReturnsEmpty(t *testing.T) {
	// Given
	input := "   \t  \t  "

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := ""
	if result != expected {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestNormalizeWhitespace_WithEmptyString_ReturnsEmpty(t *testing.T) {
	// Given
	input := ""

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := ""
	if result != expected {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestNormalizeWhitespace_WithNoWhitespace_ReturnsUnchanged(t *testing.T) {
	// Given
	input := "hello"

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestNormalizeWhitespace_WithSingleSpace_ReturnsUnchanged(t *testing.T) {
	// Given
	input := "hello world"

	// When
	result := normalizeWhitespace(input)

	// Then
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Tests for Diff with ignoreWhitespace parameter

func TestDiff_WithIgnoreWhitespaceTrue_AndIdenticalContent_ReturnsMatch(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello world"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected match for identical content with ignoreWhitespace=true")
	}
}

func TestDiff_WithIgnoreWhitespaceTrue_AndDifferentWhitespace_ReturnsMatch(t *testing.T) {
	// Given
	expected := "hello    world"
	actual := "hello world"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected match when only whitespace differs with ignoreWhitespace=true")
	}
}

func TestDiff_WithIgnoreWhitespaceFalse_AndDifferentWhitespace_ReturnsNoMatch(t *testing.T) {
	// Given
	expected := "hello    world"
	actual := "hello world"

	// When
	_, match := Diff(expected, actual, false)

	// Then
	if match {
		t.Error("Expected no match when whitespace differs with ignoreWhitespace=false")
	}
}

func TestDiff_WithIgnoreWhitespaceTrue_AndDifferentContent_ReturnsNoMatch(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello universe"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if match {
		t.Error("Expected no match when actual content differs even with ignoreWhitespace=true")
	}
}

func TestDiff_WithIgnoreWhitespaceTrue_PreservesOriginalWhitespaceInOutput(t *testing.T) {
	// Given
	expected := "hello    world"
	actual := "hello world"

	// When
	output, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected match with ignoreWhitespace=true")
	}

	// Output should show original whitespace (visualized)
	// "hello    world" should be shown as "hello␣␣␣␣world"
	if !contains(output, "hello␣␣␣␣world") {
		t.Errorf("Expected output to preserve original whitespace visualization, got: %s", output)
	}
}

func TestDiff_WithIgnoreWhitespaceTrue_AndLeadingTrailingWhitespace_ReturnsMatch(t *testing.T) {
	// Given
	expected := "  hello world  "
	actual := "hello world"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected match when only leading/trailing whitespace differs with ignoreWhitespace=true")
	}
}

func TestDiff_WithIgnoreWhitespaceTrue_AndTabsVsSpaces_ReturnsMatch(t *testing.T) {
	// Given
	expected := "hello\tworld"
	actual := "hello world"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected match when tabs vs spaces with ignoreWhitespace=true")
	}
}

// Helper function for substring checking
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ========================================
// Comprehensive Test Coverage per test-spec.md
// ========================================

// Dimension 1: Whitespace Comparison Modes

func TestDiff_Dimension1_ExactMode_WithDifferentSpaces_CausesMismatch(t *testing.T) {
	// Given
	expected := "hello  world" // 2 spaces
	actual := "hello world"    // 1 space

	// When
	_, match := Diff(expected, actual, false)

	// Then
	if match {
		t.Error("Expected exact mode to detect different number of spaces")
	}
}

func TestDiff_Dimension1_IgnoreMode_WithDifferentSpaces_StillMatches(t *testing.T) {
	// Given
	expected := "hello  world" // 2 spaces
	actual := "hello world"    // 1 space

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected ignore mode to match despite different spaces")
	}
}

func TestDiff_Dimension1_BoundaryCase_EmptyVsWhitespaceOnlyString(t *testing.T) {
	// Given
	expected := ""
	actual := "   \t  "

	// When - ignore mode
	_, matchIgnore := Diff(expected, actual, true)

	// Then - should match because whitespace-only normalizes to empty
	if !matchIgnore {
		t.Error("Expected ignore mode to match empty vs whitespace-only")
	}

	// When - exact mode
	_, matchExact := Diff(expected, actual, false)

	// Then - should not match in exact mode
	if matchExact {
		t.Error("Expected exact mode to not match empty vs whitespace-only")
	}
}

func TestDiff_Dimension1_MixedWhitespaceTypes(t *testing.T) {
	// Given - various whitespace types
	expected := "a \t\v\f b"
	actual := "a b"

	// When
	_, match := Diff(expected, actual, true)

	// Then - should match with ignore mode
	if !match {
		t.Error("Expected ignore mode to match mixed whitespace types")
	}
}

// Dimension 2: Line-by-Line Comparison

func TestDiff_Dimension2_DifferenceInFirstLine(t *testing.T) {
	// Given
	expected := "different\nline2\nline3"
	actual := "changed\nline2\nline3"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if match {
		t.Error("Expected difference in first line to be detected")
	}
}

func TestDiff_Dimension2_DifferenceInMiddleLine(t *testing.T) {
	// Given
	expected := "line1\ndifferent\nline3"
	actual := "line1\nchanged\nline3"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if match {
		t.Error("Expected difference in middle line to be detected")
	}
}

func TestDiff_Dimension2_DifferenceInLastLine(t *testing.T) {
	// Given
	expected := "line1\nline2\ndifferent"
	actual := "line1\nline2\nchanged"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if match {
		t.Error("Expected difference in last line to be detected")
	}
}

func TestDiff_Dimension2_DifferentLineCountsWithWhitespace(t *testing.T) {
	// Given
	expected := "line1   \nline2\nline3"
	actual := "line1\nline2"

	// When
	_, match := Diff(expected, actual, true)

	// Then - different line counts should not match even in ignore mode
	if match {
		t.Error("Expected different line counts to not match")
	}
}

func TestDiff_Dimension2_VeryLongLineWithWhitespaceDifference(t *testing.T) {
	// Given - line with 1000+ characters
	prefix := make([]rune, 500)
	for i := range prefix {
		prefix[i] = 'a'
	}
	suffix := make([]rune, 500)
	for i := range suffix {
		suffix[i] = 'b'
	}

	expected := string(prefix) + "    " + string(suffix) // 4 spaces
	actual := string(prefix) + " " + string(suffix)      // 1 space

	// When
	_, match := Diff(expected, actual, true)

	// Then - should match with ignore mode
	if !match {
		t.Error("Expected long lines with different whitespace to match in ignore mode")
	}
}

// Dimension 3: Whitespace Visualization

func TestDiff_Dimension3_PreservesOriginalInIgnoreMode(t *testing.T) {
	// Given
	expected := "hello    world" // 4 spaces
	actual := "hello world"      // 1 space

	// When
	output, match := Diff(expected, actual, true)

	// Then - should match
	if !match {
		t.Error("Expected match in ignore mode")
	}

	// Output should show original whitespace (4 spaces as ␣␣␣␣)
	if !contains(output, "hello␣␣␣␣world") {
		t.Errorf("Expected output to preserve original 4 spaces visualization, got: %s", output)
	}
}

func TestDiff_Dimension3_AllWhitespaceSymbolsRender(t *testing.T) {
	// Given - string with all whitespace types (Note: \r is normalized to \n)
	expected := " \t\v\f"
	actual := " \t\v\f"

	// When
	output, _ := Diff(expected, actual, false)

	// Then - should show all symbols
	if !contains(output, "␣") { // space
		t.Error("Expected space symbol ␣ in output")
	}
	if !contains(output, "␉") { // tab
		t.Error("Expected tab symbol ␉ in output")
	}
	if !contains(output, "␋") { // vtab
		t.Error("Expected vtab symbol ␋ in output")
	}
	if !contains(output, "␌") { // form feed
		t.Error("Expected form feed symbol ␌ in output")
	}
}

// Dimension 4: Cross-Platform Line Endings (now implemented!)

func TestDiff_Dimension4_UnixVsWindowsLineEndings_BothModesNormalize(t *testing.T) {
	// Given
	unix := "line1\nline2"
	windows := "line1\r\nline2"

	// When - exact mode
	_, matchExact := Diff(unix, windows, false)

	// When - ignore mode
	_, matchIgnore := Diff(unix, windows, true)

	// Then - should match in both modes due to line ending normalization
	if !matchExact {
		t.Error("Expected Unix and Windows line endings to match in exact mode after normalization")
	}
	if !matchIgnore {
		t.Error("Expected Unix and Windows line endings to match in ignore mode after normalization")
	}
}

// Dimension 5: Unicode and Special Characters

func TestDiff_Dimension5_UnicodeWithWhitespaceDifferences(t *testing.T) {
	// Given
	expected := "hello    世界"
	actual := "hello 世界"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected Unicode strings with whitespace diff to match in ignore mode")
	}
}

func TestDiff_Dimension5_EmojiWithWhitespace(t *testing.T) {
	// Given
	expected := "test  🌍  emoji"
	actual := "test 🌍 emoji"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected emoji strings with whitespace diff to match in ignore mode")
	}
}

func TestDiff_Dimension5_RTLTextWithWhitespace(t *testing.T) {
	// Given - Arabic text
	expected := "Hello    مرحبا    World"
	actual := "Hello مرحبا World"

	// When
	_, match := Diff(expected, actual, true)

	// Then
	if !match {
		t.Error("Expected RTL text with whitespace diff to match in ignore mode")
	}
}

// ========================================
// Mode Toggle Pattern Tests
// ========================================

func TestDiff_ModeToggle_WhitespaceOnlyDifference_OppositeBehavior(t *testing.T) {
	// Given - same input with whitespace-only difference
	expected := "hello    world"
	actual := "hello world"

	// When - exact mode
	_, matchExact := Diff(expected, actual, false)

	// When - ignore mode
	_, matchIgnore := Diff(expected, actual, true)

	// Then - should have opposite results
	if matchExact {
		t.Error("Expected exact mode to not match whitespace difference")
	}
	if !matchIgnore {
		t.Error("Expected ignore mode to match whitespace difference")
	}
}

func TestDiff_ModeToggle_ContentDifference_SameBehavior(t *testing.T) {
	// Given - actual content difference
	expected := "hello world"
	actual := "hello mars"

	// When - exact mode
	_, matchExact := Diff(expected, actual, false)

	// When - ignore mode
	_, matchIgnore := Diff(expected, actual, true)

	// Then - both should not match
	if matchExact {
		t.Error("Expected exact mode to not match content difference")
	}
	if matchIgnore {
		t.Error("Expected ignore mode to not match content difference")
	}
}

func TestDiff_ModeToggle_MultilineWhitespaceOnly_OppositeBehavior(t *testing.T) {
	// Given
	expected := "line1  \nline2\t\nline3   "
	actual := "line1\nline2\nline3"

	// When - exact mode
	_, matchExact := Diff(expected, actual, false)

	// When - ignore mode
	_, matchIgnore := Diff(expected, actual, true)

	// Then
	if matchExact {
		t.Error("Expected exact mode to not match multiline whitespace difference")
	}
	if !matchIgnore {
		t.Error("Expected ignore mode to match multiline whitespace difference")
	}
}

func TestDiff_ModeToggle_LeadingWhitespace_OppositeBehavior(t *testing.T) {
	// Given
	expected := "   hello world"
	actual := "hello world"

	// When
	_, matchExact := Diff(expected, actual, false)
	_, matchIgnore := Diff(expected, actual, true)

	// Then
	if matchExact {
		t.Error("Expected exact mode to detect leading whitespace difference")
	}
	if !matchIgnore {
		t.Error("Expected ignore mode to ignore leading whitespace difference")
	}
}

func TestDiff_ModeToggle_TrailingWhitespace_OppositeBehavior(t *testing.T) {
	// Given
	expected := "hello world   "
	actual := "hello world"

	// When
	_, matchExact := Diff(expected, actual, false)
	_, matchIgnore := Diff(expected, actual, true)

	// Then
	if matchExact {
		t.Error("Expected exact mode to detect trailing whitespace difference")
	}
	if !matchIgnore {
		t.Error("Expected ignore mode to ignore trailing whitespace difference")
	}
}

func TestDiff_ModeToggle_TabsVsSpaces_OppositeBehavior(t *testing.T) {
	// Given
	expected := "hello\t\tworld"
	actual := "hello  world"

	// When
	_, matchExact := Diff(expected, actual, false)
	_, matchIgnore := Diff(expected, actual, true)

	// Then
	if matchExact {
		t.Error("Expected exact mode to detect tabs vs spaces difference")
	}
	if !matchIgnore {
		t.Error("Expected ignore mode to ignore tabs vs spaces difference")
	}
}

func TestDiff_ModeToggle_ComplexWhitespace_OppositeBehavior(t *testing.T) {
	// Given - complex whitespace scenario
	expected := "  hello  \t  world  \v  test  "
	actual := "hello world test"

	// When
	_, matchExact := Diff(expected, actual, false)
	_, matchIgnore := Diff(expected, actual, true)

	// Then
	if matchExact {
		t.Error("Expected exact mode to detect complex whitespace difference")
	}
	if !matchIgnore {
		t.Error("Expected ignore mode to ignore complex whitespace difference")
	}
}

// ========================================
// Line Ending Normalization Tests (Dimension 4 from spec)
// ========================================

func TestDiff_LineEndings_UnixVsWindows_ShouldMatch(t *testing.T) {
	// Given
	unix := "line1\nline2\nline3"
	windows := "line1\r\nline2\r\nline3"

	// When - both modes should normalize line endings
	_, matchExact := Diff(unix, windows, false)
	_, matchIgnore := Diff(unix, windows, true)

	// Then - should match after normalization
	if !matchExact {
		t.Error("Expected Unix and Windows line endings to match after normalization (exact mode)")
	}
	if !matchIgnore {
		t.Error("Expected Unix and Windows line endings to match after normalization (ignore mode)")
	}
}

func TestDiff_LineEndings_UnixVsMac_ShouldMatch(t *testing.T) {
	// Given
	unix := "line1\nline2\nline3"
	mac := "line1\rline2\rline3"

	// When
	_, matchExact := Diff(unix, mac, false)
	_, matchIgnore := Diff(unix, mac, true)

	// Then - should match after normalization
	if !matchExact {
		t.Error("Expected Unix and Mac line endings to match after normalization (exact mode)")
	}
	if !matchIgnore {
		t.Error("Expected Unix and Mac line endings to match after normalization (ignore mode)")
	}
}

func TestDiff_LineEndings_WindowsVsMac_ShouldMatch(t *testing.T) {
	// Given
	windows := "line1\r\nline2\r\nline3"
	mac := "line1\rline2\rline3"

	// When
	_, matchExact := Diff(windows, mac, false)

	// Then - should match after normalization
	if !matchExact {
		t.Error("Expected Windows and Mac line endings to match after normalization")
	}
}

func TestDiff_LineEndings_MixedInSameText_ShouldNormalize(t *testing.T) {
	// Given - mixed line endings in same string
	mixed1 := "line1\nline2\r\nline3\rline4"
	mixed2 := "line1\r\nline2\nline3\nline4"

	// When
	_, match := Diff(mixed1, mixed2, false)

	// Then - should match after normalization
	if !match {
		t.Error("Expected mixed line endings to match after normalization")
	}
}

func TestDiff_LineEndings_ConsecutiveNewlines_PreservesEmptyLines(t *testing.T) {
	// Given
	unix := "line1\n\n\nline2"          // 2 empty lines
	windows := "line1\r\n\r\n\r\nline2" // 2 empty lines

	// When
	_, match := Diff(unix, windows, false)

	// Then - should match with same empty line count
	if !match {
		t.Error("Expected consecutive newlines to preserve empty lines after normalization")
	}
}

func TestDiff_LineEndings_TrailingNewlineUnix_ShouldMatch(t *testing.T) {
	// Given
	unix := "line1\nline2\n"
	windows := "line1\r\nline2\r\n"

	// When
	_, match := Diff(unix, windows, false)

	// Then
	if !match {
		t.Error("Expected trailing newlines to match after normalization")
	}
}

func TestDiff_LineEndings_NoTrailingNewline_ShouldMatch(t *testing.T) {
	// Given
	unix := "line1\nline2"
	windows := "line1\r\nline2"

	// When
	_, match := Diff(unix, windows, false)

	// Then
	if !match {
		t.Error("Expected no trailing newlines to match after normalization")
	}
}

func TestDiff_LineEndings_MultipleConsecutiveNewlines_ShouldMatch(t *testing.T) {
	// Given
	unix := "line1\n\n\n\nline2"
	windows := "line1\r\n\r\n\r\n\r\nline2"

	// When
	_, match := Diff(unix, windows, false)

	// Then
	if !match {
		t.Error("Expected multiple consecutive newlines to match after normalization")
	}
}

func TestDiff_LineEndings_StartingWithNewline_ShouldMatch(t *testing.T) {
	// Given
	unix := "\nline1\nline2"
	windows := "\r\nline1\r\nline2"

	// When
	_, match := Diff(unix, windows, false)

	// Then
	if !match {
		t.Error("Expected newlines at start to match after normalization")
	}
}

func TestDiff_LineEndings_OnlyNewlines_ShouldMatch(t *testing.T) {
	// Given
	unix := "\n\n\n"
	windows := "\r\n\r\n\r\n"

	// When
	_, match := Diff(unix, windows, false)

	// Then
	if !match {
		t.Error("Expected strings with only newlines to match after normalization")
	}
}

func TestDiff_LineEndings_ComplexMixed_ShouldMatch(t *testing.T) {
	// Given - realistic scenario with mixed endings
	text1 := "line1\r\nline2\nline3\rline4\n"
	text2 := "line1\nline2\nline3\nline4\n"

	// When
	_, match := Diff(text1, text2, false)

	// Then
	if !match {
		t.Error("Expected complex mixed line endings to match after normalization")
	}
}
