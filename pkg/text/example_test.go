package text_test

import (
	"fmt"
	"github.com/shapestone/textsmith/pkg/text"
)

// ExampleStripMargin demonstrates basic usage of StripMargin
func ExampleStripMargin() {
	content := text.StripMargin(`
		|func Example() {
		|    fmt.Println("Hello, World!")
		|    return nil
		|}`)

	fmt.Print(content)
	// Output:
	// func Example() {
	//     fmt.Println("Hello, World!")
	//     return nil
	// }
}

// ExampleStripMargin_sql demonstrates SQL query formatting
func ExampleStripMargin_sql() {
	query := text.StripMargin(`
		|SELECT u.name, u.email
		|FROM users u
		|WHERE u.active = true
		|ORDER BY u.created_at DESC`)

	fmt.Print(query)
	// Output:
	// SELECT u.name, u.email
	// FROM users u
	// WHERE u.active = true
	// ORDER BY u.created_at DESC
}

// ExampleStripMargin_json demonstrates JSON formatting
func ExampleStripMargin_json() {
	jsonContent := text.StripMargin(`
		|{
		|  "name": "John",
		|  "age": 30,
		|  "city": "New York"
		|}`)

	fmt.Print(jsonContent)
	// Output:
	// {
	//   "name": "John",
	//   "age": 30,
	//   "city": "New York"
	// }
}

// ExampleStripColumn demonstrates basic usage of StripColumn
func ExampleStripColumn() {
	content := text.StripColumn(`
		|func Example() {|
		|    fmt.Println("Hello")|
		|    return nil|
		|}|`)

	fmt.Print(content)
	// Output:
	// func Example() {
	//     fmt.Println("Hello")
	//     return nil
	// }
}

// ExampleStripColumn_table demonstrates table-like formatting
func ExampleStripColumn_table() {
	config := text.StripColumn(`
		|server.host = localhost|
		|server.port = 8080|
		|database.url = postgres://...|`)

	fmt.Print(config)
	// Output:
	// server.host = localhost
	// server.port = 8080
	// database.url = postgres://...
}

// ExampleDiff demonstrates basic usage comparing identical strings
func ExampleDiff() {
	expected := "hello world"
	actual := "hello world"

	_, match := text.Diff(expected, actual, false)
	fmt.Printf("Strings match: %t\n", match)
	// Output:
	// Strings match: true
}

// ExampleDiff_different demonstrates detecting differences
func ExampleDiff_different() {
	expected := "hello"
	actual := "help!"

	diff, match := text.Diff(expected, actual, false)
	fmt.Printf("Strings match: %t\n", match)
	// Diff output shows the difference with ≠ symbol
	fmt.Printf("Contains difference marker: %t\n", len(diff) > 0)
	// Output:
	// Strings match: false
	// Contains difference marker: true
}

// ExampleCompareStrings demonstrates comparing identical strings
func ExampleCompareStrings() {
	actual := "hello world"
	expected := "hello world"

	result := text.CompareStrings(actual, expected)
	fmt.Print(result)
	// Output:
	// CompareStrings: ✓ [MATCH]
	//   Expected: "hello␣world"¶
	//   Actual:   "hello␣world"¶
}

// ExampleCompareStrings_different demonstrates comparing different strings
func ExampleCompareStrings_different() {
	actual := "hello world"
	expected := "hello mars"

	result := text.CompareStrings(actual, expected)
	fmt.Print(result)
	// Output:
	// CompareStrings: ✗ [ASSERTION_FAILED]
	// - Expected: "hello␣mars"¶
	// + Actual:   "hello␣world"¶
	//
	//   Difference at position 6:
	//       Expected character: 'm' (U+006D)
	//       Actual character:   'w' (U+0077)
}

// ExampleCompareStrings_whitespace demonstrates whitespace difference detection
func ExampleCompareStrings_whitespace() {
	actual := "hello\tworld"
	expected := "hello world"

	result := text.CompareStrings(actual, expected)
	fmt.Print(result)
	// Output:
	// CompareStrings: ✗ [ASSERTION_FAILED]
	// - Expected: "hello␣world"¶
	// + Actual:   "hello␉world"¶
	//
	//   Difference at position 5:
	//       Expected character: ' ' (U+0020)
	//       Actual character:   '\t' (U+0009)
}

// ExampleCompareStrings_empty demonstrates empty string handling
func ExampleCompareStrings_empty() {
	actual := ""
	expected := ""

	result := text.CompareStrings(actual, expected)
	fmt.Print(result)
	// Output:
	// CompareStrings: ✓ [MATCH]
	//   Expected: <empty>¶
	//   Actual:   <empty>¶
}

// ExampleCompareStringsRaw demonstrates raw comparison without visualization
func ExampleCompareStringsRaw() {
	actual := "hello world"
	expected := "hello mars"

	result := text.CompareStringsRaw(actual, expected)
	fmt.Print(result)
	// Output:
	// CompareStrings: ✗ [ASSERTION_FAILED]
	// - Expected: "hello mars"¶
	// + Actual:   "hello world"¶
	//
	//   Difference at position 6:
	//       Expected character: 'm' (U+006D)
	//       Actual character:   'w' (U+0077)
}

// ExampleCompareStringsRaw_whitespace demonstrates raw whitespace handling
func ExampleCompareStringsRaw_whitespace() {
	actual := "hello\tworld"
	expected := "hello world"

	result := text.CompareStringsRaw(actual, expected)
	fmt.Print(result)
	// Output:
	// CompareStrings: ✗ [ASSERTION_FAILED]
	// - Expected: "hello world"¶
	// + Actual:   "hello	world"¶
	//
	//   Difference at position 5:
	//       Expected character: ' ' (U+0020)
	//       Actual character:   '\t' (U+0009)
}
