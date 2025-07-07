# textsmith

Textsmith is a lightweight utility library for crafting, comparing, and transforming text. It provides composable functions for common text operations such as diffing, margin stripping, formatting, and more — designed to help developers work with plain text more effectively.

## Installation

### As a library
```bash
go get github.com/shapestone/textsmith
```

## Usage

### Library Usage

Import the textsmith package in your Go projects:

```
import "github.com/shapestone/textsmith/pkg/text"
```

#### StripMargin Function

The `StripMargin` function lets you define multiline strings where each line is prepended with optional whitespace and a pipeline symbol. This is particularly useful for embedding formatted text, code, or configuration in your Go source.

```
package main

import (
	"fmt"
	"github.com/shapestone/textsmith/pkg/text"
)

func main() {
	content := text.StripMargin(`
        |func Example() {
        |    fmt.Println("Hello, World!")
        |    return nil
        |}`)

	fmt.Print(content)
}
```

**Output:**
```
func Example() {
    fmt.Println("Hello, World!")
    return nil
}
```

##### Examples

**Basic multiline text:**
```
message := text.StripMargin(`
    |Welcome to our application!
    |Please follow these steps:
    |1. Login with your credentials
    |2. Navigate to the dashboard
    |3. Start using the features`)

fmt.Print(message)
```

**Output:**
```
Welcome to our application!
Please follow these steps:
1. Login with your credentials
2. Navigate to the dashboard
3. Start using the features
```

**Code templates:**
```
sqlQuery := text.StripMargin(`
    |SELECT u.name, u.email, p.title
    |FROM users u
    |JOIN profiles p ON u.id = p.user_id
    |WHERE u.active = true
    |ORDER BY u.created_at DESC`)

fmt.Print(sqlQuery)
```

**Output:**
```sql
SELECT u.name, u.email, p.title
FROM users u
JOIN profiles p ON u.id = p.user_id
WHERE u.active = true
ORDER BY u.created_at DESC
```

**Configuration or markup:**
```
yamlConfig := text.StripMargin(`
    |database:
    |  host: localhost
    |  port: 5432
    |  name: myapp
    |logging:
    |  level: info
    |  file: /var/log/app.log`)

fmt.Print(yamlConfig)
```

**Output:**
```yaml
database:
  host: localhost
  port: 5432
  name: myapp
logging:
  level: info
  file: /var/log/app.log
```

#### StripColumn Function

The `StripColumn` function lets you define multiline strings where each line is enclosed by pipeline symbols on both sides. This provides a more structured approach to multiline text with clear boundaries.

```
content := text.StripColumn(`
    |func Example() {|
    |    fmt.Println("Hello, World!")|
    |    return nil|
    |}|`)

fmt.Print(content)
```

**Output:**
```
func Example() {
    fmt.Println("Hello, World!")
    return nil
}
```

##### StripColumn vs StripMargin

- **StripMargin**: Uses opening pipe only (`|content`)
- **StripColumn**: Uses enclosing pipes (`|content|`)
- **StripColumn** requires both opening and closing pipes on each line
- Lines without proper column format are ignored

##### StripColumn Examples

**Structured templates:**
```
template := text.StripColumn(`
    |<!DOCTYPE html>|
    |<html>|
    |  <head><title>{{.Title}}</title></head>|
    |  <body>{{.Content}}</body>|
    |</html>|`)

fmt.Print(template)
```

**Table-like configuration:**
```
config := text.StripColumn(`
    |server.host     = localhost|
    |server.port     = 8080|
    |database.url    = postgres://...|
    |logging.level   = info|`)

fmt.Print(config)
```

#### Diff Function

The `Diff` function compares two strings and produces a visual side-by-side diff output, making it easy to spot differences between expected and actual text. It returns both the formatted diff string and a boolean indicating whether the strings match.

**Function signature:**
```
func Diff(expected string, actual string) (string, bool)
```

