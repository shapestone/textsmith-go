package text_test

import (
	"github.com/shapestone/textsmith/pkg/text"
	"strings"
	"testing"
)

// TestIntegration_TestFrameworkUsage tests real-world test framework usage
func TestIntegration_TestFrameworkUsage(t *testing.T) {
	// Simulate a test framework comparing JSON output
	expected := text.StripMargin(`
		|{
		|  "name": "John Doe",
		|  "age": 30,
		|  "email": "john@example.com"
		|}`)

	actual := text.StripMargin(`
		|{
		|  "name": "John Doe",
		|  "age": 31,
		|  "email": "john@example.com"
		|}`)

	// Use CompareStrings for detailed failure message
	result := text.CompareStrings(actual, expected)

	if !strings.Contains(result, "✗ [ASSERTION_FAILED]") {
		t.Errorf("Expected assertion to fail for mismatched JSON")
	}

	if !strings.Contains(result, "Difference at position") {
		t.Errorf("Expected difference position to be shown")
	}
}

// TestIntegration_SQLQueryFormatting tests SQL query formatting scenario
func TestIntegration_SQLQueryFormatting(t *testing.T) {
	// Real-world scenario: formatting SQL queries
	query := text.StripMargin(`
		|SELECT
		|    u.id,
		|    u.name,
		|    u.email,
		|    p.title AS profile_title
		|FROM users u
		|INNER JOIN profiles p ON u.id = p.user_id
		|WHERE u.active = true
		|    AND u.created_at > '2024-01-01'
		|ORDER BY u.created_at DESC
		|LIMIT 100`)

	// Verify query is properly formatted
	if !strings.Contains(query, "SELECT") {
		t.Errorf("Expected query to contain SELECT")
	}

	if !strings.Contains(query, "INNER JOIN") {
		t.Errorf("Expected query to contain INNER JOIN")
	}

	if strings.HasPrefix(query, " ") || strings.HasPrefix(query, "\t") {
		t.Errorf("Expected query to not have leading whitespace")
	}

	// Count lines - should have proper structure
	lines := strings.Split(query, "\n")
	if len(lines) < 9 {
		t.Errorf("Expected query to have at least 9 lines, got %d", len(lines))
	}
}

// TestIntegration_ConfigFileGeneration tests configuration file generation
func TestIntegration_ConfigFileGeneration(t *testing.T) {
	// Real-world scenario: generating YAML config
	config := text.StripMargin(`
		|server:
		|  host: localhost
		|  port: 8080
		|  timeout: 30s
		|database:
		|  host: db.example.com
		|  port: 5432
		|  name: myapp_production
		|  pool:
		|    min: 5
		|    max: 20
		|logging:
		|  level: info
		|  output: /var/log/app.log`)

	// Verify config structure
	if !strings.Contains(config, "server:") {
		t.Errorf("Expected config to contain server section")
	}

	if !strings.Contains(config, "database:") {
		t.Errorf("Expected config to contain database section")
	}

	// Verify indentation is preserved
	if !strings.Contains(config, "  host: localhost") {
		t.Errorf("Expected proper indentation for nested properties")
	}
}

// TestIntegration_DiffInTestFailure tests using Diff in test failure messages
func TestIntegration_DiffInTestFailure(t *testing.T) {
	// Real-world scenario: test assertion failure with detailed diff
	expected := "Line 1\nLine 2\nLine 3\nLine 4"
	actual := "Line 1\nLine 2\nLine 5\nLine 4"

	diff, match := text.Diff(expected, actual)

	if match {
		t.Errorf("Expected strings to not match")
	}

	// Diff should show where they differ (stops at first difference)
	// Note: Diff visualizes spaces with ␣ symbol
	if !strings.Contains(diff, "Line␣3") {
		t.Errorf("Expected diff to show expected line, got:\n%s", diff)
	}

	if !strings.Contains(diff, "Line␣5") {
		t.Errorf("Expected diff to show actual line, got:\n%s", diff)
	}

	// Should have visual indicators
	if !strings.Contains(diff, "≠") {
		t.Errorf("Expected diff to contain difference indicator")
	}
}

// TestIntegration_CodeGeneration tests code generation scenario
func TestIntegration_CodeGeneration(t *testing.T) {
	// Real-world scenario: generating Go code
	generatedCode := text.StripColumn(`
		|package main|
		||
		|import (|
		|	"fmt"|
		|	"time"|
		|)|
		||
		|func main() {|
		|	fmt.Println("Hello, World!")|
		|	fmt.Println("Current time:", time.Now())|
		|}|`)

	// Verify code structure
	if !strings.Contains(generatedCode, "package main") {
		t.Errorf("Expected generated code to have package declaration")
	}

	if !strings.Contains(generatedCode, "import (") {
		t.Errorf("Expected generated code to have imports")
	}

	if !strings.Contains(generatedCode, "func main()") {
		t.Errorf("Expected generated code to have main function")
	}

	// Verify proper formatting (no leading whitespace)
	lines := strings.Split(generatedCode, "\n")
	if strings.HasPrefix(lines[0], " ") || strings.HasPrefix(lines[0], "\t") {
		t.Errorf("Expected generated code to not have leading whitespace")
	}
}

