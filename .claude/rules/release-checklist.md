# Release Checklist

## Before Creating a Release

ALWAYS complete these steps before pushing a tag or creating a release:

### 1. Update Documentation

- [ ] Update `README.md` with new features, commands, and keyboard shortcuts
- [ ] Update `FEATURES.md` if it exists
- [ ] Update `CLAUDE.md` if architecture changed

### 2. Create Manual Test File

Create/update `MANUAL_TEST.md` (gitignored) with test cases for new features:

```markdown
# Manual Test Checklist - vX.X.X

## New Features
- [ ] Feature 1: Description and how to test
- [ ] Feature 2: Description and how to test

## Regression Tests
- [ ] Timer starts/pauses/resets correctly
- [ ] Modes switch correctly
- [ ] Settings save and persist
- [ ] Notifications work
```

### 3. Version Bump

- Ensure version follows semver (vMAJOR.MINOR.PATCH)
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

### 4. Final Checks

- [ ] `go build` succeeds
- [ ] `go test ./...` passes
- [ ] Manual testing completed
- [ ] README reflects all new features
- [ ] Commit message follows conventional commits

### 5. Release Commands

```bash
# Commit changes
git add .
git commit -m "feat: description of changes"

# Create and push tag
git tag -a vX.X.X -m "vX.X.X - Release title"
git push origin main
git push origin vX.X.X

# Update README after release if needed
git add README.md
git commit -m "docs: update README for vX.X.X"
git push origin main
```

## After Release

- [ ] Verify GitHub Actions release workflow completed
- [ ] Check release page has all binaries
- [ ] Test install script works: `curl -sSL https://raw.githubusercontent.com/tauantcamargo/workpulse/main/install.sh | sh`
