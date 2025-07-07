package text_test

import (
	"encoding/xml"
	"github.com/shapestone/textsmith/pkg/text"
	"testing"
)

// XML StripMargin Tests

func TestStripMargin_WithValidXMLDocument_ReturnsValidXML(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0" encoding="UTF-8"?>
	|<person>
	|  <name>John Doe</name>
	|  <age>30</age>
	|  <email>john@example.com</email>
	|</person>
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<person>\n  <name>John Doe</name>\n  <age>30</age>\n  <email>john@example.com</email>\n</person>"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid XML
	var person struct {
		Name  string `xml:"name"`
		Age   int    `xml:"age"`
		Email string `xml:"email"`
	}
	if err := xml.Unmarshal([]byte(result), &person); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}
}

func TestStripMargin_WithXMLAttributes_PreservesAttributes(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<book id="123" category="fiction" available="true">
	|  <title lang="en">The Great Adventure</title>
	|  <author nationality="US">Jane Smith</author>
	|  <price currency="USD">19.99</price>
	|</book>
	`

	// When
	result := text.StripMargin(input)

	// Then
	var book struct {
		ID        string `xml:"id,attr"`
		Category  string `xml:"category,attr"`
		Available string `xml:"available,attr"`
		Title     struct {
			Lang string `xml:"lang,attr"`
			Text string `xml:",chardata"`
		} `xml:"title"`
		Author struct {
			Nationality string `xml:"nationality,attr"`
			Text        string `xml:",chardata"`
		} `xml:"author"`
		Price struct {
			Currency string `xml:"currency,attr"`
			Text     string `xml:",chardata"`
		} `xml:"price"`
	}
	if err := xml.Unmarshal([]byte(result), &book); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify attributes are preserved
	if book.ID != "123" {
		t.Fatalf("Expected ID '123', got %q", book.ID)
	}
	if book.Title.Lang != "en" {
		t.Fatalf("Expected title lang 'en', got %q", book.Title.Lang)
	}
}

func TestStripMargin_WithXMLNamespaces_PreservesNamespaces(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<root xmlns:book="http://example.com/book" xmlns:author="http://example.com/author">
	|  <book:catalog>
	|    <book:item id="1">
	|      <book:title>XML Processing</book:title>
	|      <author:name>John Developer</author:name>
	|    </book:item>
	|  </book:catalog>
	|</root>
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Just verify it parses as valid XML - namespace handling is complex
	var root interface{}
	if err := xml.Unmarshal([]byte(result), &root); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify namespace declarations are preserved in the output
	if !containsSubstring(result, `xmlns:book="http://example.com/book"`) {
		t.Fatal("Book namespace declaration not preserved")
	}
	if !containsSubstring(result, `xmlns:author="http://example.com/author"`) {
		t.Fatal("Author namespace declaration not preserved")
	}
}

func TestStripMargin_WithXMLCDATA_PreservesCDATA(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<document>
	|  <description><![CDATA[This contains <special> characters & symbols]]></description>
	|  <code><![CDATA[
	|    function example() {
	|      return "Hello & Goodbye";
	|    }
	|  ]]></code>
	|</document>
	`

	// When
	result := text.StripMargin(input)

	// Then
	var doc struct {
		Description string `xml:"description"`
		Code        string `xml:"code"`
	}
	if err := xml.Unmarshal([]byte(result), &doc); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify CDATA content is preserved
	if !containsSubstring(doc.Description, "This contains <special> characters & symbols") {
		t.Fatalf("CDATA content not preserved in description: %q", doc.Description)
	}
}

func TestStripMargin_WithXMLComments_PreservesComments(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<!-- This is a root comment -->
	|<document>
	|  <!-- User information section -->
	|  <user>
	|    <name>John Doe</name>
	|    <!-- TODO: Add more user details -->
	|  </user>
	|</document>
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Verify comments are preserved in the output
	if !containsSubstring(result, "<!-- This is a root comment -->") {
		t.Fatal("Root comment not preserved")
	}
	if !containsSubstring(result, "<!-- User information section -->") {
		t.Fatal("Section comment not preserved")
	}
	if !containsSubstring(result, "<!-- TODO: Add more user details -->") {
		t.Fatal("TODO comment not preserved")
	}

	// Verify it's still valid XML (comments should be ignored during parsing)
	var doc struct {
		User struct {
			Name string `xml:"name"`
		} `xml:"user"`
	}
	if err := xml.Unmarshal([]byte(result), &doc); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}
}

