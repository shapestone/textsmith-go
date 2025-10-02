# Automated Release Process Documentation

This document describes the automated release process used for library projects.

## Overview

The release process is fully automated using GitHub Actions. When a version tag is pushed, the workflow automatically:

1. Runs comprehensive tests
2. Executes benchmarks
3. Generates a categorized changelog from commit history
4. Creates a GitHub Release with formatted notes
5. Makes the release available via package managers (e.g., `go get`)

## Architecture

### Key Components

1. **Version Tags**: Semantic version tags (e.g., `v1.2.3`) trigger the release
2. **GitHub Actions Workflow**: `.github/workflows/release.yml` orchestrates the entire process
3. **Conventional Commits**: Commit message prefixes enable automatic changelog categorization
4. **Helper Script**: Optional `scripts/bump-version.sh` for interactive version bumping

### Workflow Jobs

The workflow consists of three sequential jobs:

```
test → release → notify
```

## GitHub Actions Workflow Configuration

### File: `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'

permissions:
  contents: write
```

**Trigger**: Any tag matching the pattern `v*.*.*` (semantic versioning)

**Permissions**: `contents: write` required to create releases

### Job 1: Pre-release Tests

```yaml
test:
  name: Pre-release Tests
  runs-on: ubuntu-latest
  steps:
    - name: Checkout code
      uses: actions/checkout@v5

    - name: Set up Go
      uses: actions/setup-go@v6
      with:
        go-version: '1.21'

    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Run benchmarks
      run: go test -bench=. -benchmem ./pkg/text
```

**Purpose**: Validates the code before creating a release

**Steps**:
- Checks out the repository
- Sets up the programming language environment (Go in this example)
- Runs full test suite with race detection and coverage
- Executes benchmarks to ensure performance

**Adaptation for Other Projects**:
- Replace Go setup with appropriate language setup action
- Modify test commands for your testing framework
- Adjust benchmark commands or remove if not applicable

### Job 2: Create Release

```yaml
release:
  name: Create Release
  needs: test
  runs-on: ubuntu-latest
```

**Dependencies**: Runs only after `test` job succeeds

#### Step 2.1: Extract Version

```yaml
- name: Get version from tag
  id: get_version
  run: echo "VERSION=${GITHUB_REF#refs/tags/}" >> $GITHUB_OUTPUT
```

**Purpose**: Extracts version number from the git tag (e.g., `v1.2.3`)

#### Step 2.2: Generate Changelog

```yaml
- name: Generate changelog
  id: changelog
  run: |
    # Get previous tag
    PREV_TAG=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")

    if [ -z "$PREV_TAG" ]; then
      # First release, get all commits
      COMMITS=$(git log --pretty=format:"- %s (%h)" --no-merges)
    else
      # Get commits since previous tag
      COMMITS=$(git log ${PREV_TAG}..HEAD --pretty=format:"- %s (%h)" --no-merges)
    fi

    # Categorize commits
    FEATURES=$(echo "$COMMITS" | grep -i "^- feat" || echo "")
    FIXES=$(echo "$COMMITS" | grep -i "^- fix" || echo "")
    DOCS=$(echo "$COMMITS" | grep -i "^- docs\|^- doc:" || echo "")
    TESTS=$(echo "$COMMITS" | grep -i "^- test" || echo "")
    REFACTOR=$(echo "$COMMITS" | grep -i "^- refactor" || echo "")
    CHORE=$(echo "$COMMITS" | grep -i "^- chore" || echo "")
    OTHER=$(echo "$COMMITS" | grep -vi "^- feat\|^- fix\|^- docs\|..." || echo "")

    # Build changelog with categories
    CHANGELOG="## What's Changed"$'\n\n'

    if [ ! -z "$FEATURES" ]; then
      CHANGELOG+="### ✨ Features"$'\n'"$FEATURES"$'\n\n'
    fi

    if [ ! -z "$FIXES" ]; then
      CHANGELOG+="### 🐛 Bug Fixes"$'\n'"$FIXES"$'\n\n'
    fi

    # ... additional categories ...

    echo "$CHANGELOG" > /tmp/changelog.md
    echo "CHANGELOG<<EOF" >> $GITHUB_OUTPUT
    echo "$CHANGELOG" >> $GITHUB_OUTPUT
    echo "EOF" >> $GITHUB_OUTPUT
