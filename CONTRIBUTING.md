# Contributing to textsmith

Thank you for your interest in contributing to textsmith! This document provides guidelines and information for contributors.

## Table of Contents
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)
- [Project Structure](#project-structure)

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- Make (for using the Makefile commands)

### Fork and Clone
1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/shapestone/textsmith.git
   cd textsmith
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/shapestone/textsmith.git
   ```

## Development Setup

### Install Dependencies
```bash
# Download and verify dependencies
go mod download
go mod verify

# Tidy up module dependencies
make mod-tidy
```

### Verify Setup
```bash
# Run all quality checks
make check

# This runs: format + vet + test
```

## Making Changes

### Branch Naming
Create descriptive branch names:
- `feature/add-new-function` - for new features
- `fix/handle-empty-input` - for bug fixes
- `docs/update-readme` - for documentation changes
- `perf/optimize-regex` - for performance improvements

### Development Workflow
1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes
3. Run tests frequently:
   ```bash
   make test
   ```

4. Run all quality checks before committing:
   ```bash
   make check
   ```

## Testing

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests in short mode (faster)
make test-short

# Run tests with race detection
make test-race

# Run benchmarks
make bench
```

### Writing Tests
- Prefer given, when, then style tests
- Use table-driven tests for comprehensive value coverage
- Test edge cases (empty strings, unicode, large inputs)
- Include benchmarks for performance validation
- Follow the existing test patterns in `*_test.go` files

### Test Requirements
- All new functions must have comprehensive unit tests
- Maintain 100% test coverage
- Include both positive and negative test cases
- Test cross-platform behavior (line endings)

## Code Style

### Formatting
```bash
# Auto-format code
make fmt-fix

# Check formatting (without fixing)
go fmt ./...
```

### Static Analysis
```bash
# Run Go vet
make vet

# Or directly:
go vet ./...
```

### Guidelines
- Follow standard Go conventions
- Use clear, descriptive function and variable names
- Add comments for exported functions
- Keep functions focused and single-purpose
- Handle edge cases gracefully (no panics)

## Submitting Changes

### Before Submitting
1. Ensure all tests pass:
   ```bash
   make test
   ```

2. Run the full quality check:
   ```bash
   make check
   ```

3. Update documentation if needed
4. Add tests for new functionality

### Pull Request Process
1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a Pull Request on GitHub with:
    - Clear title describing the change
    - Detailed description of what was changed and why
    - Reference any related issues
    - Screenshots or examples if applicable

3. Ensure CI checks pass
4. Respond to any review feedback
5. Keep your branch up to date with main:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### Pull Request Requirements
- All tests must pass
- Code coverage must be maintained
- Code must be properly formatted
- Include appropriate documentation updates
- Follow the existing code patterns and style

## Release Process

### Overview
textsmith uses semantic versioning and Git tags for releases. The process is designed to be straightforward and automated.

### Creating a Release
**Note: Only maintainers can create releases.**

1. **Prepare for Release**
    - Ensure all changes are merged to `main`
    - All tests pass on `main` branch
    - Update any version references if needed
    - Review and update CHANGELOG if maintained

2. **Create Git Tag**
   ```bash
   # Ensure you're on the latest main
   git checkout main
   git pull origin main
   
   # Create and push tag (use semantic versioning)
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **Automatic Distribution**
    - Go's module system automatically detects the new tag
    - No additional steps needed for distribution
    - Users can install with: `go get github.com/shapestone/textsmith@v1.0.0`

### Version Numbering
Follow [Semantic Versioning](https://semver.org/):
- **MAJOR** (v2.0.0): Breaking changes to public API
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

### CI/CD Pipeline
- GitHub Actions workflow in `.github/workflows/go.yml`
- Automatically runs on every push and pull request
- Tests multiple Go versions and platforms
- No manual intervention required for releases

### Release Checklist
- [ ] All tests pass on main branch
- [ ] Code coverage is maintained
- [ ] Documentation is up to date
- [ ] Breaking changes are documented
- [ ] Version tag follows semantic versioning
- [ ] Tag is pushed to origin

## Project Structure

```
textsmith/
├── .github/
│   └── workflows/
│       └── go.yml              # GitHub Actions CI/CD
├── docs/
│   └── MODULE_SPEC.md          # Technical specification
├── pkg/
│   └── text/                   # Main package
│       ├── strip_margin.go     # StripMargin implementation
│       ├── strip_column.go     # StripColumn implementation
│       └── text_diff.go        # Diff implementation
├── strip_margin_test.go        # StripMargin tests
├── strip_column_test.go        # StripColumn tests
├── text_diff_test.go           # Diff tests
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── Makefile                    # Build and test commands
├── README.md                   # User documentation
├── CONTRIBUTING.md             # This file
└── LICENSE                     # License information
```

### Key Directories
- **`pkg/text/`** - Main library code
- **`docs/`** - Technical documentation
- **`.github/workflows/`** - CI/CD configuration
- **Root directory** - Tests, configuration, and documentation

## Getting Help

### Resources
- Check existing issues and pull requests
- Review the [README.md](README.md) for usage examples
- Read the [Technical Specification](docs/MODULE_SPEC.md) for implementation details

### Communication
- Open an issue for questions or discussion
- Use pull request comments for code-specific questions
- Be respectful and constructive in all interactions

### Reporting Issues
When reporting bugs or requesting features:
1. Search existing issues first
2. Use issue templates if available
3. Provide clear reproduction steps
4. Include Go version and platform information
5. Add relevant code examples or error messages

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help maintain a welcoming environment
- Follow GitHub's community guidelines

Thank you for contributing to textsmith!