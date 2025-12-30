package tui

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dvd/cliptui/internal/clipboard"
	"github.com/dvd/cliptui/internal/search"
	"github.com/dvd/cliptui/internal/storage"
	"github.com/dvd/cliptui/pkg/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// maxItemsToFetch is the maximum number of clipboard items to fetch from storage
	maxItemsToFetch = 100
	// clipboardPollInterval is how often to check for clipboard changes
	clipboardPollInterval = 500 * time.Millisecond
	// previewTruncateLength is the maximum length for list preview text
	previewTruncateLength = 80
	// previewFormatMaxLength is the maximum length for formatted preview content
	previewFormatMaxLength = 1000
	// searchInputHeight is the height of the search input widget
	searchInputHeight = 3
)

// AppState holds the application state
type AppState struct {
	mu            sync.RWMutex
	storage       *storage.Storage
	items         []types.ClipboardItem
	filteredItems []types.ClipboardItem
	cursor        int
	currentMode   mode
	searchQuery   string
}

// App represents the tview application
type App struct {
	app   *tview.Application
	pages *tview.Pages
	state *AppState

	// Widgets
	listWidget    *tview.Table
	listContainer *tview.Flex
	listHelp      *tview.TextView
	searchInput   *tview.InputField
	mainFlex      *tview.Flex

	previewView   *tview.TextView
	previewHeader *tview.TextView
	previewHelp   *tview.TextView
}

// New creates a new TUI application
func New(store *storage.Storage) (*App, error) {
	items, err := store.GetRecent(maxItemsToFetch)
	if err != nil {
		return nil, err
	}

	// Configure tview to use terminal default colors
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorDefault
	tview.Styles.PrimaryTextColor = tcell.ColorDefault
	tview.Styles.InverseTextColor = tcell.ColorDefault
	tview.Styles.ContrastSecondaryTextColor = tcell.ColorDefault

	app := &App{
		app: tview.NewApplication(),
		state: &AppState{
			storage:       store,
			items:         items,
			filteredItems: items,
			cursor:        0,
			currentMode:   modeList,
			searchQuery:   "",
		},
	}

	app.pages = tview.NewPages()

	listPage := app.buildListPage()
	previewPage := app.buildPreviewPage()

	app.pages.AddPage("list", listPage, true, true)
	app.pages.AddPage("preview", previewPage, true, false)

	app.app.SetRoot(app.pages, true)
	app.setupGlobalKeys()

	// Enable mouse capture (prevents terminal text selection, enables mouse events)
	app.app.EnableMouse(true)

	app.updateListDisplay()

	return app, nil
}

// Run starts the tview application
func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.monitorClipboard(ctx)

	err := a.app.Run()

	cancel()

	return err
}

// monitorClipboard checks for new clipboard items periodically
func (a *App) monitorClipboard(ctx context.Context) {
	ticker := time.NewTicker(clipboardPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			items, err := a.state.storage.GetRecent(maxItemsToFetch)
			if err != nil {
				continue
			}

			a.state.mu.Lock()
			needsUpdate := len(items) != len(a.state.items) ||
				(len(items) > 0 && len(a.state.items) > 0 && items[0].ID != a.state.items[0].ID)

			if needsUpdate {
				a.state.items = items

				// Update filtered items if not searching
				if a.state.searchQuery == "" {
					a.state.filteredItems = items
				} else {
					a.state.filteredItems = search.Filter(items, a.state.searchQuery)
				}

				if a.state.cursor >= len(a.state.filteredItems) {
					a.state.cursor = len(a.state.filteredItems) - 1
				}
				if a.state.cursor < 0 {
					a.state.cursor = 0
				}
			}
			a.state.mu.Unlock()

			if needsUpdate {
				a.app.QueueUpdateDraw(func() {
					a.updateListDisplay()
				})
			}
		}
	}
}

