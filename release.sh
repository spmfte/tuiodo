#!/bin/bash

# Exit on error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print step information
print_step() {
    echo -e "${YELLOW}[STEP]${NC} $1"
}

# Function to print success messages
print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Function to print error messages
print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

# Check if we have uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    print_error "You have uncommitted changes. Please commit or stash them first."
    git status --porcelain
    exit 1
fi

# Build the project to ensure it compiles
print_step "Building project..."
go build -o tuiodo
if [ $? -ne 0 ]; then
    print_error "Build failed"
    exit 1
fi
print_success "Build successful"

# Get the current version from the binary
CURRENT_VERSION=$(./tuiodo --version | awk '/TUIODO/ {print $2}')
if [ -z "$CURRENT_VERSION" ]; then
    print_error "Could not determine current version"
    exit 1
fi

# Remove the 'v' prefix if it exists
CURRENT_VERSION=${CURRENT_VERSION#v}

# Ask for the new version
echo -e "\nCurrent version is: ${GREEN}$CURRENT_VERSION${NC}"
read -p "Enter new version number (without 'v' prefix): " NEW_VERSION

# Validate version number format
if ! [[ $NEW_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version number format. Please use semantic versioning (e.g., 1.2.3)"
    exit 1
fi

# Check if tag already exists (locally or remotely)
if git tag -l "v$NEW_VERSION" | grep -q "v$NEW_VERSION" || git ls-remote --tags origin "v$NEW_VERSION" | grep -q "v$NEW_VERSION"; then
    print_error "Version v$NEW_VERSION already exists. Please use a different version number."
    exit 1
fi

# Get commit message
echo
read -p "Enter commit message for this release: " COMMIT_MESSAGE

# Confirm actions
echo -e "\nThe following actions will be performed:"
echo -e "1. Create and push tag v$NEW_VERSION"
echo -e "2. GitHub Actions will automatically:"
echo -e "   - Build binaries for all platforms"
echo -e "   - Create GitHub release with assets"
echo -e "   - Update Homebrew formula"
echo
read -p "Do you want to proceed? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_error "Operation cancelled by user"
    exit 1
fi

# Create and push tag
print_step "Creating tag v$NEW_VERSION..."
git tag -a "v$NEW_VERSION" -m "Release version $NEW_VERSION"
print_success "Tag created locally"

# Push tag to trigger GitHub Actions
print_step "Pushing tag to trigger automated release..."
git push origin "v$NEW_VERSION"
print_success "Tag pushed successfully"

# Final success message
echo -e "\n${GREEN}Release v$NEW_VERSION initiated successfully!${NC}"
echo -e "GitHub Actions is now building and releasing your software."
echo -e "You can monitor progress at: https://github.com/spmfte/tuiodo/actions"
echo -e "\nOnce complete, users can:"
echo -e "1. Download binaries from: https://github.com/spmfte/tuiodo/releases/tag/v$NEW_VERSION"
echo -e "2. Update Homebrew: brew upgrade tuiodo" 