// TestIntegration_APIResponseComparison tests API response comparison
func TestIntegration_APIResponseComparison(t *testing.T) {
	// Real-world scenario: comparing API responses
	expectedResponse := text.StripMargin(`
		|{
		|  "status": "success",
		|  "data": {
		|    "users": [
		|      {"id": 1, "name": "Alice"},
		|      {"id": 2, "name": "Bob"}
		|    ]
		|  }
		|}`)

	actualResponse := text.StripMargin(`
		|{
		|  "status": "success",
		|  "data": {
		|    "users": [
		|      {"id": 1, "name": "Alice"},
		|      {"id": 2, "name": "Bob"}
		|    ]
		|  }
		|}`)

	diff, match := text.Diff(expectedResponse, actualResponse)

	if !match {
		t.Errorf("Expected identical API responses to match. Diff:\n%s", diff)
	}
}

// TestIntegration_MultilineStringComparison tests multiline string debugging
func TestIntegration_MultilineStringComparison(t *testing.T) {
	// Real-world scenario: debugging multiline string differences
	template1 := text.StripMargin(`
		|Dear {{.Name}},
		|
		|Thank you for your order #{{.OrderID}}.
		|Your items will be shipped within 3-5 business days.
		|
		|Best regards,
		|The Team`)

	template2 := text.StripMargin(`
		|Dear {{.Name}},
		|
		|Thank you for your order #{{.OrderID}}.
		|Your items will be shipped within 5-7 business days.
		|
		|Best regards,
		|The Team`)

	// Use CompareStrings for detailed comparison
	result := text.CompareStrings(template1, template2)

	// Should detect the difference in shipping time
	if strings.Contains(result, "✓ [MATCH]") {
		t.Errorf("Expected templates to differ")
	}

	// Should show the difference
	if !strings.Contains(result, "Difference at position") {
		t.Errorf("Expected difference position to be shown")
	}
}

// TestIntegration_TableDataFormatting tests table-like data formatting
func TestIntegration_TableDataFormatting(t *testing.T) {
	// Real-world scenario: formatting table data
	tableData := text.StripColumn(`
		|ID   | Name          | Email                |
		|-----|---------------|----------------------|
		|1    | Alice Smith   | alice@example.com    |
		|2    | Bob Jones     | bob@example.com      |
		|3    | Carol White   | carol@example.com    |`)

	// Verify table structure
	lines := strings.Split(tableData, "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 lines in table, got %d", len(lines))
	}

	// Each line should contain column data
	for i, line := range lines {
		if !strings.Contains(line, "|") {
			t.Errorf("Line %d should contain pipe separator, got: %s", i, line)
		}
	}
}

// TestIntegration_ComplexDiffScenario tests complex real-world diff
func TestIntegration_ComplexDiffScenario(t *testing.T) {
	// Real-world scenario: comparing configuration files
	oldConfig := text.StripMargin(`
		|version: 1.0
		|environment: production
		|services:
		|  - api
		|  - web
		|  - worker
		|features:
		|  new_ui: false
		|  analytics: true`)

	newConfig := text.StripMargin(`
		|version: 1.1
		|environment: production
		|services:
		|  - api
		|  - web
		|  - worker
		|features:
		|  new_ui: true
		|  analytics: true`)

	diff, match := text.Diff(oldConfig, newConfig)

	if match {
		t.Errorf("Expected configs to differ")
	}

	// Should show the version difference (first difference encountered)
	if !strings.Contains(diff, "≠") {
		t.Errorf("Expected diff to show difference indicator")
	}
}

// TestIntegration_WhitespaceDebugging tests whitespace debugging scenario
func TestIntegration_WhitespaceDebugging(t *testing.T) {
	// Real-world scenario: debugging whitespace issues
	withTabs := "column1\tcolumn2\tcolumn3"
	withSpaces := "column1 column2 column3"

	result := text.CompareStrings(withTabs, withSpaces)

	// Should visualize the difference
	if !strings.Contains(result, "␉") { // tab symbol
		t.Errorf("Expected tab symbol in output")
	}

	if !strings.Contains(result, "␣") { // space symbol
		t.Errorf("Expected space symbol in output")
	}

	// Should show exact character codes
	if !strings.Contains(result, "U+0009") { // tab unicode
		t.Errorf("Expected tab unicode in output")
	}

	if !strings.Contains(result, "U+0020") { // space unicode
		t.Errorf("Expected space unicode in output")
	}
}