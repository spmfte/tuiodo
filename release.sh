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

# Function to get confirmation
confirm() {
    read -p "$1 (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Operation cancelled by user"
        exit 1
    fi
}

# Check if we're in the right directory
if [ ! -f "main.go" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

# Build the project to ensure it compiles
print_step "Building project..."
go build -ldflags="-X main.Version=$CURRENT_VERSION -X main.BuildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ') -X main.GitCommit=$(git rev-parse --short HEAD)"
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

# Show changes to be committed
print_step "Files changed since last commit:"
git status --porcelain

# Get commit message
echo
read -p "Enter commit message: " COMMIT_MESSAGE

# Confirm actions
echo -e "\nThe following actions will be performed:"
echo -e "1. Stage and commit changes with message: ${GREEN}$COMMIT_MESSAGE${NC}"
echo -e "2. Create and push tag ${GREEN}v$NEW_VERSION${NC}"
echo -e "3. Push changes to remote repository"
echo
confirm "Do you want to proceed?"

# Stage all changes
print_step "Staging changes..."
git add .
print_success "Changes staged"

# Commit changes
print_step "Committing changes..."
git commit -m "$COMMIT_MESSAGE"
print_success "Changes committed"

# Create and push tag
print_step "Creating tag v$NEW_VERSION..."
git tag -a "v$NEW_VERSION" -m "Release version $NEW_VERSION"
print_success "Tag created"

# Push changes and tags
print_step "Pushing changes and tags to remote..."
git push origin master
git push origin "v$NEW_VERSION"
print_success "Changes and tags pushed"

# Final success message
echo -e "\n${GREEN}Release v$NEW_VERSION completed successfully!${NC}"
echo -e "Don't forget to:"
echo "1. Update the Homebrew formula with the new version"
echo "2. Update CHANGELOG.md with the release notes"
echo "3. Create a GitHub release with the changelog"

# Make the script executable
chmod +x release.sh 