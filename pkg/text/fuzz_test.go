package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"testing"
)

// FuzzStripMargin tests StripMargin with random inputs
func FuzzStripMargin(f *testing.F) {
	// Seed corpus with known inputs
	f.Add("|line 1\n|line 2")
	f.Add("\t|content")
	f.Add("  |test\n  |data")
	f.Add("")
	f.Add("|")
	f.Add("|||")
	f.Add("no pipes here")
	f.Add("|line with\ttabs")
	f.Add("|line with\r\nCRLF")

	f.Fuzz(func(t *testing.T, input string) {
		// Should not panic with any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("StripMargin panicked with input %q: %v", input, r)
			}
		}()

		result := text.StripMargin(input)

		// Basic sanity checks
		// Result should not be longer than input (we're removing characters)
		if len(result) > len(input) {
			t.Errorf("Result longer than input: input=%q, result=%q", input, result)
		}
	})
}

// FuzzStripColumn tests StripColumn with random inputs
func FuzzStripColumn(f *testing.F) {
	// Seed corpus
	f.Add("|line 1|\n|line 2|")
	f.Add("\t|content|")
	f.Add("|test|")
	f.Add("")
	f.Add("||")
	f.Add("|incomplete")
	f.Add("no pipes")
	f.Add("|with\ttabs|")

	f.Fuzz(func(t *testing.T, input string) {
		// Should not panic with any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("StripColumn panicked with input %q: %v", input, r)
			}
		}()

		result := text.StripColumn(input)

		// Result should not be longer than input
		if len(result) > len(input) {
			t.Errorf("Result longer than input: input=%q, result=%q", input, result)
		}
	})
}

// FuzzDiff tests Diff with random inputs
func FuzzDiff(f *testing.F) {
	// Seed corpus
	f.Add("hello", "hello")
	f.Add("hello", "world")
	f.Add("line1\nline2", "line1\nline3")
	f.Add("", "")
	f.Add("test", "")
	f.Add("", "test")
	f.Add("hello\tworld", "hello world")
	f.Add("a\nb\nc", "a\nb\nc")

	f.Fuzz(func(t *testing.T, expected, actual string) {
		// Should not panic with any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Diff panicked with expected=%q, actual=%q: %v", expected, actual, r)
			}
		}()

		output, match := text.Diff(expected, actual, false)

		// Sanity checks
		// If strings are identical, match should be true
		if expected == actual && !match {
			t.Errorf("Identical strings should match: expected=%q, actual=%q", expected, actual)
		}

		// If match is true, strings should be equal
		if match && expected != actual {
			t.Errorf("Match is true but strings differ: expected=%q, actual=%q", expected, actual)
		}

		// Output should not be empty
		if len(output) == 0 {
			t.Errorf("Diff output should not be empty")
		}
	})
}

// FuzzCompareStrings tests CompareStrings with random inputs
func FuzzCompareStrings(f *testing.F) {
	// Seed corpus
	f.Add("hello", "hello")
	f.Add("hello", "world")
	f.Add("", "")
	f.Add("test", "")
	f.Add("", "test")
	f.Add("hello\tworld", "hello world")
	f.Add("a b c", "a b c")

	f.Fuzz(func(t *testing.T, expected, actual string) {
		// Should not panic with any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CompareStrings panicked with expected=%q, actual=%q: %v", expected, actual, r)
			}
		}()

		result := text.CompareStrings(expected, actual)

		// Sanity checks
		// Result should contain either MATCH or ASSERTION_FAILED
		if !fuzzContainsAny(result, "MATCH", "ASSERTION_FAILED") {
			t.Errorf("Result should contain status indicator, got: %q", result)
		}

		// If strings are identical, result should contain MATCH
		if expected == actual {
			if !fuzzContains(result, "✓ [MATCH]") {
				t.Errorf("Identical strings should show MATCH, got: %q", result)
			}
		} else {
			if !fuzzContains(result, "✗ [ASSERTION_FAILED]") {
				t.Errorf("Different strings should show ASSERTION_FAILED, got: %q", result)
			}
		}

		// Output should not be empty
		if len(result) == 0 {
			t.Errorf("CompareStrings output should not be empty")
		}
	})
}

// FuzzCompareStringsRaw tests CompareStringsRaw with random inputs
func FuzzCompareStringsRaw(f *testing.F) {
	// Seed corpus
	f.Add("hello", "hello")
	f.Add("hello", "world")
	f.Add("", "")
	f.Add("test", "")
	f.Add("hello\tworld", "hello world")

	f.Fuzz(func(t *testing.T, expected, actual string) {
		// Should not panic with any input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CompareStringsRaw panicked with expected=%q, actual=%q: %v", expected, actual, r)
			}
		}()

		result := text.CompareStringsRaw(expected, actual)

		// Sanity checks
		// Result should contain either MATCH or ASSERTION_FAILED
		if !fuzzContainsAny(result, "MATCH", "ASSERTION_FAILED") {
			t.Errorf("Result should contain status indicator, got: %q", result)
		}

		// If strings are identical, result should contain MATCH
		if expected == actual {
			if !fuzzContains(result, "✓ [MATCH]") {
				t.Errorf("Identical strings should show MATCH, got: %q", result)
			}
		} else {
			if !fuzzContains(result, "✗ [ASSERTION_FAILED]") {
				t.Errorf("Different strings should show ASSERTION_FAILED, got: %q", result)
			}
		}

		// Output should not be empty
		if len(result) == 0 {
			t.Errorf("CompareStringsRaw output should not be empty")
		}
	})
}

// Helper function to check if a string contains a substring
func fuzzContains(s, substr string) bool {
	return len(s) >= len(substr) && fuzzFindSubstring(s, substr) >= 0
}

// Helper function to check if a string contains any of the given substrings
func fuzzContainsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if fuzzContains(s, substr) {
			return true
		}
	}
	return false
}

// Helper function to find substring position
func fuzzFindSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