**Basic usage:**
```
package main

import (
    "fmt"
    "github.com/shapestone/textsmith/pkg/text"
)

func main() {
    expected := "hello world"
    actual := "hello universe"
    
    diff, match := text.Diff(expected, actual)
    fmt.Printf("Strings match: %t\n", match)
    fmt.Print(diff)
}
```

##### Diff Examples

**Identical strings:**
```
expected := text.StripMargin(`
    |line 1
    |line 2
    |line 3`)
actual := text.StripMargin(`
    |line 1
    |line 2
    |line 3`)

diff, match := text.Diff(expected, actual)
fmt.Printf("Match: %t\n", match)
fmt.Print(diff)
```

**Output:**
```
Match: true
Expected | Actual
-------- | ------
line 1   | line 1
line 2   | line 2
line 3   | line 3
```

**Different content:**
```
expected := "hello world"
actual := "hello universe"

diff, match := text.Diff(expected, actual)
fmt.Printf("Match: %t\n", match)
fmt.Print(diff)
```

**Output:**
```
Match: false
Expected      | Actual
------------- | -------------
hello world   ≠ hello universe
      △       |       △
```

**Different line counts:**
```
expected := text.StripMargin(`
    |line 1
    |line 2`)
actual := text.StripMargin(`
    |line 1
    |line 2
    |line 3`)

diff, match := text.Diff(expected, actual)
fmt.Printf("Match: %t\n", match)
fmt.Print(diff)
```

**Output:**
```
Match: false
Expected | Actual
-------- | --------
line 1   | line 1
line 2   | line 2
       → | line 3
```

**Whitespace differences:**
```
expected := "hello\tworld"  // tab character
actual := "hello world"     // space character

diff, match := text.Diff(expected, actual)
fmt.Printf("Match: %t\n", match)
fmt.Print(diff)
```

**Output:**
```
Match: false
Expected    | Actual
----------- | -----------
hello␉world ≠ hello␣world
     △      |      △
```

##### Cross-Platform Line Ending Support

The diff function automatically normalizes different line ending formats:

- **Unix/Linux**: `\n` (LF)
- **Windows**: `\r\n` (CRLF)
- **Classic Mac**: `\r` (CR)

All line endings are converted to Unix format (`\n`) before comparison, ensuring consistent behavior across platforms.

```
unixText := "line1\nline2"
windowsText := "line1\r\nline2"

diff, match := text.Diff(unixText, windowsText)
// match will be true - line endings are normalized
```

##### Diff Symbols

The diff output uses special Unicode symbols to indicate different types of changes:

- **≠** - Lines that differ between expected and actual
- **△** - Points to the exact position where strings start to differ
- **←** - Expected has more content (missing from actual)
- **→** - Actual has more content (extra in actual)
- **␉** - Tab characters (shown when whitespace differs)
- **␣** - Space characters (shown when whitespace differs)
- **␤** - Empty lines (shown when line is empty but significant)

##### How Diff Works

The diff function processes multiline strings by:

1. **Line-by-line comparison** - Splits both strings by newlines and compares each line
2. **Side-by-side layout** - Creates a formatted table with Expected | Actual columns
3. **Whitespace visualization** - Converts invisible characters (tabs, spaces) to visible symbols when differences are found
4. **Precise difference location** - Shows exactly where strings start to differ using the △ symbol
5. **Length difference handling** - Uses arrows to indicate when one string has more lines than the other
6. **Proper alignment** - Ensures consistent column widths for readable output
7. **Cross-platform normalization** - Handles different line ending formats automatically

#### CompareStrings Function

The `CompareStrings` function provides a test framework style comparison between actual and expected strings with detailed diff highlighting. It's specifically designed for testing purposes and converts invisible characters to visible symbols for better debugging.

**Function signatures:**
```go
func CompareStrings(actual, expected string) string
func CompareStringsRaw(actual, expected string) string
```

