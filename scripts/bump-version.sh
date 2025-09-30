#!/bin/bash

# Version bumping script for textsmith
# Usage: ./scripts/bump-version.sh [major|minor|patch]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get current version from latest tag
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${BLUE}Current version: ${CURRENT_VERSION}${NC}"

# Remove 'v' prefix for version parsing
VERSION_NUMBER=${CURRENT_VERSION#v}

# Parse version components
IFS='.' read -r -a VERSION_PARTS <<< "$VERSION_NUMBER"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]}"
PATCH="${VERSION_PARTS[2]}"

echo -e "${BLUE}Parsed: MAJOR=$MAJOR, MINOR=$MINOR, PATCH=$PATCH${NC}\n"

# Determine version bump type
if [ -z "$1" ]; then
    echo -e "${YELLOW}What type of release is this?${NC}"
    echo "1) MAJOR - Breaking changes (incompatible API changes)"
    echo "2) MINOR - New features (backwards compatible)"
    echo "3) PATCH - Bug fixes (backwards compatible)"
    echo ""
    read -p "Enter choice (1-3): " CHOICE

    case $CHOICE in
        1) BUMP_TYPE="major" ;;
        2) BUMP_TYPE="minor" ;;
        3) BUMP_TYPE="patch" ;;
        *) echo -e "${RED}Invalid choice${NC}"; exit 1 ;;
    esac
else
    BUMP_TYPE="$1"
fi

# Calculate new version
case $BUMP_TYPE in
    major)
        NEW_MAJOR=$((MAJOR + 1))
        NEW_MINOR=0
        NEW_PATCH=0
        ;;
    minor)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$((MINOR + 1))
        NEW_PATCH=0
        ;;
    patch)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$MINOR
        NEW_PATCH=$((PATCH + 1))
        ;;
    *)
        echo -e "${RED}Invalid bump type: $BUMP_TYPE${NC}"
        echo "Usage: $0 [major|minor|patch]"
        exit 1
        ;;
esac

NEW_VERSION="v${NEW_MAJOR}.${NEW_MINOR}.${NEW_PATCH}"

echo -e "\n${GREEN}New version will be: ${NEW_VERSION}${NC}\n"

# Check for uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    echo -e "${YELLOW}Warning: You have uncommitted changes${NC}"
    git status --short
    echo ""
    read -p "Continue anyway? (y/N): " CONTINUE
    if [ "$CONTINUE" != "y" ] && [ "$CONTINUE" != "Y" ]; then
        echo -e "${RED}Aborted${NC}"
        exit 1
    fi
fi

# Get unreleased changes from git log
echo -e "${BLUE}Fetching changes since ${CURRENT_VERSION}...${NC}"
CHANGES=$(git log ${CURRENT_VERSION}..HEAD --pretty=format:"- %s (%h)" --no-merges)

if [ -z "$CHANGES" ]; then
    echo -e "${RED}No changes found since ${CURRENT_VERSION}${NC}"
    exit 1
fi

echo -e "\n${BLUE}Changes to be released:${NC}"
echo "$CHANGES"
echo ""

# Update CHANGELOG.md
read -p "Update CHANGELOG.md? (Y/n): " UPDATE_CHANGELOG
if [ "$UPDATE_CHANGELOG" != "n" ] && [ "$UPDATE_CHANGELOG" != "N" ]; then
    echo -e "\n${BLUE}Updating CHANGELOG.md...${NC}"

    # Get today's date
    TODAY=$(date +%Y-%m-%d)

    # Create temporary file with new changelog entry
    {
        # Read until we find [Unreleased]
        sed -n '1,/## \[Unreleased\]/p' CHANGELOG.md

        # Add new version section
        echo ""
        echo "## [$NEW_VERSION] - $TODAY"
        echo ""

        # Categorize changes
        echo "$CHANGES" | grep -i "^- feat" > /tmp/features.txt 2>/dev/null || true
        echo "$CHANGES" | grep -i "^- fix" > /tmp/fixes.txt 2>/dev/null || true
        echo "$CHANGES" | grep -vi "^- feat\|^- fix" > /tmp/other.txt 2>/dev/null || true

        if [ -s /tmp/features.txt ]; then
            echo "### Added"
            cat /tmp/features.txt | sed 's/^- feat: /- /' | sed 's/^- feat:/- /'
            echo ""
        fi

        if [ -s /tmp/fixes.txt ]; then
            echo "### Fixed"
            cat /tmp/fixes.txt | sed 's/^- fix: /- /' | sed 's/^- fix:/- /'
            echo ""
        fi

        if [ -s /tmp/other.txt ]; then
            echo "### Changed"
            cat /tmp/other.txt
            echo ""
        fi

        # Add rest of changelog
        sed -n '/## \[Unreleased\]/,$p' CHANGELOG.md | tail -n +2
    } > CHANGELOG.md.tmp

    # Update version links at the bottom
    sed -i.bak "s|\[Unreleased\]:.*|[Unreleased]: https://github.com/shapestone/textsmith-go/compare/${NEW_VERSION}...HEAD\n[${NEW_VERSION}]: https://github.com/shapestone/textsmith-go/compare/${CURRENT_VERSION}...${NEW_VERSION}|" CHANGELOG.md.tmp

    mv CHANGELOG.md.tmp CHANGELOG.md
    rm -f CHANGELOG.md.bak /tmp/features.txt /tmp/fixes.txt /tmp/other.txt

    echo -e "${GREEN}✓ CHANGELOG.md updated${NC}"

    # Stage the changelog
    git add CHANGELOG.md
fi

# Commit changelog if changes were made
if git diff --cached --quiet; then
    echo -e "${YELLOW}No changes to commit${NC}"
else
    echo -e "\n${BLUE}Committing CHANGELOG.md...${NC}"
    git commit -m "chore: prepare release ${NEW_VERSION}"
    echo -e "${GREEN}✓ Changes committed${NC}"
fi

# Create and push tag
echo -e "\n${BLUE}Creating tag ${NEW_VERSION}...${NC}"
git tag -a "$NEW_VERSION" -m "Release ${NEW_VERSION}"
echo -e "${GREEN}✓ Tag created${NC}"

echo -e "\n${YELLOW}Ready to push tag to trigger release workflow${NC}"
read -p "Push tag to origin? (Y/n): " PUSH_TAG

if [ "$PUSH_TAG" != "n" ] && [ "$PUSH_TAG" != "N" ]; then
    echo -e "\n${BLUE}Pushing changes and tag...${NC}"
    git push origin main
    git push origin "$NEW_VERSION"

    echo -e "\n${GREEN}✓ Release ${NEW_VERSION} initiated!${NC}"
    echo -e "${BLUE}Check GitHub Actions for release progress:${NC}"
    echo -e "https://github.com/shapestone/textsmith-go/actions"
else
    echo -e "\n${YELLOW}Tag created but not pushed. To push manually:${NC}"
    echo -e "  git push origin main"
    echo -e "  git push origin ${NEW_VERSION}"
fi

echo -e "\n${GREEN}Done!${NC}"