// setupGlobalKeys sets up global keyboard shortcuts
func (a *App) setupGlobalKeys() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Global quit on Ctrl+C
		if event.Key() == tcell.KeyCtrlC {
			a.app.Stop()
			return nil
		}

		// Mode-specific 'q' quit
		if event.Rune() == 'q' && a.state.currentMode != modeSearch {
			a.app.Stop()
			return nil
		}

		return event
	})
}

// switchToListMode switches to list mode
func (a *App) switchToListMode() {
	a.state.mu.Lock()
	a.state.currentMode = modeList
	a.state.mu.Unlock()

	a.pages.SwitchToPage("list")
	a.app.SetFocus(a.listWidget)
	a.updateListDisplay()
}

// switchToPreviewMode switches to preview mode
func (a *App) switchToPreviewMode() {
	a.state.mu.RLock()
	hasItems := len(a.state.filteredItems) > 0
	a.state.mu.RUnlock()

	if !hasItems {
		return
	}

	a.state.mu.Lock()
	a.state.currentMode = modePreview
	a.state.mu.Unlock()

	a.updatePreviewContent()
	a.pages.SwitchToPage("preview")
	a.app.SetFocus(a.previewView)
}

// switchToSearchMode switches to search mode (shows search box at bottom)
func (a *App) switchToSearchMode() {
	a.state.mu.Lock()
	a.state.currentMode = modeSearch
	a.state.searchQuery = ""
	a.state.mu.Unlock()

	a.searchInput.SetText("")

	// Replace help with search input
	a.mainFlex.RemoveItem(a.listHelp)
	a.mainFlex.AddItem(a.searchInput, searchInputHeight, 0, true)

	a.app.SetFocus(a.searchInput)
}

// exitSearchMode exits search mode and returns to list
func (a *App) exitSearchMode() {
	a.state.mu.Lock()
	a.state.currentMode = modeList
	a.state.mu.Unlock()

	// Replace search input with help
	a.mainFlex.RemoveItem(a.searchInput)
	a.mainFlex.AddItem(a.listHelp, searchInputHeight, 0, false)

	a.app.SetFocus(a.listWidget)
}

// handleCopyAction copies the selected item and quits
func (a *App) handleCopyAction() {
	a.state.mu.RLock()
	if len(a.state.filteredItems) == 0 {
		a.state.mu.RUnlock()
		return
	}
	item := a.state.filteredItems[a.state.cursor]
	a.state.mu.RUnlock()

	clipboard.SetClipboard(item.Content)
	a.app.Stop()
}

// handleDeleteAction deletes the selected item
func (a *App) handleDeleteAction() {
	a.state.mu.RLock()
	if len(a.state.filteredItems) == 0 {
		a.state.mu.RUnlock()
		return
	}
	itemID := a.state.filteredItems[a.state.cursor].ID
	a.state.mu.RUnlock()

	a.state.storage.Delete(itemID)
	a.reloadItems()
	a.updateListDisplay()
}

// handleClearAllAction clears all clipboard history
func (a *App) handleClearAllAction() {
	a.state.storage.Clear()

	a.state.mu.Lock()
	a.state.items = []types.ClipboardItem{}
	a.state.filteredItems = []types.ClipboardItem{}
	a.state.cursor = 0
	a.state.mu.Unlock()

	a.updateListDisplay()
}

// reloadItems reloads items from storage
func (a *App) reloadItems() {
	items, _ := a.state.storage.GetRecent(maxItemsToFetch)

	a.state.mu.Lock()
	defer a.state.mu.Unlock()

	a.state.items = items

	if a.state.searchQuery != "" {
		a.state.filteredItems = search.Filter(items, a.state.searchQuery)
	} else {
		a.state.filteredItems = items
	}

	if a.state.cursor >= len(a.state.filteredItems) {
		a.state.cursor = len(a.state.filteredItems) - 1
	}
	if a.state.cursor < 0 {
		a.state.cursor = 0
	}
}

