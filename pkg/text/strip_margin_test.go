package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"testing"
)

func TestStripMargin_WithSimpleMultilineString_ReturnsStrippedContent(t *testing.T) {
	// Given
	input := `
	|line 1
	|line 2
	|line 3
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithEmptyString_ReturnsEmptyString(t *testing.T) {
	// Given
	input := ""

	// When
	result := text.StripMargin(input)

	// Then
	if result != "" {
		t.Fatalf("Expected empty string, got %q", result)
	}
}

func TestStripMargin_WithNoMarginSymbol_ReturnsEmptyString(t *testing.T) {
	// Given
	input := `
	line 1
	line 2
	line 3
	`

	// When
	result := text.StripMargin(input)

	// Then
	if result != "" {
		t.Fatalf("Expected empty string, got %q", result)
	}
}

func TestStripMargin_WithSingleLine_ReturnsSingleLine(t *testing.T) {
	// Given
	input := "|single line"

	// When
	result := text.StripMargin(input)

	// Then
	expected := "single line"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithEmptyLines_PreservesEmptyLines(t *testing.T) {
	// Given
	input := `
	|line 1
	|
	|line 3
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line 1\n\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithVariousWhitespace_StripsWhitespaceBeforeMargin(t *testing.T) {
	// Given
	input := `
	|line 1
		|line 2
		    |line 3
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithTabsAndSpaces_StripsTabsAndSpacesBeforeMargin(t *testing.T) {
	// Given
	input := "\t |line 1\n  \t|line 2\n\t\t |line 3"

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithContentAfterMargin_PreservesContentSpacing(t *testing.T) {
	// Given
	input := `
	|  content with spaces  
	|	content with tab
	|normal content
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "  content with spaces  \n	content with tab\nnormal content"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithMixedLinesWithAndWithoutMargin_OnlyProcessesMarginLines(t *testing.T) {
	// Given
	input := `
	|line with margin
	line without margin
	|another line with margin
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line with margin\nanother line with margin"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithOnlyMarginSymbol_ReturnsEmptyLine(t *testing.T) {
	// Given
	input := "|"

	// When
	result := text.StripMargin(input)

	// Then
	expected := ""
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithMultipleMarginSymbolsInLine_ProcessesFirstMarginOnly(t *testing.T) {
	// Given
	input := `
	|content | more content
	|another | line
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "content | more content\nanother | line"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithSpecialCharacters_PreservesSpecialCharacters(t *testing.T) {
	// Given
	input := `
	|special chars: !@#$%^&*()
	|unicode: 🚀 ñ é
	|quotes: "hello" 'world'
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "special chars: !@#$%^&*()\nunicode: 🚀 ñ é\nquotes: \"hello\" 'world'"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithTrailingNewlines_HandlesTrailingNewlinesCorrectly(t *testing.T) {
	// Given
	input := "|line 1\n|line 2\n"

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line 1\nline 2"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithNoNewlineAtEnd_HandlesCorrectly(t *testing.T) {
	// Given
	input := "|line 1\n|line 2"

	// When
	result := text.StripMargin(input)

	// Then
	expected := "line 1\nline 2"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithMultiByteEmoji_PreservesEmoji(t *testing.T) {
	// Given
	input := `
	|🚀 rocket launch
	|🌟 star bright
	|🎉 celebration time
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "🚀 rocket launch\n🌟 star bright\n🎉 celebration time"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithCJKCharacters_PreservesCJKCharacters(t *testing.T) {
	// Given
	input := `
	|こんにちは world
	|你好 world
	|안녕하세요 world
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "こんにちは world\n你好 world\n안녕하세요 world"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithArabicText_PreservesArabicText(t *testing.T) {
	// Given
	input := `
	|مرحبا بك
	|اللغة العربية
	|نص تجريبي
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "مرحبا بك\nاللغة العربية\nنص تجريبي"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithCombiningCharacters_PreservesCombiningCharacters(t *testing.T) {
	// Given
	input := `
	|café with é (e + ́)
	|naïve with ï (i + ̈)
	|résumé with é (e + ́)
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "café with é (e + ́)\nnaïve with ï (i + ̈)\nrésumé with é (e + ́)"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithComplexUnicodeEmoji_PreservesComplexEmoji(t *testing.T) {
	// Given
	input := `
	|👨‍💻 developer
	|👩‍🚀 astronaut
	|🏳️‍🌈 rainbow flag
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "👨‍💻 developer\n👩‍🚀 astronaut\n🏳️‍🌈 rainbow flag"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithMixedUnicodeAndASCII_PreservesBothCorrectly(t *testing.T) {
	// Given
	input := `
	|ASCII text
	|🌍 Unicode emoji
	|普通话 Chinese
	|العربية Arabic
	|More ASCII
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "ASCII text\n🌍 Unicode emoji\n普通话 Chinese\nالعربية Arabic\nMore ASCII"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripMargin_WithUnicodeWhitespace_HandlesUnicodeWhitespace(t *testing.T) {
	// Given - using regular space and tab before margin, unicode content after
	input := "\t |🚀\u0020space\u00A0nbsp\u2003emspace\n  |test\u2009thinspace"

	// When
	result := text.StripMargin(input)

	// Then
	expected := "🚀\u0020space\u00A0nbsp\u2003emspace\ntest\u2009thinspace"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithSimpleMultilineString_ReturnsStrippedContent(t *testing.T) {
	// Given
	input := `
	|line 1|
	|line 2|
	|line 3|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithEmptyString_ReturnsEmptyString(t *testing.T) {
	// Given
	input := ""

	// When
	result := text.StripColumn(input)

	// Then
	if result != "" {
		t.Fatalf("Expected empty string, got %q", result)
	}
}

func TestStripColumn_WithNoColumnSymbols_ReturnsEmptyString(t *testing.T) {
	// Given
	input := `
	line 1
	line 2
	line 3
	`

	// When
	result := text.StripColumn(input)

	// Then
	if result != "" {
		t.Fatalf("Expected empty string, got %q", result)
	}
}

func TestStripColumn_WithOnlyOpeningPipe_ReturnsEmptyString(t *testing.T) {
	// Given
	input := `
	|line 1
	|line 2
	|line 3
	`

	// When
	result := text.StripColumn(input)

	// Then
	if result != "" {
		t.Fatalf("Expected empty string, got %q", result)
	}
}

func TestStripColumn_WithSingleLine_ReturnsSingleLine(t *testing.T) {
	// Given
	input := "|single line|"

	// When
	result := text.StripColumn(input)

	// Then
	expected := "single line"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithEmptyLines_PreservesEmptyLines(t *testing.T) {
	// Given
	input := `
	|line 1|
	||
	|line 3|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 1\n\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithVariousWhitespace_StripsWhitespaceBeforeOpeningPipe(t *testing.T) {
	// Given
	input := `
	|line 1|
		|line 2|
		    |line 3|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithTabsAndSpaces_StripsTabsAndSpacesAroundPipes(t *testing.T) {
	// Given
	input := "\t |line 1| \n  \t|line 2|\t\n\t\t |line 3|  \t"

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 1\nline 2\nline 3"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithContentSpacing_PreservesInternalSpacing(t *testing.T) {
	// Given
	input := `
	|  content with spaces  |
	|	content with tab	|
	|normal content|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "  content with spaces  \n	content with tab	\nnormal content"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithMixedLinesWithAndWithoutColumns_OnlyProcessesColumnLines(t *testing.T) {
	// Given
	input := `
	|line with columns|
	line without columns
	|another line with columns|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line with columns\nanother line with columns"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithOnlyPipeSymbols_ReturnsEmptyLine(t *testing.T) {
	// Given
	input := "||"

	// When
	result := text.StripColumn(input)

	// Then
	expected := ""
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithMultiplePipeSymbolsInContent_PreservesInternalPipes(t *testing.T) {
	// Given
	input := `
	|content | with | pipes|
	|another | line | here|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "content | with | pipes\nanother | line | here"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithSpecialCharacters_PreservesSpecialCharacters(t *testing.T) {
	// Given
	input := `
	|special chars: !@#$%^&*()|
	|unicode: 🚀 ñ é|
	|quotes: "hello" 'world'|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "special chars: !@#$%^&*()\nunicode: 🚀 ñ é\nquotes: \"hello\" 'world'"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithTrailingNewlines_HandlesTrailingNewlinesCorrectly(t *testing.T) {
	// Given
	input := "|line 1|\n|line 2|\n"

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 1\nline 2"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithNoNewlineAtEnd_HandlesCorrectly(t *testing.T) {
	// Given
	input := "|line 1|\n|line 2|"

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 1\nline 2"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithMultiByteEmoji_PreservesEmoji(t *testing.T) {
	// Given
	input := `
	|🚀 rocket launch|
	|🌟 star bright|
	|🎉 celebration time|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "🚀 rocket launch\n🌟 star bright\n🎉 celebration time"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithCJKCharacters_PreservesCJKCharacters(t *testing.T) {
	// Given
	input := `
	|こんにちは world|
	|你好 world|
	|안녕하세요 world|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "こんにちは world\n你好 world\n안녕하세요 world"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithArabicText_PreservesArabicText(t *testing.T) {
	// Given
	input := `
	|مرحبا بك|
	|اللغة العربية|
	|نص تجريبي|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "مرحبا بك\nاللغة العربية\nنص تجريبي"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithCombiningCharacters_PreservesCombiningCharacters(t *testing.T) {
	// Given
	input := `
	|café with é (e + ́)|
	|naïve with ï (i + ̈)|
	|résumé with é (e + ́)|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "café with é (e + ́)\nnaïve with ï (i + ̈)\nrésumé with é (e + ́)"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithComplexUnicodeEmoji_PreservesComplexEmoji(t *testing.T) {
	// Given
	input := `
	|👨‍💻 developer|
	|👩‍🚀 astronaut|
	|🏳️‍🌈 rainbow flag|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "👨‍💻 developer\n👩‍🚀 astronaut\n🏳️‍🌈 rainbow flag"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithMixedUnicodeAndASCII_PreservesBothCorrectly(t *testing.T) {
	// Given
	input := `
	|ASCII text|
	|🌍 Unicode emoji|
	|普通话 Chinese|
	|العربية Arabic|
	|More ASCII|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "ASCII text\n🌍 Unicode emoji\n普通话 Chinese\nالعربية Arabic\nMore ASCII"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithUnicodeWhitespace_HandlesUnicodeWhitespace(t *testing.T) {
	// Given - using regular space and tab before/after pipes, unicode content within
	input := "\t |🚀\u0020space\u00A0nbsp\u2003emspace| \n  |test\u2009thinspace|\t"

	// When
	result := text.StripColumn(input)

	// Then
	expected := "🚀\u0020space\u00A0nbsp\u2003emspace\ntest\u2009thinspace"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithIncompleteColumns_ReturnsEmptyString(t *testing.T) {
	// Given
	input := `
	|line 1
	|line 2|
	line 3|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "line 2"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestStripColumn_WithMalformedColumnSyntax_OnlyProcessesValidLines(t *testing.T) {
	// Given
	input := `
	|valid line|
	|incomplete line
	incomplete| line
	|another valid|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "valid line\nanother valid"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}
