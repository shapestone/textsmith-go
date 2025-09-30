package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"strings"
	"testing"
	"time"
)

// TestStripMargin_WithVeryLargeInput_CompletesInReasonableTime tests performance with large input
func TestStripMargin_WithVeryLargeInput_CompletesInReasonableTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 100,000 lines (approximately 5MB of text)
	lines := make([]string, 100000)
	for i := 0; i < 100000; i++ {
		lines[i] = "\t|This is line content with some text that takes up space"
	}
	input := strings.Join(lines, "\n")

	// When
	start := time.Now()
	result := text.StripMargin(input)
	duration := time.Since(start)

	// Then: Should complete in less than 2 seconds (relaxed for CI)
	if duration > 2*time.Second {
		t.Errorf("StripMargin took too long: %v (expected < 2s)", duration)
	}

	// Verify result is correct
	if !strings.Contains(result, "This is line content") {
		t.Errorf("Expected result to contain processed content")
	}
}

// TestStripColumn_WithVeryLargeInput_CompletesInReasonableTime tests performance with large input
func TestStripColumn_WithVeryLargeInput_CompletesInReasonableTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 10,000 lines (StripColumn uses regex which is slower on very large inputs)
	lines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		lines[i] = "\t|This is line content with some text that takes up space|"
	}
	input := strings.Join(lines, "\n")

	// When
	start := time.Now()
	result := text.StripColumn(input)
	duration := time.Since(start)

	// Then: Should complete in less than 5 seconds (regex performance consideration)
	if duration > 5*time.Second {
		t.Errorf("StripColumn took too long: %v (expected < 5s)", duration)
	}

	// Verify result is correct
	if !strings.Contains(result, "This is line content") {
		t.Errorf("Expected result to contain processed content")
	}
}

// TestDiff_WithLargeIdenticalFiles_CompletesInReasonableTime tests large identical files
func TestDiff_WithLargeIdenticalFiles_CompletesInReasonableTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 50,000 identical lines
	lines := make([]string, 50000)
	for i := 0; i < 50000; i++ {
		lines[i] = "This is line content with some text that takes up space"
	}
	content := strings.Join(lines, "\n")

	// When
	start := time.Now()
	_, match := text.Diff(content, content)
	duration := time.Since(start)

	// Then: Should complete in less than 4 seconds and match (relaxed for CI)
	if duration > 4*time.Second {
		t.Errorf("Diff took too long: %v (expected < 4s)", duration)
	}

	if !match {
		t.Errorf("Expected identical large files to match")
	}
}

// TestDiff_WithLargeDifferentFiles_StopsAtFirstDifference tests early exit behavior
func TestDiff_WithLargeDifferentFiles_StopsAtFirstDifference(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 50,000 lines with difference on line 10
	expectedLines := make([]string, 50000)
	actualLines := make([]string, 50000)
	for i := 0; i < 50000; i++ {
		if i == 10 {
			expectedLines[i] = "This line is different in expected"
			actualLines[i] = "This line is different in actual"
		} else {
			expectedLines[i] = "This is line content with some text"
			actualLines[i] = "This is line content with some text"
		}
	}
	expected := strings.Join(expectedLines, "\n")
	actual := strings.Join(actualLines, "\n")

	// When
	start := time.Now()
	output, match := text.Diff(expected, actual)
	duration := time.Since(start)

	// Then: Should complete quickly due to early exit (< 1s, relaxed for CI)
	if duration > time.Second {
		t.Errorf("Diff took too long: %v (expected < 1s due to early exit)", duration)
	}

	if match {
		t.Errorf("Expected different files to not match")
	}

	// Should show the difference indicator (≠)
	if !strings.Contains(output, "≠") {
		t.Errorf("Expected output to show the difference indicator, got: %s", output)
	}
}

// TestDiff_WithVeryLongSingleLine_HandlesCorrectly tests very long single line
func TestDiff_WithVeryLongSingleLine_HandlesCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 1MB single line (no newlines)
	longString := strings.Repeat("a", 1024*1024)

	// When
	start := time.Now()
	_, match := text.Diff(longString, longString)
	duration := time.Since(start)

	// Then: Should complete in less than 2 seconds (relaxed for CI)
	if duration > 2*time.Second {
		t.Errorf("Diff with very long single line took too long: %v (expected < 2s)", duration)
	}

	if !match {
		t.Errorf("Expected identical long strings to match")
	}
}

