package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"strings"
	"testing"
)

// YAML StripMargin Tests

func TestStripMargin_WithValidYAMLDocument_ReturnsValidYAML(t *testing.T) {
	// Given
	input := `
	|name: John Doe
	|age: 30
	|email: john@example.com
	|active: true
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "name: John Doe\nage: 30\nemail: john@example.com\nactive: true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify YAML structure is preserved (basic validation)
	if !containsSubstring(result, "name: John Doe") {
		t.Fatal("YAML key-value pair not preserved")
	}
	if !containsSubstring(result, "active: true") {
		t.Fatal("YAML boolean value not preserved")
	}
}

func TestStripMargin_WithYAMLNestedStructure_PreservesIndentation(t *testing.T) {
	// Given
	input := `
	|person:
	|  name: Jane Smith
	|  address:
	|    street: 123 Main St
	|    city: Anytown
	|    zipcode: "12345"
	|  hobbies:
	|    - reading
	|    - swimming
	|    - coding
	|active: true
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "person:\n  name: Jane Smith\n  address:\n    street: 123 Main St\n    city: Anytown\n    zipcode: \"12345\"\n  hobbies:\n    - reading\n    - swimming\n    - coding\nactive: true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify nested structure indentation is preserved
	lines := strings.Split(result, "\n")
	if !strings.HasPrefix(lines[1], "  name:") {
		t.Fatal("First level indentation not preserved")
	}
	if !strings.HasPrefix(lines[3], "    street:") {
		t.Fatal("Second level indentation not preserved")
	}
	if !strings.HasPrefix(lines[7], "    - reading") {
		t.Fatal("Array item indentation not preserved")
	}
}

func TestStripMargin_WithYAMLArrays_PreservesArrayStructure(t *testing.T) {
	// Given
	input := `
	|fruits:
	|  - apple
	|  - banana
	|  - cherry
	|numbers:
	|  - 1
	|  - 2
	|  - 3.14
	|mixed_array:
	|  - "string"
	|  - 42
	|  - true
	|  - null
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "fruits:\n  - apple\n  - banana\n  - cherry\nnumbers:\n  - 1\n  - 2\n  - 3.14\nmixed_array:\n  - \"string\"\n  - 42\n  - true\n  - null"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify array structure is preserved
	if !containsSubstring(result, "fruits:\n  - apple") {
		t.Fatal("Array structure not preserved")
	}
	if !containsSubstring(result, "- 3.14") {
		t.Fatal("Float value in array not preserved")
	}
	if !containsSubstring(result, "- null") {
		t.Fatal("Null value in array not preserved")
	}
}

func TestStripMargin_WithYAMLComplexObjects_PreservesComplexStructure(t *testing.T) {
	// Given
	input := `
	|database:
	|  host: localhost
	|  port: 5432
	|  credentials:
	|    username: admin
	|    password: secret123
	|  pools:
	|    - name: main
	|      size: 10
	|      timeout: 30
	|    - name: backup
	|      size: 5
	|      timeout: 60
	|logging:
	|  level: info
	|  file: /var/log/app.log
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Verify complex nested structure is preserved
	if !containsSubstring(result, "database:\n  host: localhost") {
		t.Fatal("Database section structure not preserved")
	}

	// Check nested credentials
	if !containsSubstring(result, "  credentials:\n    username: admin") {
		t.Fatal("Nested credentials structure not preserved")
	}

	// Check array of objects
	if !containsSubstring(result, "  pools:\n    - name: main\n      size: 10") {
		t.Fatal("Array of objects structure not preserved")
	}
}

