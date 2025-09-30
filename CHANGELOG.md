# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-09-29

### Added
- Comprehensive test suite with 197+ tests achieving 97.0% coverage
- `CompareStrings` function for test framework style string comparison with visualization
- `CompareStringsRaw` function for string comparison without character visualization
- Performance tests validating behavior with files up to 100MB
- Edge case tests for invalid UTF-8, null bytes, mixed line endings, zero-width characters
- Godoc examples for all public functions
- Fuzz tests for robustness testing (Go 1.18+)
- Integration tests for real-world usage scenarios
- Benchmark suite for all public functions
- GitHub Actions CI/CD with multi-platform testing (Ubuntu, macOS, Windows)
- GitHub Actions workflow for automated releases
- Comprehensive performance documentation in README
- `CONTRIBUTING.md` with development guidelines
- GitHub issue templates (bug report, feature request)
- `RELEASE.md` with release process documentation
- Version bumping helper script

### Changed
- Updated README with comprehensive performance characteristics and benchmarks
- Improved CONTRIBUTING.md with detailed guidelines

### Fixed
- Fixed empty string formatting in CompareStrings (added ¶ marker)

## [1.0.1] - 2024-06-22

### Fixed
- Respect input string format in Diff function for trailing newlines
- Fixed handling of trailing newline differences in text comparison

### Changed
- Modified text_diff.go functionality to only perform pure diffs
- Improved test coverage and reliability

## [1.0.0] - 2024-06-22

### Added
- Initial release of textsmith
- `StripMargin` function for clean multiline string handling
- `StripColumn` function for column-based multiline string handling
- `Diff` function for visual side-by-side text comparison
- Whitespace visualization (tabs, spaces, line endings)
- Cross-platform line ending support (Unix, Windows, Mac)
- Full Unicode support including emojis and complex scripts
- Comprehensive test suite
- Performance optimizations
- Documentation with examples

[Unreleased]: https://github.com/shapestone/textsmith-go/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/shapestone/textsmith-go/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/shapestone/textsmith-go/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/shapestone/textsmith-go/releases/tag/v1.0.0