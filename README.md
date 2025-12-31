<div align="center">

![screenrecording-2025-12-31_10-50-29 (1)](https://github.com/user-attachments/assets/edbf16d6-26f4-4e75-b5ad-56b8361fe5d0)

# ClipTUI

[![Release](https://img.shields.io/github/v/release/ddVital/ClipTUI?style=for-the-badge)](https://github.com/ddVital/ClipTUI/releases/latest)
[![License](https://img.shields.io/github/license/ddVital/ClipTUI?style=for-the-badge)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ddVital/ClipTUI?style=for-the-badge)](go.mod)
[![AUR](https://img.shields.io/aur/version/cliptui?style=for-the-badge)](https://aur.archlinux.org/packages/cliptui)

**A modern, fast, and elegant clipboard history manager for Linux**

Built with Go • Beautiful TUI • Lightning Fast

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Configuration](#-configuration)

</div>

---

ClipTUI watches your system clipboard in the background, stores every item locally, and lets you browse, search, preview, and restore previous clipboard entries—directly from your terminal.

## Features

- **Live clipboard tracking** — Automatically captures anything you copy: text, code, links, commands
- **Beautiful TUI interface** — Clean aesthetics with smooth keyboard navigation
- **Powerful fuzzy search** — Instantly find old snippets, code blocks, or anything you've copied
- **Quick copy** — Number keys (0-9) for instant access to recent items
- **Item previews** — Full-screen preview mode for detailed viewing
- **Cross-desktop support** — Works on X11 and Wayland (GNOME, KDE, Sway, etc.)
- **Local-first storage** — Secure, offline history stored in SQLite
- **Lightning fast** — Pure Go binary with minimal dependencies
- **Systemd integration** — Run as a background service

## Installation

### Arch Linux (AUR)

```bash
# Build from source
yay -S cliptui

# Or use prebuilt binary
yay -S cliptui-bin
```

### Using Go

```bash
go install github.com/dvd/cliptui/cmd/cliptui@latest
```

### From Source

```bash
git clone https://github.com/ddVital/ClipTUI.git
cd ClipTUI
go build -o cliptui ./cmd/cliptui
sudo install -Dm755 cliptui /usr/local/bin/cliptui
```

### Download Binary

Download the latest binary from the [releases page](https://github.com/ddVital/ClipTUI/releases/latest).

```bash
# Download and install (replace VERSION with actual version)
curl -LO https://github.com/ddVital/ClipTUI/releases/download/vVERSION/cliptui_VERSION_linux_amd64.tar.gz
tar -xzf cliptui_VERSION_linux_amd64.tar.gz
sudo install -Dm755 cliptui /usr/local/bin/cliptui
```

## Usage

### Quick Start

```bash
# Start the background daemon
cliptui daemon

# Browse your clipboard history
cliptui
```

### Run as Systemd Service

```bash
# Enable and start the service
systemctl --user enable --now cliptui.service

# Check status
systemctl --user status cliptui.service

# View logs
journalctl --user -u cliptui.service -f
```

### Keyboard Shortcuts

<table>
<tr><th>List View</th><th>Action</th></tr>
<tr><td><kbd>0</kbd>-<kbd>9</kbd></td><td>Quick copy items 1-10</td></tr>
<tr><td><kbd>↑</kbd> / <kbd>k</kbd></td><td>Move up</td></tr>
<tr><td><kbd>↓</kbd> / <kbd>j</kbd></td><td>Move down</td></tr>
<tr><td><kbd>Enter</kbd> / <kbd>y</kbd></td><td>Copy selected item to clipboard</td></tr>
<tr><td><kbd>p</kbd></td><td>Preview item</td></tr>
<tr><td><kbd>/</kbd></td><td>Search mode</td></tr>
<tr><td><kbd>d</kbd></td><td>Delete selected item</td></tr>
<tr><td><kbd>D</kbd></td><td>Clear all history</td></tr>
<tr><td><kbd>q</kbd> / <kbd>Esc</kbd></td><td>Quit</td></tr>
</table>

<table>
<tr><th>Preview Mode</th><th>Action</th></tr>
<tr><td><kbd>Enter</kbd> / <kbd>y</kbd></td><td>Copy item to clipboard</td></tr>
<tr><td><kbd>Esc</kbd> / <kbd>q</kbd></td><td>Back to list</td></tr>
</table>

<table>
<tr><th>Search Mode</th><th>Action</th></tr>
<tr><td>Type to search</td><td>Fuzzy search through history</td></tr>
<tr><td><kbd>Enter</kbd></td><td>Confirm search</td></tr>
<tr><td><kbd>Esc</kbd></td><td>Cancel search</td></tr>
</table>

### CLI Commands

```bash
# Show clipboard history (interactive TUI)
cliptui

# Start background daemon
cliptui daemon

# Clear all history
cliptui clear

# Show help
cliptui --help
```

## Configuration

ClipTUI stores its data in `~/.local/share/cliptui/clipboard.db` by default.

### Custom Database Location

```bash
cliptui --db /custom/path/clipboard.db daemon
```

### Configuration Options

```bash
# Custom database location
cliptui --db ~/.config/cliptui/history.db daemon

# Limit maximum stored items
cliptui --max-items 500 daemon

# Show version
cliptui version
```

### Systemd Service Customization

Edit the service file to customize daemon behavior:

```bash
systemctl --user edit cliptui.service
```

Add your custom flags:
```ini
[Service]
ExecStart=
ExecStart=/usr/bin/cliptui --db /custom/path/clipboard.db daemon
```

## Technology Stack

- **[Go](https://go.dev/)** — High-performance compiled language
- **[tview](https://github.com/rivo/tview)** — Terminal UI framework
- **[tcell](https://github.com/gdamore/tcell)** — Terminal handling
- **[SQLite](https://www.sqlite.org/)** — Local database storage
- **[Chroma](https://github.com/alecthomas/chroma)** — Syntax highlighting
- **[Cobra](https://github.com/spf13/cobra)** — CLI framework

## Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details

## Support

- [Report a bug](https://github.com/ddVital/ClipTUI/issues/new)
- [Request a feature](https://github.com/ddVital/ClipTUI/issues/new)
- [Read the docs](https://github.com/ddVital/ClipTUI/wiki)
- [Discussions](https://github.com/ddVital/ClipTUI/discussions)

## ⭐ Star History

If you find ClipTUI useful, please consider giving it a star on GitHub!

---

<div align="center">

Made with ❤️ by the ddVital

[GitHub](https://github.com/ddVital/ClipTUI) • [Issues](https://github.com/ddVital/ClipTUI/issues) • [Releases](https://github.com/ddVital/ClipTUI/releases)

</div>
