# Variables
GO_FILES=$(shell find . -name "*.go" -type f -not -name "*_test.go")
TEST_FILES=$(shell find . -name "*_test.go" -type f)

# Default target
all: check

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in short mode
test-short:
	@echo "Running tests in short mode..."
	go test -short -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Clean test artifacts
clean:
	rm -f coverage.out coverage.html

# Check code formatting
fmt:
	@echo "Checking code formatting..."
	gofmt -s -l .

# Format code
fmt-fix:
	@echo "Formatting code..."
	gofmt -s -w .

# Run static analysis
vet:
	@echo "Running go vet..."
	go vet ./...

# Run module tidy
mod-tidy:
	@echo "Running go mod tidy..."
	go mod tidy

# Run all quality checks
check: fmt vet test

# Show help
help:
	@echo "Available targets:"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-short     - Run tests in short mode"
	@echo "  test-race      - Run tests with race detection"
	@echo "  bench          - Run benchmarks"
	@echo "  clean          - Remove test artifacts"
	@echo "  fmt            - Check code formatting"
	@echo "  fmt-fix        - Format code"
	@echo "  vet            - Run static analysis"
	@echo "  mod-tidy       - Run go mod tidy"
	@echo "  check          - Run all quality checks (fmt + vet + test)"
	@echo "  help           - Show this help message"

.PHONY: all test test-coverage test-short test-race bench clean fmt fmt-fix vet mod-tidy check help
