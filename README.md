# 📋 clipTUI — A beautiful terminal-based clipboard history manager

**clipTUI** is a modern, fast, and elegant clipboard history manager for Linux, built with Go and designed with the same visual polish found in Charm's terminal apps.

It watches your system clipboard in the background, stores every item locally, and lets you browse, search, preview, and restore previous clipboard entries—directly from your terminal.

![clipTUI Demo](https://via.placeholder.com/800x400.png?text=clipTUI+Demo)

## ✨ Features

- **Live clipboard tracking** — Automatically captures anything you copy: text, code, links, commands
- **Beautiful TUI interface** — Built with Bubble Tea, featuring smooth transitions and clean aesthetics
- **Powerful fuzzy search** — Instantly find old snippets, code blocks, or anything you've copied
- **Quick paste** — Select an item and instantly send it back to the clipboard
- **Item previews** with:
  - Syntax highlighting for code
  - Markdown preview
  - Truncated or full-screen view
- **Cross-desktop support** — Works on X11, Wayland, GNOME, KDE, Sway
- **Local-first storage** — Secure, offline history stored in SQLite
- **Lightning fast** — Pure Go binary with zero dependencies
- **Easy packaging** — Distribute through pacman/AUR, .deb, .rpm, or static binaries

## 🚀 Installation

### From Source

```bash
git clone https://github.com/dvd/cliptui
cd cliptui
go build -o cliptui ./cmd/cliptui
sudo mv cliptui /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/dvd/cliptui/cmd/cliptui@latest
```

### Arch Linux (AUR)

```bash
yay -S cliptui
```

### Debian/Ubuntu

```bash
wget https://github.com/dvd/cliptui/releases/latest/download/cliptui_linux_amd64.deb
sudo dpkg -i cliptui_linux_amd64.deb
```

### RPM-based (Fedora, RHEL, etc.)

```bash
wget https://github.com/dvd/cliptui/releases/latest/download/cliptui_linux_amd64.rpm
sudo rpm -i cliptui_linux_amd64.rpm
```

## 🎯 Usage

### Start the clipboard monitor daemon

```bash
# Run in foreground
cliptui daemon

# Or enable as systemd user service
systemctl --user enable --now cliptui.service
```

### Browse clipboard history

```bash
cliptui
# or
cliptui show
```

### Keyboard shortcuts

**List view:**
- `↑/k` — Move up
- `↓/j` — Move down
- `Enter/y` — Copy selected item to clipboard
- `p` — Preview item
- `/` — Search mode
- `d` — Delete selected item
- `D` — Clear all history
- `q` — Quit

**Preview mode:**
- `Enter/y` — Copy item to clipboard
- `Esc/q` — Back to list

**Search mode:**
- Type to search
- `Enter` — Confirm search
- `Esc` — Cancel search

### Clear history

```bash
cliptui clear
```

## 🛠️ Technology Stack

- **Go** — High-performance static binary
- **Bubble Tea** — Terminal UI framework
- **Lipgloss** — Styles & layouts
- **SQLite** — History storage
- **atotto/clipboard** — Clipboard reading
- **sahilm/fuzzy** — Fuzzy search
- **alecthomas/chroma** — Syntax highlighting
- **cobra** — CLI framework
- **goreleaser + nfpm** — Packaging automation

## 📁 Project Structure

```
cliptui/
├── cmd/cliptui/          # Main application entry point
├── internal/
│   ├── clipboard/        # Clipboard monitoring
│   ├── config/           # Configuration management
│   ├── search/           # Fuzzy search implementation
│   ├── storage/          # SQLite database layer
│   └── tui/              # Bubble Tea UI components
├── pkg/types/            # Shared types
├── systemd/              # Systemd service files
└── scripts/              # Installation scripts
```

## 🎨 Configuration

clipTUI stores its data in `~/.local/share/cliptui/clipboard.db` by default.

You can customize the behavior with flags:

```bash
cliptui --db /custom/path/clipboard.db --max-items 500 daemon
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) for details

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) — For the amazing terminal UI libraries
- The Go community — For excellent tooling and libraries

## 🐛 Known Issues

- Clipboard monitoring requires X11 or Wayland with `wl-clipboard` installed
- May require `xsel` or `xclip` on some systems

## 📮 Support

If you encounter any issues or have questions:
- Open an issue on [GitHub](https://github.com/dvd/cliptui/issues)
- Check existing issues for solutions

---

Made with ❤️ by the community