```

**Purpose**: Automatically generates a categorized changelog from commit messages

**How it works**:
1. Finds the previous release tag using `git describe`
2. Extracts all commits between previous tag and current tag
3. Categorizes commits by prefix using grep patterns
4. Builds formatted changelog with sections for each category
5. Outputs changelog for use in release notes

**Commit Prefixes Used**:
- `feat:` → ✨ Features
- `fix:` → 🐛 Bug Fixes
- `docs:` → 📚 Documentation
- `test:` → 🧪 Tests
- `refactor:` → ♻️ Refactoring
- `chore:` → (typically not shown in user-facing changelog)
- Others → Other Changes

#### Step 2.3: Create GitHub Release

```yaml
- name: Create GitHub Release
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    gh release create "${{ steps.get_version.outputs.VERSION }}" \
      --repo="${GITHUB_REPOSITORY}" \
      --title="Release ${{ steps.get_version.outputs.VERSION }}" \
      --notes="$(cat <<'EOF'
    # ${{ steps.get_version.outputs.VERSION }}

    ${{ steps.changelog.outputs.CHANGELOG }}

    ## Installation

    \`\`\`bash
    go get github.com/shapestone/textsmith@${{ steps.get_version.outputs.VERSION }}
    \`\`\`

    ## Full Changelog

    See [CHANGELOG.md](https://github.com/shapestone/textsmith-go/blob/${{ steps.get_version.outputs.VERSION }}/CHANGELOG.md) for complete details.
    EOF
    )"
```

**Purpose**: Creates the GitHub Release with formatted notes

**Components**:
- Release title: "Release v1.2.3"
- Release notes include:
  - Version number as heading
  - Auto-generated categorized changelog
  - Installation instructions
  - Link to full CHANGELOG.md

### Job 3: Notifications

```yaml
notify:
  name: Post-release Notifications
  needs: release
  runs-on: ubuntu-latest
  if: success()
  steps:
    - name: Release successful
      run: |
        echo "✅ Release ${{ github.ref_name }} published successfully!"
        echo "📦 Users can now install via: go get github.com/shapestone/textsmith@${{ github.ref_name }}"
```

**Purpose**: Confirms successful release and logs installation command

**Extension Ideas**:
- Send Slack/Discord notifications
- Update documentation sites
- Trigger downstream builds
- Post to social media

## Conventional Commit Format

The changelog generation relies on conventional commit messages:

### Format

```
<type>: <description>

[optional body]

[optional footer]
```

### Standard Types

| Type | Description | Changelog Section | Version Impact |
|------|-------------|-------------------|----------------|
| `feat:` | New feature | ✨ Features | MINOR bump |
| `fix:` | Bug fix | 🐛 Bug Fixes | PATCH bump |
| `docs:` | Documentation only | 📚 Documentation | - |
| `test:` | Test changes | 🧪 Tests | - |
| `refactor:` | Code refactoring | ♻️ Refactoring | - |
| `perf:` | Performance improvements | ✨ Features or ♻️ Refactoring | MINOR/PATCH |
| `chore:` | Maintenance tasks | (usually omitted) | - |
| `ci:` | CI/CD changes | (usually omitted) | - |

### Breaking Changes

For MAJOR version bumps, add `BREAKING CHANGE:` in the commit footer:

```
feat: redesign API interface

BREAKING CHANGE: The Diff() function now returns a Result struct instead of string
```

### Examples

```
feat: add CompareStrings function for test assertions
fix: respect trailing newlines in Diff output
docs: update README with performance benchmarks
test: add fuzz tests for StripMargin
refactor: simplify internal diff algorithm
perf: optimize memory allocation in text processing
chore: update dependencies
```

## Setting Up for a New Project

### Step 1: Create the Workflow File

Create `.github/workflows/release.yml` with the content shown above.

**Customizations needed**:
1. Update programming language setup (lines 19-22)
2. Modify test commands (lines 24-28)
3. Update installation instructions in release notes (lines 119-121)
4. Adjust repository URLs

### Step 2: Configure Repository Permissions

Ensure GitHub Actions has permission to create releases:

1. Go to repository Settings → Actions → General
2. Under "Workflow permissions", select "Read and write permissions"
3. Enable "Allow GitHub Actions to create and approve pull requests" (if needed)

### Step 3: Adopt Conventional Commits

1. Document commit message conventions in CONTRIBUTING.md
2. Optional: Add commit message linting (commitlint)
3. Optional: Use git hooks to enforce format (husky + commitlint)

### Step 4: Create Helper Script (Optional)

Create `scripts/bump-version.sh` to automate version bumping:

```bash
#!/bin/bash
# Interactive version bump script
# - Shows current version
# - Prompts for bump type (major/minor/patch)
# - Updates CHANGELOG.md
# - Creates and pushes tag
```

### Step 5: Document the Process

