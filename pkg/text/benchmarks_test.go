package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"strings"
	"testing"
)

// Benchmark StripMargin with small input
func BenchmarkStripMargin_Small(b *testing.B) {
	input := `
		|line 1
		|line 2
		|line 3
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripMargin(input)
	}
}

// Benchmark StripMargin with medium input (100 lines)
func BenchmarkStripMargin_Medium(b *testing.B) {
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = "\t|This is line content with some text"
	}
	input := strings.Join(lines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripMargin(input)
	}
}

// Benchmark StripMargin with large input (10,000 lines)
func BenchmarkStripMargin_Large(b *testing.B) {
	lines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		lines[i] = "\t|This is line content with some text"
	}
	input := strings.Join(lines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripMargin(input)
	}
}

// Benchmark StripColumn with small input
func BenchmarkStripColumn_Small(b *testing.B) {
	input := `
		|line 1|
		|line 2|
		|line 3|
	`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripColumn(input)
	}
}

// Benchmark StripColumn with medium input (100 lines)
func BenchmarkStripColumn_Medium(b *testing.B) {
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = "\t|This is line content with some text|"
	}
	input := strings.Join(lines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripColumn(input)
	}
}

// Benchmark StripColumn with large input (10,000 lines)
func BenchmarkStripColumn_Large(b *testing.B) {
	lines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		lines[i] = "\t|This is line content with some text|"
	}
	input := strings.Join(lines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripColumn(input)
	}
}

// Benchmark Diff with identical small strings
func BenchmarkDiff_IdenticalSmall(b *testing.B) {
	str := "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(str, str)
	}
}

// Benchmark Diff with different small strings
func BenchmarkDiff_DifferentSmall(b *testing.B) {
	expected := "hello world"
	actual := "hello mars"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(expected, actual)
	}
}

// Benchmark Diff with identical medium strings (100 lines)
func BenchmarkDiff_IdenticalMedium(b *testing.B) {
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = "This is line content with some text"
	}
	str := strings.Join(lines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(str, str)
	}
}

// Benchmark Diff with different medium strings (difference on line 50)
func BenchmarkDiff_DifferentMedium(b *testing.B) {
	expectedLines := make([]string, 100)
	actualLines := make([]string, 100)
	for i := 0; i < 100; i++ {
		if i == 50 {
			expectedLines[i] = "This line is different in expected"
			actualLines[i] = "This line is different in actual"
		} else {
			expectedLines[i] = "This is line content with some text"
			actualLines[i] = "This is line content with some text"
		}
	}
	expected := strings.Join(expectedLines, "\n")
	actual := strings.Join(actualLines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(expected, actual)
	}
}

// Benchmark Diff with identical large strings (10,000 lines)
func BenchmarkDiff_IdenticalLarge(b *testing.B) {
	lines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		lines[i] = "This is line content with some text"
	}
	str := strings.Join(lines, "\n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(str, str)
	}
}

// Benchmark Diff with different large strings (difference early)
func BenchmarkDiff_DifferentLargeEarly(b *testing.B) {
	expectedLines := make([]string, 10000)
	actualLines := make([]string, 10000)
	for i := 0; i < 10000; i++ {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(expected, actual)
	}
}

// Benchmark Diff with Unicode content
func BenchmarkDiff_Unicode(b *testing.B) {
	expected := "Hello 世界! 🌍 This is a test with Unicode characters"
	actual := "Hello 世界! 🌎 This is a test with Unicode characters"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(expected, actual)
	}
}

// Benchmark CompareStrings with matching strings
func BenchmarkCompareStrings_Matching(b *testing.B) {
	str := "hello world with some content"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.CompareStrings(str, str)
	}
}

// Benchmark CompareStrings with different strings
func BenchmarkCompareStrings_Different(b *testing.B) {
	expected := "hello world"
	actual := "hello mars"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.CompareStrings(expected, actual)
	}
}

// Benchmark CompareStrings with whitespace differences
func BenchmarkCompareStrings_Whitespace(b *testing.B) {
	expected := "hello world"
	actual := "hello\tworld"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.CompareStrings(expected, actual)
	}
}

// Benchmark CompareStringsRaw with matching strings
func BenchmarkCompareStringsRaw_Matching(b *testing.B) {
	str := "hello world with some content"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.CompareStringsRaw(str, str)
	}
}

// Benchmark CompareStringsRaw with different strings
func BenchmarkCompareStringsRaw_Different(b *testing.B) {
	expected := "hello world"
	actual := "hello mars"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.CompareStringsRaw(expected, actual)
	}
}

// Benchmark memory allocation for StripMargin
func BenchmarkStripMargin_Allocations(b *testing.B) {
	input := `
		|line 1
		|line 2
		|line 3
		|line 4
		|line 5
	`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.StripMargin(input)
	}
}

// Benchmark memory allocation for Diff
func BenchmarkDiff_Allocations(b *testing.B) {
	expected := "hello world\nline 2\nline 3"
	actual := "hello mars\nline 2\nline 3"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = text.Diff(expected, actual)
	}
}

// Benchmark memory allocation for CompareStrings
func BenchmarkCompareStrings_Allocations(b *testing.B) {
	expected := "hello world"
	actual := "hello mars"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = text.CompareStrings(expected, actual)
	}
}
