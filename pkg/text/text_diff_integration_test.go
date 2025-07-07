package text

import (
	"strings"
	"testing"
)

// Integration tests for function composition and workflow

func TestDiffIntegration_SimpleStringComparison_ProducesExpectedOutput(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello mars"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␣world ≠ hello␣mars |
		|      △             △    |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_IdenticalStrings_ProducesMatchOutput(t *testing.T) {
	// Given
	expected := "hello world"
	actual := "hello world"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true")
	}

	expectedOutput := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␣world | hello␣world|
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_MultilineWithMissingLines_ProducesCorrectArrows(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    ←         |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_ExtraLinesInActual_ProducesRightArrow(t *testing.T) {
	// Given
	expected := "line1\nline2"
	actual := "line1\nline2\nline3"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|         → line3   |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_WhitespaceVisualization_ShowsAllInvisibleChars(t *testing.T) {
	// Given
	expected := "hello\tworld"
	actual := "hello world"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected    | Actual     |
		|----------- | -----------|
		|hello␉world ≠ hello␣world|
		|     △             △     |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_TrailingNewlines_HandlesCorrectly(t *testing.T) {
	// Given
	expected := "line1\nline2\n"
	actual := "line1\nline2"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}
	if !result.HasTrailingNL {
		t.Errorf("Expected HasTrailingNL to be true")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|␤        ←         |
		||
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_EmptyStrings_ProducesEmptyComparison(t *testing.T) {
	// Given
	expected := ""
	actual := ""

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|         |         |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_UnicodeContent_HandlesCorrectly(t *testing.T) {
	// Given
	expected := "Hello 🌍 World"
	actual := "Hello 🚀 World"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected      | Actual       |
		|------------- | -------------|
		|Hello␣🌍␣World ≠ Hello␣🚀␣World|
		|      △               △      |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}
func TestDiffIntegration_CrossPlatformLineEndings_NormalizesAndMatches(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\r\nline2\r\nline3"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true after normalization")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    | line3   |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_ComplexWhitespace_ShowsAllSymbols(t *testing.T) {
	// Given
	expected := "hello\t\t world"
	actual := "hello   world"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected      | Actual       |
		|------------- | -------------|
		|hello␉␉␣world ≠ hello␣␣␣world|
		|     △               △       |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_FirstLineDiffers_StopsAtFirstDifference(t *testing.T) {
	// Given
	expected := "different line\nsame line\nanother same line"
	actual := "changed line\nsame line\nanother same line"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	// Should only show up to the first different line
	expectedOutput := StripColumn(`
		|Expected          | Actual           |
		|----------------- | -----------------|
		|different␣line    ≠ changed␣line     |
		|△                   △                |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}

	// Verify it doesn't contain the subsequent lines
	if strings.Contains(output, "same line") {
		t.Errorf("Expected output to stop at first difference, but found 'same line'")
	}
}

func TestDiffIntegration_MiddleLineDiffers_ShowsUpToFirstDifference(t *testing.T) {
	// Given
	expected := "same line 1\ndifferent line\nsame line 3"
	actual := "same line 1\nchanged line\nsame line 3"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected       | Actual        |
		|-------------- | --------------|
		|same␣line␣1    | same␣line␣1   |
		|different␣line ≠ changed␣line  |
		|△                △             |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}

	// Verify it doesn't contain the third line
	if strings.Contains(output, "same line 3") {
		t.Errorf("Expected output to stop at first difference, but found 'same line 3'")
	}
}

func TestDiffIntegration_LargeWidthDifference_AlignsProperly(t *testing.T) {
	// Given
	expected := "short"
	actual := "very long string indeed here"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected                     | Actual                      |
		|---------------------------- | ----------------------------|
		|short                        ≠ very␣long␣string␣indeed␣here|
		|△                              △                           |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_OnlyWhitespace_ShowsWhitespaceSymbols(t *testing.T) {
	// Given
	expected := "   "
	actual := "\t\t\t"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|␣␣␣      ≠ ␉␉␉     |
		|△          △       |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_MixedLineEndingsInInput_NormalizesCorrectly(t *testing.T) {
	// Given
	expected := "line1\nline2\nline3"
	actual := "line1\nline2\r\nline3\rline4"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|line2    | line2   |
		|line3    | line3   |
		|         → line4   |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_EdgeCaseEmptyLineHandling_WorksCorrectly(t *testing.T) {
	// Given
	expected := "line1\n\nline3"
	actual := "line1\n\nline3"

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if !result.Match {
		t.Errorf("Expected Match to be true")
	}

	expectedOutput := StripColumn(`
		|Expected | Actual  |
		|-------- | --------|
		|line1    | line1   |
		|         |         |
		|line3    | line3   |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_RealWorldExample_JSONComparison(t *testing.T) {
	// Given - simulating a common use case of comparing JSON-like structures
	expected := StripMargin(`
		|{
		|  "name": "John",
		|  "age": 30
		|}
	`)
	actual := StripMargin(`
		|{
		|  "name": "Jane",
		|  "age": 25
		|}
	`)

	// When
	result := computeDiff(expected, actual)
	output := renderDiff(result)

	// Then
	if result.Match {
		t.Errorf("Expected Match to be false")
	}

	expectedOutput := StripColumn(`
		|Expected          | Actual           |
		|----------------- | -----------------|
		|{                 | {                |
		|␣␣"name":␣"John", ≠ ␣␣"name":␣"Jane",|
		|            △                   △    |
	`)

	if output != expectedOutput {
		t.Errorf("Output mismatch:\n%s", compareMultilineStrings(output, expectedOutput))
	}
}

func TestDiffIntegration_PublicAPIConsistency_MatchesDiffFunction(t *testing.T) {
	// Given
	expected := "hello\tworld"
	actual := "hello world"

	// When - using the public API
	publicOutput, publicMatch := Diff(expected, actual)

	// When - using internal functions
	internalResult := computeDiff(expected, actual)
	internalOutput := renderDiff(internalResult)

	// Then - both approaches should produce identical results
	if publicMatch != internalResult.Match {
		t.Errorf("Match results differ: public=%v, internal=%v", publicMatch, internalResult.Match)
	}

	if publicOutput != internalOutput {
		t.Errorf("Output differs between public and internal APIs:\n%s", compareMultilineStrings(publicOutput, internalOutput))
	}
}
