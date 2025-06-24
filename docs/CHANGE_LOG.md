# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2025-06-23

### Added
- Cross-platform line ending normalization in text diff functionality
- Support for Windows (CRLF), Unix (LF), and Classic Mac (CR) line endings
- Enhanced Unicode and emoji support in string comparisons
- Comprehensive test coverage for edge cases including combining characters and RTL text
- Character-level difference detection with Unicode code point information
- Visual indicators for whitespace characters (spaces, tabs, carriage returns, etc.)
- Multi-line string comparison utilities with detailed line-by-line analysis
- Change log information docs/CHANGE_LOG.md

### Changed
- Improved diff output formatting with better alignment and padding
- Enhanced whitespace visualization symbols for better readability
- Refined error messaging and diff indicators for clearer test output

### Fixed
- Proper handling of empty strings and empty line comparisons
- Accurate positioning of difference markers in Unicode text
- Consistent behavior across different text encodings and character sets
- **Diff function trailing newline behavior to respect input string format** - The Diff function now only adds trailing newlines to output when input strings actually end with newlines, ensuring output format consistency with input format
- **Improved test debugging output** - Enhanced compareMultilineStrings function with detailed character-by-character analysis, length difference detection, and better error messaging for failed string comparisons

## [1.0.1] - 2025-06-23

### Added
- Initial release
- Core functionality
- Documentation

### Fixed
- Initial bug fixes

[Unreleased]: https://github.com/shapestone/textsmith-go/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/shapestone/textsmith-go/releases/tag/v1.1.0
[1.0.1]: https://github.com/shapestone/textsmith-go/releases/tag/v1.0.1