**Basic usage:**
```go
package main

import (
    "fmt"
    "github.com/shapestone/textsmith/pkg/text"
)

func main() {
    actual := "hello world"
    expected := "hello world"
    
    result := text.CompareStrings(actual, expected)
    fmt.Print(result)
}
```

##### CompareStrings vs CompareStringsRaw

- **CompareStrings**: Converts invisible characters to visible symbols (␣ for space, ␉ for tab, etc.)
- **CompareStringsRaw**: Shows strings as-is without visualization, useful when strings are already formatted

##### CompareStrings Examples

**Matching strings:**
```go
actual := "hello world"
expected := "hello world"

result := text.CompareStrings(actual, expected)
fmt.Print(result)
```

**Output:**
```
CompareStrings: ✓ [MATCH]
  Expected: "hello␣world"¶
  Actual:   "hello␣world"¶
```

**Different strings:**
```go
actual := "hello world"
expected := "hello mars"

result := text.CompareStrings(actual, expected)
fmt.Print(result)
```

**Output:**
```
CompareStrings: ✗ [ASSERTION_FAILED]
- Expected: "hello␣mars"¶
+ Actual:   "hello␣world"¶

  Difference at position 6:
      Expected character: 'm' (U+006D)
      Actual character:   'w' (U+0077)
```

**Whitespace differences:**
```go
actual := "hello\tworld"
expected := "hello world"

result := text.CompareStrings(actual, expected)
fmt.Print(result)
```

**Output:**
```
CompareStrings: ✗ [ASSERTION_FAILED]
- Expected: "hello␣world"¶
+ Actual:   "hello␉world"¶

  Difference at position 5:
      Expected character: ' ' (U+0020)
      Actual character:   '	' (U+0009)
```

**Empty strings:**
```go
actual := ""
expected := ""

result := text.CompareStrings(actual, expected)
fmt.Print(result)
```

**Output:**
```
CompareStrings: ✓ [MATCH]
  Expected: <empty>¶
  Actual:   <empty>¶
```

**Raw comparison (no visualization):**
```go
actual := "hello world"
expected := "hello mars"

result := text.CompareStringsRaw(actual, expected)
fmt.Print(result)
```

**Output:**
```
CompareStrings: ✗ [ASSERTION_FAILED]
- Expected: "hello mars"¶
+ Actual:   "hello world"¶

  Difference at position 6:
      Expected character: 'm' (U+006D)
      Actual character:   'w' (U+0077)
```

##### Visualization Symbols

The `CompareStrings` function uses special Unicode symbols to make invisible characters visible:

- **␣** - Space characters (U+0020)
- **␉** - Tab characters (U+0009)
- **␊** - Line feed (U+000A)
- **␍** - Carriage return (U+000D)
- **␋** - Vertical tab (U+000B)
- **␌** - Form feed (U+000C)
- **¶** - End of line marker
- **<empty>** - Empty string indicator

##### Unicode Support

Both functions fully support Unicode characters including emojis and complex scripts:

```go
actual := "Hello 世界! 🌍"
expected := "Hello 世界! 🌎"

result := text.CompareStrings(actual, expected)
// Shows exact Unicode code points for differences
```

## Building and Testing

### Test
```shell
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests in short mode
make test-short

# Run tests with race detection
make test-race

# Run benchmarks
make bench
```

### Development
```shell
# Run all quality checks (format + vet + test)
make check

# Format code
make fmt-fix

# Run static analysis
make vet

# Tidy module dependencies
make mod-tidy
```

## Features

- **StripMargin**: Clean multiline string handling with margin indicators
- **StripColumn**: Column-based multiline string handling with enclosing pipes
- **Diff**: Visual side-by-side text comparison with precise difference highlighting
- **CompareStrings**: Test framework style string comparison with invisible character visualization
- **Whitespace visualization**: Shows invisible characters when comparing text
- **Cross-platform line endings**: Automatic normalization of Unix, Windows, and Mac line endings
- **Unicode support**: Works with international characters and emojis
- **Performance optimized**: Efficient regex-based processing
- **Comprehensive tests**: Full test coverage with benchmarks