func TestStripMargin_WithXMLSpecialCharacters_PreservesEscaping(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<data>
	|  <text>AT&amp;T Corporation</text>
	|  <formula>5 &lt; 10 &amp;&amp; 10 &gt; 5</formula>
	|  <quote>He said &quot;Hello&quot; to me</quote>
	|  <apostrophe>That&apos;s correct</apostrophe>
	|</data>
	`

	// When
	result := text.StripMargin(input)

	// Then
	var data struct {
		Text       string `xml:"text"`
		Formula    string `xml:"formula"`
		Quote      string `xml:"quote"`
		Apostrophe string `xml:"apostrophe"`
	}
	if err := xml.Unmarshal([]byte(result), &data); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify special characters are properly decoded
	if data.Text != "AT&T Corporation" {
		t.Fatalf("Expected 'AT&T Corporation', got %q", data.Text)
	}
	if data.Formula != "5 < 10 && 10 > 5" {
		t.Fatalf("Expected '5 < 10 && 10 > 5', got %q", data.Formula)
	}
	if data.Quote != "He said \"Hello\" to me" {
		t.Fatalf("Expected 'He said \"Hello\" to me', got %q", data.Quote)
	}
}

func TestStripMargin_WithXMLUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0" encoding="UTF-8"?>
	|<multilingual>
	|  <chinese>你好世界</chinese>
	|  <arabic>مرحبا بالعالم</arabic>
	|  <emoji>🚀🌟🎉</emoji>
	|  <complex_emoji>👨‍💻👩‍🚀🏳️‍🌈</complex_emoji>
	|  <combined>café with é (e + ́)</combined>
	|</multilingual>
	`

	// When
	result := text.StripMargin(input)

	// Then
	var multilingual struct {
		Chinese      string `xml:"chinese"`
		Arabic       string `xml:"arabic"`
		Emoji        string `xml:"emoji"`
		ComplexEmoji string `xml:"complex_emoji"`
		Combined     string `xml:"combined"`
	}
	if err := xml.Unmarshal([]byte(result), &multilingual); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify Unicode content is preserved
	if multilingual.Chinese != "你好世界" {
		t.Fatalf("Chinese text not preserved: %q", multilingual.Chinese)
	}
	if multilingual.Arabic != "مرحبا بالعالم" {
		t.Fatalf("Arabic text not preserved: %q", multilingual.Arabic)
	}
	if multilingual.Emoji != "🚀🌟🎉" {
		t.Fatalf("Emoji not preserved: %q", multilingual.Emoji)
	}
}

func TestStripMargin_WithEmptyXMLElements_PreservesEmptyElements(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<document>
	|  <empty_self_closing/>
	|  <empty_with_content></empty_with_content>
	|  <container>
	|    <nested_empty/>
	|    <with_content>Some text</with_content>
	|  </container>
	|</document>
	`

	// When
	result := text.StripMargin(input)

	// Then
	var doc struct {
		EmptySelfClosing string `xml:"empty_self_closing"`
		EmptyWithContent string `xml:"empty_with_content"`
		Container        struct {
			NestedEmpty string `xml:"nested_empty"`
			WithContent string `xml:"with_content"`
		} `xml:"container"`
	}
	if err := xml.Unmarshal([]byte(result), &doc); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify empty elements are handled correctly
	if doc.Container.WithContent != "Some text" {
		t.Fatalf("Expected 'Some text', got %q", doc.Container.WithContent)
	}
}

func TestStripMargin_WithMalformedXML_ReturnsStringContent(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>
	|<document>
	|  <unclosed_tag>Some content
	|  <another>Valid content</another>
	|</document>
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "<?xml version=\"1.0\"?>\n<document>\n  <unclosed_tag>Some content\n  <another>Valid content</another>\n</document>"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's invalid XML (should fail to parse)
	var doc interface{}
	if err := xml.Unmarshal([]byte(result), &doc); err == nil {
		t.Fatal("Expected invalid XML to fail parsing, but it succeeded")
	}
}

// XML StripColumn Tests

func TestStripColumn_WithValidXMLDocument_ReturnsValidXML(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0" encoding="UTF-8"?>|
	|<person>|
	|  <name>John Doe</name>|
	|  <age>30</age>|
	|  <email>john@example.com</email>|
	|</person>|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<person>\n  <name>John Doe</name>\n  <age>30</age>\n  <email>john@example.com</email>\n</person>"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid XML
	var person struct {
		Name  string `xml:"name"`
		Age   int    `xml:"age"`
		Email string `xml:"email"`
	}
	if err := xml.Unmarshal([]byte(result), &person); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}
}

func TestStripColumn_WithXMLAttributes_PreservesAttributes(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>|
	|<book id="123" category="fiction" available="true">|
	|  <title lang="en">The Great Adventure</title>|
	|  <author nationality="US">Jane Smith</author>|
	|  <price currency="USD">19.99</price>|
	|</book>|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var book struct {
		ID        string `xml:"id,attr"`
		Category  string `xml:"category,attr"`
		Available string `xml:"available,attr"`
		Title     struct {
			Lang string `xml:"lang,attr"`
			Text string `xml:",chardata"`
		} `xml:"title"`
		Author struct {
			Nationality string `xml:"nationality,attr"`
			Text        string `xml:",chardata"`
		} `xml:"author"`
		Price struct {
			Currency string `xml:"currency,attr"`
			Text     string `xml:",chardata"`
		} `xml:"price"`
	}
	if err := xml.Unmarshal([]byte(result), &book); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify attributes are preserved
	if book.ID != "123" {
		t.Fatalf("Expected ID '123', got %q", book.ID)
	}
	if book.Title.Lang != "en" {
		t.Fatalf("Expected title lang 'en', got %q", book.Title.Lang)
	}
}

func TestStripColumn_WithXMLNamespaces_PreservesNamespaces(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>|
	|<root xmlns:book="http://example.com/book" xmlns:author="http://example.com/author">|
	|  <book:catalog>|
	|    <book:item id="1">|
	|      <book:title>XML Processing</book:title>|
	|      <author:name>John Developer</author:name>|
	|    </book:item>|
	|  </book:catalog>|
	|</root>|
	`

	// When
	result := text.StripColumn(input)

	// Then
	// Just verify it parses as valid XML - namespace handling is complex
	var root interface{}
	if err := xml.Unmarshal([]byte(result), &root); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify namespace declarations are preserved in the output
	if !containsSubstring(result, `xmlns:book="http://example.com/book"`) {
		t.Fatal("Book namespace declaration not preserved")
	}
	if !containsSubstring(result, `xmlns:author="http://example.com/author"`) {
		t.Fatal("Author namespace declaration not preserved")
	}
}

