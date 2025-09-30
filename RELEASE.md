# Release Process

This document describes the release process for textsmith maintainers.

## Overview

textsmith follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner
- **PATCH** version for backwards compatible bug fixes

## Automated Release Process

Releases are automated via GitHub Actions. When you push a version tag, the workflow:

1. Runs the full test suite
2. Runs benchmarks
3. Generates a changelog from commits
4. Creates a GitHub Release
5. Makes the release available via `go get`

## Release Checklist

### Pre-Release

- [ ] Ensure all tests pass locally: `make test`
- [ ] Ensure code is formatted: `make fmt-fix`
- [ ] Run static analysis: `make vet`
- [ ] Review open issues and PRs
- [ ] Update documentation if needed
- [ ] Ensure CI is green on main branch

### Creating a Release

#### Option 1: Using the Helper Script (Recommended)

```bash
# Interactive version bump
./scripts/bump-version.sh

# The script will:
# 1. Show current version
# 2. Ask for bump type (major/minor/patch)
# 3. Update CHANGELOG.md with unreleased changes
# 4. Create and push the tag
```

#### Option 2: Manual Release

1. **Update CHANGELOG.md**

   Move unreleased changes to a new version section:

   ```markdown
   ## [1.2.0] - 2024-09-29

   ### Added
   - New feature X
   - New feature Y

   ### Fixed
   - Bug fix Z
   ```

   Update the links at the bottom:

   ```markdown
   [Unreleased]: https://github.com/shapestone/textsmith-go/compare/v1.2.0...HEAD
   [1.2.0]: https://github.com/shapestone/textsmith-go/compare/v1.1.0...v1.2.0
   ```

2. **Commit the changelog**

   ```bash
   git add CHANGELOG.md
   git commit -m "chore: prepare release v1.2.0"
   git push origin main
   ```

3. **Create and push the tag**

   ```bash
   # Create annotated tag
   git tag -a v1.2.0 -m "Release v1.2.0"

   # Push tag to trigger release workflow
   git push origin v1.2.0
   ```

4. **Verify the release**

   - Go to [GitHub Releases](https://github.com/shapestone/textsmith-go/releases)
   - Verify the release was created
   - Check that release notes are correct
   - Verify the workflow completed successfully

### Post-Release

1. **Verify installation**

   ```bash
   # Test that users can install the new version
   go get github.com/shapestone/textsmith@v1.2.0
   ```

2. **Announce the release** (if significant)
   - Update project README badges if needed
   - Post to relevant channels/forums
   - Update documentation sites

3. **Monitor for issues**
   - Watch for bug reports
   - Monitor GitHub issues
   - Be ready to release a patch if needed

## Determining Version Bump

### MAJOR version (x.0.0)

Increment when making incompatible API changes:

- Removing public functions or methods
- Changing function signatures
- Removing or renaming public types
- Breaking behavioral changes

**Example:** Removing `Diff()` function or changing its signature

### MINOR version (1.x.0)

Increment when adding functionality in a backwards compatible manner:

- Adding new public functions
- Adding new parameters with defaults
- Adding new types
- Performance improvements
- New features

**Example:** Adding `CompareStrings()` function

### PATCH version (1.0.x)

Increment when making backwards compatible bug fixes:

- Bug fixes
- Documentation updates
- Internal refactoring
- Performance improvements without API changes
- Test improvements

**Example:** Fixing a bug in `Diff()` output formatting

## Emergency Rollback

If a release has critical issues:

1. **Create a patch release immediately**

   ```bash
   # Create patch with fix
   git checkout -b hotfix/critical-bug
   # ... make fixes ...
   git commit -m "fix: critical bug in release"
   git push origin hotfix/critical-bug

   # Merge to main
   # Create new patch release tag
   git tag -a v1.2.1 -m "Release v1.2.1 - Critical hotfix"
   git push origin v1.2.1
   ```

2. **Mark the problematic release**
   - Edit the GitHub Release
   - Mark as "This release contains critical bugs"
   - Add link to the fixed version

3. **Communicate the issue**
   - Post issue/announcement
   - Notify users to upgrade

## Commit Message Conventions

Use conventional commits to help with changelog generation:

- `feat:` - New feature (MINOR bump)
- `fix:` - Bug fix (PATCH bump)
- `docs:` - Documentation changes
- `test:` - Test changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks
- `perf:` - Performance improvements

**Examples:**
```
feat: add CompareStrings function for test assertions
fix: respect trailing newlines in Diff output
docs: update README with performance benchmarks
test: add fuzz tests for StripMargin
```

## Release Schedule

There is no fixed release schedule. Releases are made when:

- Significant features are ready
- Critical bugs need fixing
- Enough changes have accumulated

## Testing Releases

To test the release process without publishing:

1. Use a feature branch
2. Create a tag with `-rc` suffix: `v1.2.0-rc1`
3. Push to verify workflow runs
4. Delete the tag if testing: `git tag -d v1.2.0-rc1`

## Troubleshooting

### Release workflow failed

1. Check GitHub Actions logs
2. Fix the issue
3. Delete the tag: `git tag -d v1.2.0 && git push origin :refs/tags/v1.2.0`
4. Fix the problem
5. Recreate and push the tag

### Tag was pushed prematurely

```bash
# Delete local tag
git tag -d v1.2.0

# Delete remote tag
git push origin :refs/tags/v1.2.0

# Fix issues and recreate tag
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0
```

### Need to update release notes

1. Go to GitHub Releases page
2. Click "Edit" on the release
3. Update the release notes
4. Save changes

## Questions?

If you have questions about the release process:

1. Review this document
2. Check the [GitHub Actions workflow](.github/workflows/release.yml)
3. Open an issue for discussion
4. Reach out to maintainers