## Library API

### Functions

- `StripMargin(s string) string` - Process multiline strings with margin pipes
- `StripColumn(s string) string` - Process multiline strings with enclosing pipes
- `Diff(expected string, actual string) (string, bool)` - Compare two strings and return visual diff
- `CompareStrings(actual, expected string) string` - Test framework style string comparison with visualization
- `CompareStringsRaw(actual, expected string) string` - String comparison without character visualization

### How StripMargin Works

The function uses a regular expression to find lines that start with optional whitespace followed by a pipe character (`|`). It then:

1. Matches lines with the pattern: `(?m)^[ \t]*\|(.*)(?:\r?\n|$)`
    - `(?m)` enables multiline mode
    - `\r?` handles Windows line endings
2. Extracts the content after the pipe for each line
3. Joins the extracted content with newlines
4. Returns the clean multiline string

**Input processing:**
- Leading tabs and spaces before `|` are removed
- The `|` character itself is removed
- Content after `|` is preserved exactly (including trailing spaces)
- Empty lines (just `|`) become empty lines in output
- Lines without `|` are ignored

### How StripColumn Works

The function uses a regular expression to find lines that are enclosed by pipe characters (`|content|`). It then:

1. Matches lines with the pattern: `(?m)^[ \t]*\|(.*)(?:\|[ \t]*\n|\|[ \t]*$)`
2. Extracts the content between the pipes for each line
3. Joins the extracted content with newlines
4. Returns the clean multiline string

**Input processing:**
- Leading tabs and spaces before opening `|` are removed
- Both opening and closing `|` characters are removed
- Content between pipes is preserved exactly
- Lines without proper `|content|` format are ignored
- Trailing whitespace after closing `|` is ignored

## Testing

The project includes comprehensive tests:

- **Unit tests** (`strip_margin_test.go`, `text_diff_test.go`, `string_diff_test.go`) - Test all public functions
- **Edge case testing** - Empty inputs, unicode content, large inputs, whitespace differences
- **Performance benchmarks** - Ensure efficient processing
- **Table-driven tests** - Comprehensive scenario coverage
- **Cross-platform testing** - Validates line ending normalization

Tests use black box testing to ensure they only test public APIs, providing validation of the user experience.

## Use Cases

### StripMargin Use Cases
- **Code generation**: Template processing for generated code
- **Configuration**: Embedding YAML, JSON, or other config formats
- **Documentation**: Multi-line help text and usage examples
- **SQL queries**: Clean formatting of complex database queries
- **Shell scripts**: Embedding shell commands or scripts
- **Markup**: HTML, XML, or other markup within Go source
- **Test data**: Readable test fixtures and expected outputs

### StripColumn Use Cases
- **Structured templates**: Content that needs clear visual boundaries
- **Table-like data**: When you want symmetric formatting
- **Configuration blocks**: Settings that should be visually contained
- **Code blocks**: When you want clear start/end markers for embedded code
- **Documentation tables**: Content that should appear in columns
- **Test fixtures**: When you need clearly bounded test data

### Diff Use Cases
- **Testing frameworks**: Showing detailed differences in test failures
- **Code review tools**: Highlighting changes between file versions
- **Configuration validation**: Comparing expected vs actual config outputs
- **API testing**: Comparing expected vs actual JSON/XML responses
- **Documentation**: Showing before/after examples
- **Data validation**: Comparing processed vs expected data formats
- **Debugging**: Understanding discrepancies in text processing

### CompareStrings Use Cases
- **Testing frameworks**: Detailed assertion failure messages with exact difference locations
- **Unit test debugging**: Understanding why string comparisons fail, especially with whitespace
- **Data processing validation**: Ensuring text transformations produce exact expected results
- **Configuration testing**: Validating that config generation produces expected output
- **Template testing**: Verifying template rendering produces exact expected content
- **Protoc

## Future Considerations
- No special support for extremely long lines that might need wrapping