// TestDiff_WithVeryLongSingleLineDifference_DetectsDifference tests long line with difference
func TestDiff_WithVeryLongSingleLineDifference_DetectsDifference(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 1MB single line with difference in the middle
	halfMB := strings.Repeat("a", 512*1024)
	expected := halfMB + "X" + halfMB
	actual := halfMB + "Y" + halfMB

	// When
	start := time.Now()
	output, match := text.Diff(expected, actual)
	duration := time.Since(start)

	// Then: Should complete in less than 2 seconds (relaxed for CI)
	if duration > 2*time.Second {
		t.Errorf("Diff with very long single line difference took too long: %v (expected < 2s)", duration)
	}

	if match {
		t.Errorf("Expected different long strings to not match")
	}

	// Should detect the difference
	if !strings.Contains(output, "≠") {
		t.Errorf("Expected output to show difference indicator")
	}
}

// TestCompareStrings_WithLongStrings_HandlesEfficiently tests CompareStrings performance
func TestCompareStrings_WithLongStrings_HandlesEfficiently(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: Long strings with difference at position 1000
	prefix := strings.Repeat("a", 1000)
	suffix := strings.Repeat("b", 1000)
	expected := prefix + "X" + suffix
	actual := prefix + "Y" + suffix

	// When
	start := time.Now()
	result := text.CompareStrings(expected, actual)
	duration := time.Since(start)

	// Then: Should complete in less than 100ms
	if duration > 100*time.Millisecond {
		t.Errorf("CompareStrings took too long: %v (expected < 100ms)", duration)
	}

	// Should show the difference at position 1000
	if !strings.Contains(result, "Difference at position 1000") {
		t.Errorf("Expected result to show difference at position 1000, got: %s", result)
	}
}

// TestStripMargin_WithUnicodeHeavyContent_HandlesCorrectly tests Unicode performance
func TestStripMargin_WithUnicodeHeavyContent_HandlesCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 10,000 lines of Unicode text (emojis, CJK, Arabic)
	lines := make([]string, 10000)
	unicodeContent := "Hello 世界 مرحبا 🌍 🌎 🌏"
	for i := 0; i < 10000; i++ {
		lines[i] = "\t|" + unicodeContent
	}
	input := strings.Join(lines, "\n")

	// When
	start := time.Now()
	result := text.StripMargin(input)
	duration := time.Since(start)

	// Then: Should complete in less than 1 second
	if duration > time.Second {
		t.Errorf("StripMargin with Unicode took too long: %v (expected < 1s)", duration)
	}

	// Verify Unicode is preserved
	if !strings.Contains(result, "世界") {
		t.Errorf("Expected result to contain Unicode characters")
	}
}

// TestDiff_WithManyShortLines_HandlesCorrectly tests many small lines
func TestDiff_WithManyShortLines_HandlesCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 100,000 short lines
	lines := make([]string, 100000)
	for i := 0; i < 100000; i++ {
		lines[i] = "line"
	}
	content := strings.Join(lines, "\n")

	// When
	start := time.Now()
	_, match := text.Diff(content, content)
	duration := time.Since(start)

	// Then: Should complete in less than 3 seconds
	if duration > 3*time.Second {
		t.Errorf("Diff with many short lines took too long: %v (expected < 3s)", duration)
	}

	if !match {
		t.Errorf("Expected identical content to match")
	}
}

// TestStripMargin_WithEmptyLines_HandlesEfficiently tests many empty lines
func TestStripMargin_WithEmptyLines_HandlesEfficiently(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: 50,000 empty lines with pipes
	lines := make([]string, 50000)
	for i := 0; i < 50000; i++ {
		lines[i] = "\t|"
	}
	input := strings.Join(lines, "\n")

	// When
	start := time.Now()
	result := text.StripMargin(input)
	duration := time.Since(start)

	// Then: Should complete in less than 1 second
	if duration > time.Second {
		t.Errorf("StripMargin with empty lines took too long: %v (expected < 1s)", duration)
	}

	// Result should be many empty lines
	expectedLineCount := 50000
	actualLineCount := strings.Count(result, "\n") + 1
	if actualLineCount != expectedLineCount {
		t.Errorf("Expected %d lines, got %d", expectedLineCount, actualLineCount)
	}
}

// TestCompareStrings_WithManyWhitespaceCharacters_VisualizesEfficiently tests whitespace handling
func TestCompareStrings_WithManyWhitespaceCharacters_VisualizesEfficiently(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Given: String with 10,000 spaces and tabs
	expected := strings.Repeat(" \t", 5000)
	actual := strings.Repeat("\t ", 5000)

	// When
	start := time.Now()
	result := text.CompareStrings(expected, actual)
	duration := time.Since(start)

	// Then: Should complete in less than 200ms
	if duration > 200*time.Millisecond {
		t.Errorf("CompareStrings with many whitespace chars took too long: %v (expected < 200ms)", duration)
	}

	// Should show visualization
	if !strings.Contains(result, "␣") || !strings.Contains(result, "␉") {
		t.Errorf("Expected result to contain whitespace visualization symbols")
	}
}