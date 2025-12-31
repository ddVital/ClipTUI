# Contributing to ClipTUI

Thank you for considering contributing to ClipTUI! We welcome contributions from the community.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Commit Guidelines](#commit-guidelines)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

Be respectful and considerate. We're all here to build something useful together.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Test your changes
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.24.0 or higher
- Git
- SQLite3 development libraries
- `xsel` or `wl-clipboard` (for clipboard access)

### Clone and Build

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/ClipTUI.git
cd ClipTUI

# Install dependencies
go mod download

# Build
go build -o cliptui ./cmd/cliptui

# Run tests
go test ./...

# Run the application
./cliptui daemon
```

### Project Structure

```
ClipTUI/
├── cmd/cliptui/          # Main application entry point
├── internal/
│   ├── clipboard/        # Clipboard monitoring
│   ├── config/           # Configuration management
│   ├── search/           # Fuzzy search implementation
│   ├── storage/          # SQLite database layer
│   └── tui/              # TUI components
├── pkg/types/            # Shared types
├── systemd/              # Systemd service files
└── scripts/              # Build and installation scripts
```

## Making Changes

### Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add comments for exported functions and types
- Write meaningful variable names

### Testing

- Write tests for new functionality
- Ensure existing tests pass: `go test ./...`
- Test manually with the TUI before submitting

### Documentation

- Update README.md if adding user-facing features
- Add comments for complex logic
- Update help text if adding new commands

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) for clear commit history.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```bash
feat(search): add regex support for advanced filtering

fix(clipboard): prevent duplicate entries from being stored

docs(readme): update installation instructions for Arch Linux

refactor(tui): split app.go into smaller modules
```

### Scope

Use the module name: `clipboard`, `tui`, `storage`, `search`, `config`

## Submitting a Pull Request

1. **Update your fork**
   ```bash
   git checkout dev
   git pull upstream dev
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes and commit**
   ```bash
   git add .
   git commit -m "feat(scope): description"
   ```

4. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Open a Pull Request**
   - Go to the original repository
   - Click "New Pull Request"
   - Select your branch
   - Fill in the PR template
   - Submit

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Tests pass (`go test ./...`)
- [ ] Code builds without errors
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventional commits
- [ ] Branch is up to date with `dev`

## Reporting Bugs

When reporting bugs, please include:

- ClipTUI version (`cliptui version`)
- Operating system and version
- Desktop environment (X11/Wayland, GNOME/KDE/etc.)
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages or logs

### Bug Report Template

```markdown
**Version:** v1.0.0
**OS:** Arch Linux (6.x.x)
**Desktop:** Wayland + KDE Plasma

**Steps to reproduce:**
1. Start daemon with `cliptui daemon`
2. Copy text from Firefox
3. Open TUI with `cliptui`

**Expected:** Text appears in history
**Actual:** History is empty

**Logs:**
```
journalctl --user -u cliptui.service
```
\`\`\`
```

## Suggesting Features

Feature requests are welcome! Please include:

- Clear description of the feature
- Use case / why it's needed
- Proposed implementation (if you have ideas)
- Mockups or examples (if applicable)

## Development Tips

### Live Reload During Development

```bash
# Install air for live reloading
go install github.com/cosmtrek/air@latest

# Run with live reload
air
```

### Debugging

```bash
# Run with verbose logging
go run ./cmd/cliptui --debug daemon

# Check database contents
sqlite3 ~/.local/share/cliptui/clipboard.db "SELECT * FROM clipboard_history;"
```

### Building for Release

```bash
# Build with version info
go build -ldflags="-s -w -X main.version=v1.0.0" -o cliptui ./cmd/cliptui
```

## Questions?

- Open a [Discussion](https://github.com/ddVital/ClipTUI/discussions)
- Ask in an [Issue](https://github.com/ddVital/ClipTUI/issues)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