Create `RELEASE.md` documenting:
- Pre-release checklist
- How to create releases (manual vs. script)
- Version bump guidelines
- Troubleshooting tips
- Emergency rollback procedures

## Release Workflow

### For Maintainers

1. **Pre-release validation**:
   ```bash
   make test      # or equivalent
   make fmt-fix   # or equivalent
   make vet       # or equivalent
   ```

2. **Create and push tag** (Manual method):
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```

3. **Or use helper script**:
   ```bash
   ./scripts/bump-version.sh
   ```

4. **Monitor GitHub Actions**:
   - Go to Actions tab
   - Watch the release workflow
   - Verify all jobs pass

5. **Verify the release**:
   - Check GitHub Releases page
   - Test installation command
   - Review auto-generated changelog

### Rollback Procedure

If a release has issues:

1. **Delete the problematic tag**:
   ```bash
   git tag -d v1.2.0
   git push origin :refs/tags/v1.2.0
   ```

2. **Delete the GitHub Release** (via UI or gh CLI):
   ```bash
   gh release delete v1.2.0
   ```

3. **Fix the issue**

4. **Recreate the tag**:
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```

## Advantages of This Approach

1. **Fully Automated**: No manual release note writing
2. **Consistent**: Same format for every release
3. **Fast**: Releases happen in minutes, not hours
4. **Quality Gates**: Tests must pass before release
5. **Traceable**: Clear history of what changed in each version
6. **Low Effort**: Just push a tag to release
7. **Conventional Commits**: Encourages better commit messages

## Customization Options

### Language-Specific Adaptations

**Python**:
- Replace Go setup with `actions/setup-python`
- Use `pytest` for testing
- Update installation instructions to `pip install`

**JavaScript/TypeScript**:
- Replace Go setup with `actions/setup-node`
- Use `npm test` or `yarn test`
- Update installation instructions to `npm install`

**Rust**:
- Replace Go setup with `actions-rs/toolchain`
- Use `cargo test` for testing
- Update installation instructions to `cargo add`

### Additional Features

**Add Asset Uploads**:
```yaml
- name: Build binaries
  run: make build-all

- name: Upload assets
  run: |
    gh release upload "$VERSION" \
      dist/binary-linux-amd64 \
      dist/binary-darwin-amd64 \
      dist/binary-windows-amd64.exe
```

**Add Changelog File Update**:
```yaml
- name: Update CHANGELOG.md
  run: |
    # Prepend generated changelog to CHANGELOG.md
    cat /tmp/changelog.md CHANGELOG.md > CHANGELOG.new.md
    mv CHANGELOG.new.md CHANGELOG.md
    git add CHANGELOG.md
    git commit -m "chore: update CHANGELOG for $VERSION"
    git push
```

**Add Security Scanning**:
```yaml
- name: Run security scan
  run: |
    # Go example
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
```

**Add Notifications**:
```yaml
- name: Notify Slack
  uses: slackapi/slack-github-action@v1
  with:
    webhook-url: ${{ secrets.SLACK_WEBHOOK }}
    payload: |
      {
        "text": "🎉 New release ${{ github.ref_name }} published!"
      }
```

## Best Practices

1. **Always use annotated tags**: `git tag -a v1.0.0 -m "message"`
2. **Follow semantic versioning strictly**
3. **Write descriptive commit messages** using conventional format
4. **Test locally before tagging**
5. **Use release candidates for major changes**: `v2.0.0-rc1`
6. **Keep CHANGELOG.md in sync** (if maintaining manually)
7. **Document breaking changes clearly** in commit messages
8. **Monitor the first few releases** of a new project carefully

## Troubleshooting

### Workflow doesn't trigger

**Check**:
- Tag format matches `v*.*.*` pattern
- Workflow file is on the default branch
- GitHub Actions are enabled for the repository

### Tests fail during release

**Solution**:
- Delete the tag
- Fix the failing tests
- Recreate the tag

### Changelog is empty or incorrect

**Check**:
- Commits use conventional format
- Grep patterns in workflow match your commit prefixes
- Previous tag exists and is reachable

### Permission denied when creating release

**Check**:
- Workflow has `permissions: contents: write`
- Repository settings allow Actions to create releases
- GITHUB_TOKEN is available (automatic in Actions)

## References

- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GitHub CLI - gh release](https://cli.github.com/manual/gh_release)

## Conclusion

This automated release process provides a robust, consistent, and efficient way to publish library releases. By combining semantic versioning, conventional commits, and GitHub Actions, the entire release workflow becomes a simple tag push operation while maintaining high quality standards through automated testing and changelog generation.