func TestStripMargin_WithYAMLSpecialCharacters_PreservesSpecialChars(t *testing.T) {
	// Given
	input := `
	|company: "AT&T Corporation"
	|formula: "5 < 10 && 10 > 5"
	|quote: 'He said "Hello" to me'
	|path: "C:\\Users\\John\\Documents"
	|multiline: |
	|  This is a multiline
	|  string that spans
	|  multiple lines
	|folded: >
	|  This is a folded
	|  string that will be
	|  on one line
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Verify special characters and string formats are preserved
	if !containsSubstring(result, `company: "AT&T Corporation"`) {
		t.Fatal("Special characters in quoted string not preserved")
	}
	if !containsSubstring(result, `formula: "5 < 10 && 10 > 5"`) {
		t.Fatal("Formula with special characters not preserved")
	}
	if !containsSubstring(result, `quote: 'He said "Hello" to me'`) {
		t.Fatal("Mixed quotes not preserved")
	}
	if !containsSubstring(result, "multiline: |\n  This is a multiline") {
		t.Fatal("Literal block scalar not preserved")
	}
	if !containsSubstring(result, "folded: >\n  This is a folded") {
		t.Fatal("Folded block scalar not preserved")
	}
}

func TestStripMargin_WithYAMLUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	input := `
	|languages:
	|  chinese: 你好世界
	|  arabic: مرحبا بالعالم
	|  japanese: こんにちは世界
	|  korean: 안녕하세요 세계
	|emojis:
	|  simple: 🚀🌟🎉
	|  complex: 👨‍💻👩‍🚀🏳️‍🌈
	|special_chars:
	|  accented: café naïve résumé
	|  combined: "café with é (e + ́)"
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Verify Unicode content is preserved
	if !containsSubstring(result, "chinese: 你好世界") {
		t.Fatal("Chinese Unicode not preserved")
	}
	if !containsSubstring(result, "arabic: مرحبا بالعالم") {
		t.Fatal("Arabic Unicode not preserved")
	}
	if !containsSubstring(result, "simple: 🚀🌟🎉") {
		t.Fatal("Emoji Unicode not preserved")
	}
	if !containsSubstring(result, "complex: 👨‍💻👩‍🚀🏳️‍🌈") {
		t.Fatal("Complex emoji Unicode not preserved")
	}
	if !containsSubstring(result, "accented: café naïve résumé") {
		t.Fatal("Accented characters not preserved")
	}
}

func TestStripMargin_WithYAMLDataTypes_PreservesDataTypes(t *testing.T) {
	// Given
	input := `
	|string_value: "Hello World"
	|integer_value: 42
	|float_value: 3.14159
	|boolean_true: true
	|boolean_false: false
	|null_value: null
	|empty_string: ""
	|scientific_notation: 1.23e-4
	|negative_number: -273.15
	|large_integer: 1234567890123456789
	|quoted_number: "123"
	|yaml_timestamp: 2023-12-25T10:30:00Z
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Verify different data types are preserved
	if !containsSubstring(result, "integer_value: 42") {
		t.Fatal("Integer value not preserved")
	}
	if !containsSubstring(result, "float_value: 3.14159") {
		t.Fatal("Float value not preserved")
	}
	if !containsSubstring(result, "boolean_true: true") {
		t.Fatal("Boolean true not preserved")
	}
	if !containsSubstring(result, "boolean_false: false") {
		t.Fatal("Boolean false not preserved")
	}
	if !containsSubstring(result, "null_value: null") {
		t.Fatal("Null value not preserved")
	}
	if !containsSubstring(result, "scientific_notation: 1.23e-4") {
		t.Fatal("Scientific notation not preserved")
	}
	if !containsSubstring(result, `quoted_number: "123"`) {
		t.Fatal("Quoted number not preserved")
	}
}

func TestStripMargin_WithYAMLComments_PreservesComments(t *testing.T) {
	// Given
	input := `
	|# Application configuration
	|app:
	|  name: MyApp  # Application name
	|  version: 1.0.0
	|  # Database settings
	|  database:
	|    host: localhost  # Default host
	|    port: 5432       # PostgreSQL default port
	|# End of configuration
	`

	// When
	result := text.StripMargin(input)

	// Then
	// Verify comments are preserved
	if !containsSubstring(result, "# Application configuration") {
		t.Fatal("Header comment not preserved")
	}
	if !containsSubstring(result, "name: MyApp  # Application name") {
		t.Fatal("Inline comment not preserved")
	}
	if !containsSubstring(result, "  # Database settings") {
		t.Fatal("Indented comment not preserved")
	}
	if !containsSubstring(result, "# End of configuration") {
		t.Fatal("Footer comment not preserved")
	}
}

func TestStripMargin_WithYAMLMultiDocument_PreservesDocumentSeparators(t *testing.T) {
	// Given
	input := `
	|---
	|document: 1
	|name: First Document
	|---
	|document: 2
	|name: Second Document
	|items:
	|  - item1
	|  - item2
	|...
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "---\ndocument: 1\nname: First Document\n---\ndocument: 2\nname: Second Document\nitems:\n  - item1\n  - item2\n..."
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify document separators are preserved
	if !containsSubstring(result, "---\ndocument: 1") {
		t.Fatal("Document start separator not preserved")
	}
	if !containsSubstring(result, "---\ndocument: 2") {
		t.Fatal("Document separator between documents not preserved")
	}
	if !containsSubstring(result, "- item2\n...") {
		t.Fatal("Document end separator not preserved")
	}
}

