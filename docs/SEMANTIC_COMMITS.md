# Git Commit Message Guide

## Overview

This document provides guidelines for writing consistent, meaningful commit messages using semantic commit conventions. Following these standards improves code history readability and enables automated tooling for versioning and changelog generation.

## Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Structure Rules

- **Header**: Required, max 50 characters
- **Body**: Optional, wrap at 72 characters
- **Footer**: Optional, used for breaking changes and issue references

## Commit Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature for the user | `feat(auth): add JWT token validation` |
| `fix` | Bug fix for the user | `fix(api): resolve null pointer in user service` |
| `docs` | Documentation changes | `docs: update API documentation` |
| `style` | Code style changes (formatting, semicolons, etc.) | `style: format code with gofmt` |
| `refactor` | Code refactoring without feature changes | `refactor(db): optimize user queries` |
| `perf` | Performance improvements | `perf(cache): implement Redis caching` |
| `test` | Adding or updating tests | `test(auth): add login unit tests` |
| `chore` | Maintenance tasks | `chore: update Go dependencies` |
| `ci` | CI/CD configuration changes | `ci: add GitHub Actions workflow` |
| `build` | Build system or dependency changes | `build: update Dockerfile` |
| `revert` | Revert previous commit | `revert: "feat(auth): add JWT validation"` |

## Scope Guidelines

The scope is optional and should specify the area of the codebase affected:

### Go Project Scopes (based on project layout)
- `cmd`: Changes to main applications
- `internal`: Internal package changes
- `pkg`: Public library changes
- `api`: API specification changes
- `docs`: Documentation changes
- `build`: Build and CI changes
- `scripts`: Script changes

### Examples by Component
```bash
feat(cmd/server): add health check endpoint
fix(internal/auth): handle expired tokens correctly
docs(api): update OpenAPI specification
chore(build): update Docker base image
test(pkg/utils): add validation helper tests
```

## Writing Guidelines

### Header (Required)
- Use imperative mood: "add" not "added" or "adds"
- No capital letter after type
- No period at the end
- Keep under 50 characters

### Body (Optional)
- Use imperative mood
- Explain **what** and **why**, not **how**
- Wrap at 72 characters
- Separate from header with blank line

### Footer (Optional)
- Reference issues: `Closes #123` or `Fixes #456`
- Breaking changes: `BREAKING CHANGE: <description>`

## Examples

### Simple Feature
```
feat(auth): add JWT token validation

Implement JWT token validation middleware to secure API endpoints.
Tokens are validated against secret key and expiration time.

Closes #42
```

### Bug Fix
```
fix(database): prevent connection pool exhaustion

Connection pool was not properly releasing connections after failed
queries, leading to pool exhaustion under high load.

- Add connection cleanup in error handlers
- Implement connection timeout configuration
- Add monitoring for pool statistics

Fixes #156
```

### Breaking Change
```
feat(api): update user endpoint response format

BREAKING CHANGE: User API now returns ISO 8601 timestamps instead
of Unix timestamps. Update client code to handle new format.

Before: {"created_at": 1640995200}
After: {"created_at": "2022-01-01T00:00:00Z"}

Closes #89
```

### Documentation Update
```
docs: add deployment guide for production

- Add step-by-step deployment instructions
- Include environment variable configuration
- Add troubleshooting section
```

### Refactoring
```
refactor(internal/user): extract validation logic

Move user validation logic to separate package for better reusability
and testability. No functional changes.
```

## Best Practices

### DO ✅
- Write clear, descriptive messages
- Use present tense, imperative mood
- Include context in the body when needed
- Reference related issues
- Keep commits atomic (one logical change)

### DON'T ❌
- Use vague messages like "fix stuff" or "update code"
- Include debugging information in commit message
- Make massive commits with multiple unrelated changes
- Use past tense ("fixed", "added")
- Exceed character limits

## Tools and Automation

### Commit Message Validation
Consider using tools like:
- **commitlint**: Lint commit messages
- **husky**: Git hooks for validation
- **conventional-changelog**: Generate changelogs automatically

### Git Hooks Example
Add to `.githooks/commit-msg`:
```bash
#!/bin/sh
# Validate commit message format
npx commitlint --edit $1
```

## Integration with Go Projects

### Module-Specific Commits
When working with Go modules, be specific about the affected package:

```bash
feat(internal/auth): implement OAuth2 provider
fix(pkg/validator): handle empty string validation
test(cmd/cli): add integration tests for user commands
docs(internal/database): add connection pool documentation
```

### Dependency Updates
```bash
chore(deps): update golang.org/x/crypto to v0.14.0
build(go.mod): upgrade to Go 1.21
```

### Configuration Changes
```bash
feat(config): add environment variable support
fix(config): resolve YAML parsing edge case
```

## Quick Reference

| Action | Type | Example |
|--------|------|---------|
| Add new feature | `feat` | `feat(api): add user search endpoint` |
| Fix bug | `fix` | `fix(auth): handle invalid token gracefully` |
| Update docs | `docs` | `docs: update README installation steps` |
| Refactor code | `refactor` | `refactor(utils): simplify error handling` |
| Add tests | `test` | `test(service): add unit tests for user CRUD` |
| Update dependencies | `chore` | `chore: update all Go dependencies` |
| Performance improvement | `perf` | `perf(db): add database query optimization` |
| Style/formatting | `style` | `style: run gofmt on all source files` |

## Resources

- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [Angular Commit Message Guidelines](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#commit)
- [Semantic Versioning](https://semver.org/)

---

*This guide should be followed by all team members to maintain consistent commit history and enable automated tooling.*