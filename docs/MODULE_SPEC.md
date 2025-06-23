# Textsmith Module Specification

## 1. Module Overview
- **Module Name:** textsmith
- **Repository:** github.com/shapestone/textsmith
- **Package Path:** github.com/shapestone/textsmith/pkg/text
- **Type:** Utility Library (zero external dependencies)

## 2. Responsibilities
- **Primary Responsibilities:**
    - Text processing utilities for Go applications with regex-based implementation
    - Cross-platform text normalization and formatting
    - Visual diff generation for testing and debugging workflows
- **Scope:**
    - **Does:** In-memory string processing, regex-based text manipulation, visual comparison output
    - **Does NOT:** File I/O, network operations, persistent storage, or external service integration

## 3. Architecture & Implementation Details

### Core Functions
- **`func StripMargin(s string) string`**
    - **Regex Pattern:** `(?m)^[ \t]*\|(.*)(?:\r?\n|$)`
    - **Processing:** Removes leading whitespace + pipe, preserves content after pipe
    - **Edge Cases:** Empty lines become empty, malformed lines ignored

- **`func StripColumn(s string) string`**
    - **Regex Pattern:** `(?m)^[ \t]*\|(.*)(?:\|[ \t]*\n|\|[ \t]*$)`
    - **Processing:** Extracts content between enclosing pipes
    - **Edge Cases:** Lines without proper column format are ignored

- **`func Diff(expected string, actual string) (string, bool)`**
    - **Algorithm:** Line-by-line comparison with Unicode symbol rendering
    - **Output:** Formatted table with difference indicators (≠, △, ←, →, ␉, ␣, ␤)
    - **Normalization:** Automatic line ending conversion (CRLF/CR → LF)

### Performance Characteristics
- **Time Complexity:** O(n) where n = input string length
- **Memory Usage:** Minimal allocation with efficient string building
- **Regex Compilation:** Per-function call (not cached)
- **Large Input Handling:** No streaming; processes entire string in memory

## 4. Dependencies & Integration
- **External Dependencies:** Go standard library only (`regexp`, `strings`, `unicode/utf8`)
- **Integration Pattern:** Direct function imports - no initialization or configuration required
- **Error Handling:** Silent failure mode - malformed input lines are ignored, no panics or exceptions

## 5. Code Structure & File Organization
```
textsmith/
└── pkg/text/
    ├── strip_margin.go      # StripMargin and StripColumn implementation
    ├── text_diff.go         # Diff implementation + Unicode symbols
    ├── strip_margin_test.go # Tests for StripMargin and StripColumn
    └── text_diff_test.go    # Tests for Diff
```

**Key Implementation Details:**
- **Regex Compilation:** Not cached - compiled on each function call
- **Memory Management:** No pooling; relies on Go GC for string cleanup
- **Thread Safety:** Pure functions - safe for concurrent use
- **Unicode Handling:** Full UTF-8 support with proper character boundary detection

## 6. Testing & Quality Strategy
- **Coverage Requirement:** 100% test coverage maintained
- **Test Types:** Table-driven tests, edge cases, cross-platform compatibility, benchmarks
- **CI Pipeline:** GitHub Actions with Go 1.21+ matrix testing
- **Quality Gates:** `make check` runs format + vet + test before commits

**Test Categories:**
- Functional correctness (happy path + edge cases)
- Cross-platform line ending compatibility
- Unicode/emoji handling validation
- Performance benchmarks for large inputs
- Memory allocation profiling

## 7. Release & Deployment Strategy
- **Versioning:** Git tags with semantic versioning (`v1.0.0`)
- **Distribution:** Go module system automatic detection
- **CI/CD:** GitHub Actions pipeline (`.github/workflows/go.yml`)
- **Backward Compatibility:** Strong commitment within major versions

**Release Process:**
```bash
git tag v1.0.0
git push origin v1.0.0
# Go module system handles the rest
```

## 8. Performance Considerations & Limitations

### Current Limitations
- **Regex Compilation:** No caching - compiles pattern on each function call
- **Memory Usage:** No streaming support - entire input processed in memory
- **Large Input Performance:** Not optimized for multi-gigabyte text processing
- **Concurrency:** No built-in parallelization for multi-core systems

### Recommended Usage Patterns
- **Optimal:** Small to medium text processing (< 10MB)
- **Acceptable:** Configuration files, code generation, test output comparison
- **Avoid:** Real-time processing of very large files, high-frequency operations

### Performance Benchmarks
Run `make bench` to measure:
- Processing speed for different input sizes
- Memory allocation patterns
- Regex compilation overhead

## 9. Future Architecture Considerations
- **Regex Caching:** Pre-compile and cache regex patterns for repeated use
- **Streaming API:** Support for `io.Reader`/`io.Writer` interfaces for large files
- **Plugin Architecture:** Extensible processing pipeline for custom transformations
- **Parallel Processing:** Multi-core support for large text processing
- **Custom Formatters:** Pluggable output formats (JSON, XML, HTML) for diff results
- **Configuration Options:** Runtime options for margin characters, diff symbols, etc.

## 10. Integration Patterns

### Testing Framework Integration
```
// Example: Using Diff in test assertions
func TestMyFunction(t *testing.T) {
    expected := generateExpectedOutput()
    actual := myFunction(input)
    
    if diff, match := text.Diff(expected, actual); !match {
        t.Errorf("Output mismatch:\n%s", diff)
    }
}
```

### Code Generation Integration
```
// Example: Using StripMargin for templates
func generateCode(className string) string {
    return text.StripMargin(fmt.Sprintf(`
        |package main
        |
        |type %s struct {
        |    ID   int    ` + "`json:\"id\"`" + `
        |    Name string ` + "`json:\"name\"`" + `
        |}`, className))
}
```

### Configuration Processing
```
// Example: Embedded configuration with StripColumn
func getDefaultConfig() string {
    return text.StripColumn(`
        |server.host = localhost|
        |server.port = 8080|
        |db.url = postgres://localhost/mydb|
    `)
}
```