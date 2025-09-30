# Contributing to textsmith

Thank you for your interest in contributing to textsmith! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected vs actual behavior
- Your environment (Go version, OS)
- Code samples or test cases if applicable

### Suggesting Enhancements

Enhancement suggestions are welcome! Please:
- Use a clear, descriptive title
- Provide a detailed description of the proposed functionality
- Explain why this enhancement would be useful
- Include code examples if applicable

### Pull Requests

1. **Fork the repository** and create your branch from \`main\`
2. **Write tests** for your changes
3. **Ensure tests pass** by running \`make test\`
4. **Follow the coding style** - run \`make fmt-fix\` to format code
5. **Update documentation** if you're changing functionality
6. **Write clear commit messages**

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)

### Getting Started

\`\`\`shell
# Clone the repository
git clone https://github.com/shapestone/textsmith.git
cd textsmith

# Run tests
make test

# Run all quality checks
make check
\`\`\`

## Coding Guidelines

- Follow standard Go conventions
- Use descriptive test names
- Maintain or improve code coverage (currently at 97%)
- Add godoc comments for all exported functions

## Testing Requirements

- All new code must have tests
- Tests must pass on all supported platforms
- Maintain or improve code coverage

## Commit Message Conventions

We use conventional commit messages to help with changelog generation and semantic versioning:

- `feat:` - New feature (triggers MINOR version bump)
- `fix:` - Bug fix (triggers PATCH version bump)
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks
- `perf:` - Performance improvements

**Examples:**
```
feat: add new string comparison function
fix: correct trailing newline handling in Diff
docs: update README with examples
test: add edge case tests for Unicode
```

## Release Process (Maintainers Only)

Releases are automated via GitHub Actions. See [RELEASE.md](RELEASE.md) for detailed information.

### Quick Release Guide

1. Ensure all tests pass and CI is green
2. Update `CHANGELOG.md` with changes
3. Create and push a version tag:
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```
4. GitHub Actions automatically creates the release

For complete details, see [RELEASE.md](RELEASE.md).

Thank you for contributing to textsmith!
