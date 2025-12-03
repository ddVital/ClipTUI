package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dvd/cliptui/internal/clipboard"
	"github.com/dvd/cliptui/internal/search"
	"github.com/dvd/cliptui/internal/storage"
	"github.com/dvd/cliptui/pkg/types"
)

type mode int

const (
	modeList mode = iota
	modePreview
	modeSearch
)

// Model represents the TUI state
type Model struct {
	storage       *storage.Storage
	items         []types.ClipboardItem
	filteredItems []types.ClipboardItem
	cursor        int
	mode          mode
	searchInput   textinput.Model
	width         int
	height        int
	err           error
}

// New creates a new TUI model
func New(store *storage.Storage) (*Model, error) {
	items, err := store.GetRecent(100)
	if err != nil {
		return nil, err
	}

	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100
	ti.Width = 50

	return &Model{
		storage:       store,
		items:         items,
		filteredItems: items,
		cursor:        0,
		mode:          modeList,
		searchInput:   ti,
		width:         80,
		height:        24,
	}, nil
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	// Update search input if in search mode
	if m.mode == modeSearch {
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)

		// Re-filter on search input change
		query := m.searchInput.Value()
		m.filteredItems = search.Filter(m.items, query)
		m.cursor = 0

		return m, cmd
	}

	return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modeSearch:
		switch msg.String() {
		case "esc":
			m.mode = modeList
			m.searchInput.SetValue("")
			m.filteredItems = m.items
			m.cursor = 0
			return m, nil
		case "enter":
			m.mode = modeList
			return m, nil
		}

	case modePreview:
		switch msg.String() {
		case "esc", "q":
			m.mode = modeList
			return m, nil
		case "enter", "y":
			// Copy selected item to clipboard
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				item := m.filteredItems[m.cursor]
				clipboard.SetClipboard(item.Content)
				return m, tea.Quit
			}
		}

	case modeList:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.filteredItems)-1 {
				m.cursor++
			}

		case "enter", "y":
			// Copy selected item to clipboard
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				item := m.filteredItems[m.cursor]
				clipboard.SetClipboard(item.Content)
				return m, tea.Quit
			}

		case "p":
			// Enter preview mode
			m.mode = modePreview
			return m, nil

		case "/":
			// Enter search mode
			m.mode = modeSearch
			m.searchInput.Focus()
			return m, textinput.Blink

		case "d":
			// Delete selected item
			if len(m.filteredItems) > 0 && m.cursor < len(m.filteredItems) {
				item := m.filteredItems[m.cursor]
				m.storage.Delete(item.ID)

				// Reload items
				items, _ := m.storage.GetRecent(100)
				m.items = items
				m.filteredItems = search.Filter(m.items, m.searchInput.Value())

				if m.cursor >= len(m.filteredItems) && m.cursor > 0 {
					m.cursor--
				}
			}

		case "D":
			// Clear all
			m.storage.Clear()
			m.items = []types.ClipboardItem{}
			m.filteredItems = []types.ClipboardItem{}
			m.cursor = 0
		}
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	switch m.mode {
	case modePreview:
		return m.renderPreview()
	case modeSearch:
		return m.renderSearch()
	default:
		return m.renderList()
	}
}

func (m *Model) renderList() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("📋 clipTUI - Clipboard History")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Items
	viewportHeight := m.height - 8 // Reserve space for header and footer
	start := m.cursor - viewportHeight/2
	if start < 0 {
		start = 0
	}
	end := start + viewportHeight
	if end > len(m.filteredItems) {
		end = len(m.filteredItems)
	}

	for i := start; i < end; i++ {
		item := m.filteredItems[i]
		cursor := " "
		style := itemStyle

		if i == m.cursor {
			cursor = "▸"
			style = selectedItemStyle
		}

		// Format timestamp
		timestamp := formatTimestamp(item.Timestamp)
		typeBadge := getTypeBadge(item.Type)

		line := fmt.Sprintf("%s %s %s  %s",
			cursor,
			typeBadge,
			truncate(item.Preview, 60),
			timestampStyle.Render(timestamp),
		)

		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	// Status bar
	b.WriteString("\n")
	statusText := fmt.Sprintf("Items: %d | Selected: %d/%d",
		len(m.items),
		m.cursor+1,
		len(m.filteredItems),
	)
	b.WriteString(statusStyle.Render(statusText))
	b.WriteString("\n")

	// Help
	help := helpStyle.Render("↑/k up • ↓/j down • enter/y copy • p preview • / search • d delete • D clear all • q quit")
	b.WriteString(help)

	return b.String()
}

func (m *Model) renderPreview() string {
	var b strings.Builder

	if len(m.filteredItems) == 0 || m.cursor >= len(m.filteredItems) {
		return "No item selected"
	}

	item := m.filteredItems[m.cursor]

	// Title
	title := titleStyle.Render(fmt.Sprintf("📋 Preview - %s", item.Type))
	b.WriteString(title)
	b.WriteString("\n\n")

	// Metadata
	timestamp := formatTimestamp(item.Timestamp)
	meta := fmt.Sprintf("Type: %s | Time: %s | Length: %d bytes",
		getTypeBadge(item.Type),
		timestamp,
		len(item.Content),
	)
	b.WriteString(statusStyle.Render(meta))
	b.WriteString("\n\n")

	// Content with syntax highlighting
	maxLines := m.height - 10
	preview := FormatPreview(item.Content, item.Type, maxLines)
	b.WriteString(previewStyle.Render(preview))
	b.WriteString("\n\n")

	// Help
	help := helpStyle.Render("enter/y copy • esc back • q quit")
	b.WriteString(help)

	return b.String()
}

func (m *Model) renderSearch() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("🔍 Search Clipboard History")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Search input
	b.WriteString(searchStyle.Render("Search: "))
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	// Show filtered results count
	resultInfo := fmt.Sprintf("Found %d items", len(m.filteredItems))
	b.WriteString(statusStyle.Render(resultInfo))
	b.WriteString("\n\n")

	// Show preview of results
	maxShow := min(5, len(m.filteredItems))
	for i := 0; i < maxShow; i++ {
		item := m.filteredItems[i]
		line := fmt.Sprintf("  %s %s",
			getTypeBadge(item.Type),
			truncate(item.Preview, 60),
		)
		b.WriteString(itemStyle.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Help
	help := helpStyle.Render("Type to search • enter confirm • esc cancel")
	b.WriteString(help)

	return b.String()
}

func formatTimestamp(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}

func truncate(s string, maxLen int) string {
	// Replace newlines with spaces for preview
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
