# Release Guide

## Prerequisites

1. **Install GoReleaser**:
   ```bash
   go install github.com/goreleaser/goreleaser@latest
   ```

2. **Set up GitHub token** for Homebrew tap access:
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Create a token with `repo` scope
   - Add it as a secret in your repository: `HOMEBREW_TAP_TOKEN`

## Automated Release Process

### Option 1: Use the release script (Recommended)
```bash
./release.sh
```

This script will:
1. Build and test the project
2. Commit your changes
3. Create and push a git tag
4. Run GoReleaser to:
   - Build binaries for all platforms
   - Create a GitHub release
   - Update the Homebrew formula automatically
   - Push to your Homebrew tap repository

### Option 2: Manual GoReleaser
```bash
# Create and push a tag
git tag -a v1.1.4 -m "Release version 1.1.4"
git push origin v1.1.4

# Run GoReleaser
goreleaser release --clean
```

## What Gets Created

- **GitHub Release** with:
  - Source tarball
  - Binaries for Linux, macOS, and Windows (AMD64 and ARM64)
  - Checksums for verification

- **Homebrew Formula** automatically updated with:
  - New version number
  - Correct SHA256 checksum
  - Updated caveats

## Benefits

✅ **No more manual Homebrew formula updates**  
✅ **Automatic multi-platform builds**  
✅ **Consistent release process**  
✅ **Users get updates immediately**  

## Troubleshooting

If GoReleaser fails:
1. Check that `HOMEBREW_TAP_TOKEN` is set correctly
2. Verify the Homebrew tap repository exists and is accessible
3. Check the GitHub Actions logs for detailed error messages

## Manual Override

If you need to manually update the Homebrew formula:
```bash
cd /Users/aidan/projects/homebrew-tuiodo
# Make your changes
git add .
git commit -m "Manual update to v1.1.4"
git push origin master
```