func TestStripMargin_WithMalformedYAML_ReturnsStringContent(t *testing.T) {
	// Given
	input := `
	|name: John Doe
	|age: 30
	|invalid_structure:
	|  - item1
	|    - nested_without_parent
	|  normal_item: value
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "name: John Doe\nage: 30\ninvalid_structure:\n  - item1\n    - nested_without_parent\n  normal_item: value"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Note: YAML validation would require a YAML parser library
	// For now, we just ensure the string content is returned correctly
}

// YAML StripColumn Tests

func TestStripColumn_WithValidYAMLDocument_ReturnsValidYAML(t *testing.T) {
	// Given
	input := `
	|name: John Doe|
	|age: 30|
	|email: john@example.com|
	|active: true|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "name: John Doe\nage: 30\nemail: john@example.com\nactive: true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify YAML structure is preserved
	if !containsSubstring(result, "name: John Doe") {
		t.Fatal("YAML key-value pair not preserved")
	}
	if !containsSubstring(result, "active: true") {
		t.Fatal("YAML boolean value not preserved")
	}
}

func TestStripColumn_WithYAMLNestedStructure_PreservesIndentation(t *testing.T) {
	// Given
	input := `
	|person:|
	|  name: Jane Smith|
	|  address:|
	|    street: 123 Main St|
	|    city: Anytown|
	|    zipcode: "12345"|
	|  hobbies:|
	|    - reading|
	|    - swimming|
	|    - coding|
	|active: true|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "person:\n  name: Jane Smith\n  address:\n    street: 123 Main St\n    city: Anytown\n    zipcode: \"12345\"\n  hobbies:\n    - reading\n    - swimming\n    - coding\nactive: true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify nested structure indentation is preserved
	lines := strings.Split(result, "\n")
	if !strings.HasPrefix(lines[1], "  name:") {
		t.Fatal("First level indentation not preserved")
	}
	if !strings.HasPrefix(lines[3], "    street:") {
		t.Fatal("Second level indentation not preserved")
	}
}

func TestStripColumn_WithYAMLArrays_PreservesArrayStructure(t *testing.T) {
	// Given
	input := `
	|fruits:|
	|  - apple|
	|  - banana|
	|  - cherry|
	|numbers:|
	|  - 1|
	|  - 2|
	|  - 3.14|
	|mixed_array:|
	|  - "string"|
	|  - 42|
	|  - true|
	|  - null|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "fruits:\n  - apple\n  - banana\n  - cherry\nnumbers:\n  - 1\n  - 2\n  - 3.14\nmixed_array:\n  - \"string\"\n  - 42\n  - true\n  - null"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify array structure is preserved
	if !containsSubstring(result, "fruits:\n  - apple") {
		t.Fatal("Array structure not preserved")
	}
	if !containsSubstring(result, "- 3.14") {
		t.Fatal("Float value in array not preserved")
	}
}

func TestStripColumn_WithYAMLBlockScalars_PreservesBlockScalars(t *testing.T) {
	// Given
	input := `
	|description: ||
	|  This is a literal|
	|  block scalar that|
	|  preserves newlines|
	|summary: >|
	|  This is a folded|
	|  block scalar that|
	|  becomes one line|
	|code: ||
	|  function example() {|
	|    return "Hello World";|
	|  }|
	`

	// When
	result := text.StripColumn(input)

	// Then
	// Verify block scalar markers are preserved
	if !containsSubstring(result, "description: |\n  This is a literal") {
		t.Fatal("Literal block scalar not preserved")
	}
	if !containsSubstring(result, "summary: >\n  This is a folded") {
		t.Fatal("Folded block scalar not preserved")
	}
	if !containsSubstring(result, "code: |\n  function example()") {
		t.Fatal("Code block scalar not preserved")
	}
}