func TestStripColumn_WithXMLCDATA_PreservesCDATA(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>|
	|<document>|
	|  <description><![CDATA[This contains <special> characters & symbols]]></description>|
	|  <code><![CDATA[|
	|    function example() {|
	|      return "Hello & Goodbye";|
	|    }|
	|  ]]></code>|
	|</document>|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var doc struct {
		Description string `xml:"description"`
		Code        string `xml:"code"`
	}
	if err := xml.Unmarshal([]byte(result), &doc); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify CDATA content is preserved
	if !containsSubstring(doc.Description, "This contains <special> characters & symbols") {
		t.Fatalf("CDATA content not preserved in description: %q", doc.Description)
	}
}

func TestStripColumn_WithXMLUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0" encoding="UTF-8"?>|
	|<multilingual>|
	|  <chinese>你好世界</chinese>|
	|  <arabic>مرحبا بالعالم</arabic>|
	|  <emoji>🚀🌟🎉</emoji>|
	|  <complex_emoji>👨‍💻👩‍🚀🏳️‍🌈</complex_emoji>|
	|</multilingual>|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var multilingual struct {
		Chinese      string `xml:"chinese"`
		Arabic       string `xml:"arabic"`
		Emoji        string `xml:"emoji"`
		ComplexEmoji string `xml:"complex_emoji"`
	}
	if err := xml.Unmarshal([]byte(result), &multilingual); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}

	// Verify Unicode content is preserved
	if multilingual.Chinese != "你好世界" {
		t.Fatalf("Chinese text not preserved: %q", multilingual.Chinese)
	}
	if multilingual.Emoji != "🚀🌟🎉" {
		t.Fatalf("Emoji not preserved: %q", multilingual.Emoji)
	}
}

func TestStripColumn_WithCompactXML_PreservesCompactFormat(t *testing.T) {
	// Given
	input := `|<person><name>John</name><age>30</age></person>|`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "<person><name>John</name><age>30</age></person>"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid XML
	var person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}
	if err := xml.Unmarshal([]byte(result), &person); err != nil {
		t.Fatalf("Result is not valid XML: %v", err)
	}
}

func TestStripColumn_WithMalformedXML_ReturnsStringContent(t *testing.T) {
	// Given
	input := `
	|<?xml version="1.0"?>|
	|<document>|
	|  <unclosed_tag>Some content|
	|  <another>Valid content</another>|
	|</document>|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "<?xml version=\"1.0\"?>\n<document>\n  <unclosed_tag>Some content\n  <another>Valid content</another>\n</document>"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's invalid XML (should fail to parse)
	var doc interface{}
	if err := xml.Unmarshal([]byte(result), &doc); err == nil {
		t.Fatal("Expected invalid XML to fail parsing, but it succeeded")
	}
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