// updateListDisplay updates the table widget with current items
func (a *App) updateListDisplay() {
	a.listWidget.Clear()

	a.state.mu.RLock()
	filteredItems := a.state.filteredItems
	cursor := a.state.cursor
	searchQuery := a.state.searchQuery
	currentMode := a.state.currentMode
	a.state.mu.RUnlock()

	if len(filteredItems) > 0 {
		title := fmt.Sprintf(" Clipboard History (%d/%d) ",
			cursor+1, len(filteredItems))
		a.listContainer.SetTitle(title)
	} else {
		a.listContainer.SetTitle(" Clipboard History (0) ")
	}

	if len(filteredItems) == 0 {
		var message string
		if searchQuery != "" {
			message = "No results found for '" + searchQuery + "'"
		} else {
			message = "No items in clipboard history"
		}

		a.listWidget.SetCell(0, 0, tview.NewTableCell(message).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorDefault).
			SetAttributes(tcell.AttrDim).
			SetSelectable(false))
		return
	}

	// Left spacer
	a.listWidget.SetCell(0, 0, tview.NewTableCell("").
		SetExpansion(2).
		SetSelectable(false))
	// Number column
	a.listWidget.SetCell(0, 1, tview.NewTableCell(" # ").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter).
		SetSelectable(false).
		SetExpansion(0).
		SetAttributes(tcell.AttrBold))
	// Content column
	a.listWidget.SetCell(0, 2, tview.NewTableCell("Content").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetSelectable(false).
		SetExpansion(3).
		SetAttributes(tcell.AttrBold))
	// Date column
	a.listWidget.SetCell(0, 3, tview.NewTableCell(fmt.Sprintf("%12s", "Date")).
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight).
		SetSelectable(false).
		SetExpansion(0).
		SetAttributes(tcell.AttrBold))
	// Right spacer
	a.listWidget.SetCell(0, 4, tview.NewTableCell("").
		SetExpansion(2).
		SetSelectable(false))

	for i, item := range filteredItems {
		row := i + 1 // +1 because row 0 is the header
		preview := truncate(item.Preview, previewTruncateLength)
		timestamp := formatTimestamp(item.Timestamp)

		// Left spacer
		a.listWidget.SetCell(row, 0, tview.NewTableCell("").
			SetExpansion(2))

		// Number column (show numbers 0-9 for first 10 items)
		var numStr string
		if i < 10 {
			numStr = fmt.Sprintf(" %d ", i)
		} else {
			numStr = "   "
		}
		a.listWidget.SetCell(row, 1, tview.NewTableCell(numStr).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorYellow).
			SetExpansion(0).
			SetAttributes(tcell.AttrBold))

		// Content column
		a.listWidget.SetCell(row, 2, tview.NewTableCell(preview).
			SetAlign(tview.AlignLeft).
			SetExpansion(3).
			SetTextColor(tcell.ColorDefault))

		// Date column
		a.listWidget.SetCell(row, 3, tview.NewTableCell(fmt.Sprintf("%12s", timestamp)).
			SetAlign(tview.AlignRight).
			SetExpansion(0).
			SetTextColor(tcell.ColorDefault).
			SetAttributes(tcell.AttrDim))

		// Right spacer
		a.listWidget.SetCell(row, 4, tview.NewTableCell("").
			SetExpansion(2))
	}

	if len(filteredItems) > 0 && currentMode != modeSearch {
		a.listWidget.Select(cursor+1, 1) // +1 for header row, column 1 is number column
	}
}

// updatePreviewContent updates the preview view with current item
func (a *App) updatePreviewContent() {
	a.state.mu.RLock()
	if len(a.state.filteredItems) == 0 || a.state.cursor >= len(a.state.filteredItems) {
		a.state.mu.RUnlock()
		a.previewView.SetText("No item selected")
		return
	}
	item := a.state.filteredItems[a.state.cursor]
	a.state.mu.RUnlock()

	timestamp := formatTimestamp(item.Timestamp)
	title := fmt.Sprintf(" Preview - %s • %d bytes • %s ",
		item.Type, len(item.Content), timestamp)
	a.previewView.SetTitle(title)

	content := FormatPreview(item.Content, item.Type, previewFormatMaxLength)
	a.previewView.SetText(content)
	a.previewView.ScrollToBeginning()
}