func TestStripColumn_WithYAMLSpecialCharacters_PreservesSpecialChars(t *testing.T) {
	// Given
	input := `
	|company: "AT&T Corporation"|
	|formula: "5 < 10 && 10 > 5"|
	|quote: 'He said "Hello" to me'|
	|path: "C:\\Users\\John\\Documents"|
	|special: "Symbols: !@#$%^&*()"|
	`

	// When
	result := text.StripColumn(input)

	// Then
	// Verify special characters are preserved
	if !containsSubstring(result, `company: "AT&T Corporation"`) {
		t.Fatal("Special characters in quoted string not preserved")
	}
	if !containsSubstring(result, `formula: "5 < 10 && 10 > 5"`) {
		t.Fatal("Formula with special characters not preserved")
	}
	if !containsSubstring(result, `quote: 'He said "Hello" to me'`) {
		t.Fatal("Mixed quotes not preserved")
	}
	if !containsSubstring(result, `special: "Symbols: !@#$%^&*()"`) {
		t.Fatal("Special symbols not preserved")
	}
}

func TestStripColumn_WithYAMLUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	input := `
	|languages:|
	|  chinese: 你好世界|
	|  arabic: مرحبا بالعالم|
	|  japanese: こんにちは世界|
	|emojis:|
	|  simple: 🚀🌟🎉|
	|  complex: 👨‍💻👩‍🚀🏳️‍🌈|
	`

	// When
	result := text.StripColumn(input)

	// Then
	// Verify Unicode content is preserved
	if !containsSubstring(result, "chinese: 你好世界") {
		t.Fatal("Chinese Unicode not preserved")
	}
	if !containsSubstring(result, "arabic: مرحبا بالعالم") {
		t.Fatal("Arabic Unicode not preserved")
	}
	if !containsSubstring(result, "simple: 🚀🌟🎉") {
		t.Fatal("Emoji Unicode not preserved")
	}
	if !containsSubstring(result, "complex: 👨‍💻👩‍🚀🏳️‍🌈") {
		t.Fatal("Complex emoji Unicode not preserved")
	}
}

func TestStripColumn_WithYAMLComments_PreservesComments(t *testing.T) {
	// Given
	input := `
	|# Application configuration|
	|app:|
	|  name: MyApp  # Application name|
	|  version: 1.0.0|
	|  # Database settings|
	|  database:|
	|    host: localhost  # Default host|
	|# End of configuration|
	`

	// When
	result := text.StripColumn(input)

	// Then
	// Verify comments are preserved
	if !containsSubstring(result, "# Application configuration") {
		t.Fatal("Header comment not preserved")
	}
	if !containsSubstring(result, "name: MyApp  # Application name") {
		t.Fatal("Inline comment not preserved")
	}
	if !containsSubstring(result, "  # Database settings") {
		t.Fatal("Indented comment not preserved")
	}
}

func TestStripColumn_WithYAMLMultiDocument_PreservesDocumentSeparators(t *testing.T) {
	// Given
	input := `
	|---|
	|document: 1|
	|name: First Document|
	|---|
	|document: 2|
	|name: Second Document|
	|...|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "---\ndocument: 1\nname: First Document\n---\ndocument: 2\nname: Second Document\n..."
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify document separators are preserved
	if !containsSubstring(result, "---\ndocument: 1") {
		t.Fatal("Document start separator not preserved")
	}
	if !containsSubstring(result, "---\ndocument: 2") {
		t.Fatal("Document separator between documents not preserved")
	}
}

func TestStripColumn_WithCompactYAML_PreservesCompactFormat(t *testing.T) {
	// Given
	input := `|name: John, age: 30, active: true|`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "name: John, age: 30, active: true"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify compact format is preserved
	if !containsSubstring(result, "name: John, age: 30, active: true") {
		t.Fatal("Compact YAML format not preserved")
	}
}

func TestStripColumn_WithMalformedYAML_ReturnsStringContent(t *testing.T) {
	// Given
	input := `
	|name: John Doe|
	|age: 30|
	|invalid_structure:|
	|  - item1|
	|    - nested_without_parent|
	|  normal_item: value|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "name: John Doe\nage: 30\ninvalid_structure:\n  - item1\n    - nested_without_parent\n  normal_item: value"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Note: YAML validation would require a YAML parser library
	// For now, we just ensure the string content is returned correctly
}
