package text_test

import (
	"encoding/json"
	"github.com/shapestone/textsmith/pkg/text"
	"testing"
)

// JSON StripMargin Tests

func TestStripMargin_WithValidJSONObject_ReturnsValidJSON(t *testing.T) {
	// Given
	input := `
	|{
	|  "name": "John Doe",
	|  "age": 30,
	|  "email": "john@example.com"
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "{\n  \"name\": \"John Doe\",\n  \"age\": 30,\n  \"email\": \"john@example.com\"\n}"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid JSON
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestStripMargin_WithJSONArray_ReturnsValidJSONArray(t *testing.T) {
	// Given
	input := `
	|[
	|  "apple",
	|  "banana",
	|  "cherry"
	|]
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "[\n  \"apple\",\n  \"banana\",\n  \"cherry\"\n]"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid JSON
	var jsonArray []string
	if err := json.Unmarshal([]byte(result), &jsonArray); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestStripMargin_WithNestedJSONStructure_ReturnsValidNestedJSON(t *testing.T) {
	// Given
	input := `
	|{
	|  "person": {
	|    "name": "Jane Smith",
	|    "address": {
	|      "street": "123 Main St",
	|      "city": "Anytown",
	|      "zipcode": "12345"
	|    },
	|    "hobbies": ["reading", "swimming", "coding"]
	|  },
	|  "active": true
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify nested structure exists
	person, ok := jsonObj["person"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'person' to be an object")
	}
	if person["name"] != "Jane Smith" {
		t.Fatalf("Expected person name to be 'Jane Smith', got %v", person["name"])
	}
}

func TestStripMargin_WithJSONSpecialCharacters_PreservesEscaping(t *testing.T) {
	// Given
	input := `
	|{
	|  "quote": "He said \"Hello, world!\"",
	|  "backslash": "C:\\Users\\John",
	|  "newline": "Line 1\nLine 2",
	|  "tab": "Column1\tColumn2",
	|  "unicode": "Emoji: \uD83D\uDE80"
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify special characters are preserved
	if jsonObj["quote"] != "He said \"Hello, world!\"" {
		t.Fatalf("Quote not preserved correctly: %v", jsonObj["quote"])
	}
}

func TestStripMargin_WithJSONUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	input := `
	|{
	|  "chinese": "你好世界",
	|  "arabic": "مرحبا بالعالم",
	|  "emoji": "🚀🌟🎉",
	|  "complex_emoji": "👨‍💻👩‍🚀🏳️‍🌈"
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify Unicode content is preserved
	if jsonObj["chinese"] != "你好世界" {
		t.Fatalf("Chinese text not preserved: %v", jsonObj["chinese"])
	}
	if jsonObj["emoji"] != "🚀🌟🎉" {
		t.Fatalf("Emoji not preserved: %v", jsonObj["emoji"])
	}
}

func TestStripMargin_WithJSONNumbers_PreservesNumberTypes(t *testing.T) {
	// Given
	input := `
	|{
	|  "integer": 42,
	|  "float": 3.14159,
	|  "negative": -273.15,
	|  "scientific": 1.23e-4,
	|  "large": 1234567890123456789
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify numbers are preserved correctly
	if jsonObj["integer"].(float64) != 42 {
		t.Fatalf("Integer not preserved: %v", jsonObj["integer"])
	}
	if jsonObj["float"].(float64) != 3.14159 {
		t.Fatalf("Float not preserved: %v", jsonObj["float"])
	}
}

func TestStripMargin_WithJSONBooleanAndNull_PreservesTypes(t *testing.T) {
	// Given
	input := `
	|{
	|  "active": true,
	|  "disabled": false,
	|  "value": null,
	|  "empty": ""
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify boolean and null values
	if jsonObj["active"] != true {
		t.Fatalf("Boolean true not preserved: %v", jsonObj["active"])
	}
	if jsonObj["disabled"] != false {
		t.Fatalf("Boolean false not preserved: %v", jsonObj["disabled"])
	}
	if jsonObj["value"] != nil {
		t.Fatalf("Null not preserved: %v", jsonObj["value"])
	}
}

func TestStripMargin_WithMalformedJSON_ReturnsStringContent(t *testing.T) {
	// Given
	input := `
	|{
	|  "name": "John",
	|  "age": 30,
	|  "missing_quote: "value"
	|}
	`

	// When
	result := text.StripMargin(input)

	// Then
	expected := "{\n  \"name\": \"John\",\n  \"age\": 30,\n  \"missing_quote: \"value\"\n}"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's invalid JSON (should fail to parse)
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err == nil {
		t.Fatal("Expected invalid JSON to fail parsing, but it succeeded")
	}
}

// JSON StripColumn Tests

func TestStripColumn_WithValidJSONObject_ReturnsValidJSON(t *testing.T) {
	// Given
	input := `
	|{|
	|  "name": "John Doe",|
	|  "age": 30,|
	|  "email": "john@example.com"|
	|}|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "{\n  \"name\": \"John Doe\",\n  \"age\": 30,\n  \"email\": \"john@example.com\"\n}"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid JSON
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestStripColumn_WithJSONArray_ReturnsValidJSONArray(t *testing.T) {
	// Given
	input := `
	|[|
	|  "apple",|
	|  "banana",|
	|  "cherry"|
	|]|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "[\n  \"apple\",\n  \"banana\",\n  \"cherry\"\n]"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid JSON
	var jsonArray []string
	if err := json.Unmarshal([]byte(result), &jsonArray); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestStripColumn_WithNestedJSONStructure_ReturnsValidNestedJSON(t *testing.T) {
	// Given
	input := `
	|{|
	|  "person": {|
	|    "name": "Jane Smith",|
	|    "address": {|
	|      "street": "123 Main St",|
	|      "city": "Anytown",|
	|      "zipcode": "12345"|
	|    },|
	|    "hobbies": ["reading", "swimming", "coding"]|
	|  },|
	|  "active": true|
	|}|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify nested structure exists
	person, ok := jsonObj["person"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'person' to be an object")
	}
	if person["name"] != "Jane Smith" {
		t.Fatalf("Expected person name to be 'Jane Smith', got %v", person["name"])
	}
}

func TestStripColumn_WithJSONSpecialCharacters_PreservesEscaping(t *testing.T) {
	// Given
	input := `
	|{|
	|  "quote": "He said \"Hello, world!\"",|
	|  "backslash": "C:\\Users\\John",|
	|  "newline": "Line 1\nLine 2",|
	|  "tab": "Column1\tColumn2",|
	|  "unicode": "Emoji: \uD83D\uDE80"|
	|}|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify special characters are preserved
	if jsonObj["quote"] != "He said \"Hello, world!\"" {
		t.Fatalf("Quote not preserved correctly: %v", jsonObj["quote"])
	}
}

func TestStripColumn_WithJSONUnicodeContent_PreservesUnicode(t *testing.T) {
	// Given
	input := `
	|{|
	|  "chinese": "你好世界",|
	|  "arabic": "مرحبا بالعالم",|
	|  "emoji": "🚀🌟🎉",|
	|  "complex_emoji": "👨‍💻👩‍🚀🏳️‍🌈"|
	|}|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify Unicode content is preserved
	if jsonObj["chinese"] != "你好世界" {
		t.Fatalf("Chinese text not preserved: %v", jsonObj["chinese"])
	}
	if jsonObj["emoji"] != "🚀🌟🎉" {
		t.Fatalf("Emoji not preserved: %v", jsonObj["emoji"])
	}
}

func TestStripColumn_WithCompactJSON_PreservesCompactFormat(t *testing.T) {
	// Given
	input := `|{"name":"John","age":30,"active":true}|`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "{\"name\":\"John\",\"age\":30,\"active\":true}"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's valid JSON
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}
}

func TestStripColumn_WithEmptyJSONStructures_ReturnsValidEmptyStructures(t *testing.T) {
	// Given
	input := `
	|{|
	|"empty_object": {},|
	|"empty_array": [],|
	|"empty_string": ""|
	|}|
	`

	// When
	result := text.StripColumn(input)

	// Then
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
		t.Fatalf("Result is not valid JSON: %v", err)
	}

	// Verify empty structures
	emptyObj, ok := jsonObj["empty_object"].(map[string]interface{})
	if !ok || len(emptyObj) != 0 {
		t.Fatalf("Expected empty object, got %v", jsonObj["empty_object"])
	}

	emptyArray, ok := jsonObj["empty_array"].([]interface{})
	if !ok || len(emptyArray) != 0 {
		t.Fatalf("Expected empty array, got %v", jsonObj["empty_array"])
	}
}

func TestStripColumn_WithMalformedJSON_ReturnsStringContent(t *testing.T) {
	// Given
	input := `
	|{|
	|  "name": "John",|
	|  "age": 30,|
	|  "missing_quote: "value"|
	|}|
	`

	// When
	result := text.StripColumn(input)

	// Then
	expected := "{\n  \"name\": \"John\",\n  \"age\": 30,\n  \"missing_quote: \"value\"\n}"
	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}

	// Verify it's invalid JSON (should fail to parse)
	var jsonObj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &jsonObj); err == nil {
		t.Fatal("Expected invalid JSON to fail parsing, but it succeeded")
	}
}
