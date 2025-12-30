package tui

import (
	"fmt"
	"time"

	"github.com/dvd/cliptui/internal/clipboard"
	"github.com/dvd/cliptui/internal/search"
	"github.com/dvd/cliptui/internal/storage"
	"github.com/dvd/cliptui/pkg/types"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AppState holds the application state
type AppState struct {
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
	items, err := store.GetRecent(100)
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
	stopChan := make(chan bool)
	go a.monitorClipboard(stopChan)

	err := a.app.Run()

	stopChan <- true

	return err
}

// monitorClipboard checks for new clipboard items periodically
func (a *App) monitorClipboard(stopChan chan bool) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			items, err := a.state.storage.GetRecent(100)
			if err != nil {
				continue
			}

			if len(items) != len(a.state.items) || (len(items) > 0 && len(a.state.items) > 0 && items[0].ID != a.state.items[0].ID) {
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
	a.state.currentMode = modeList
	a.pages.SwitchToPage("list")
	a.app.SetFocus(a.listWidget)
	a.updateListDisplay()
}

// switchToPreviewMode switches to preview mode
func (a *App) switchToPreviewMode() {
	if len(a.state.filteredItems) == 0 {
		return
	}

	a.state.currentMode = modePreview
	a.updatePreviewContent()
	a.pages.SwitchToPage("preview")
	a.app.SetFocus(a.previewView)
}

// switchToSearchMode switches to search mode (shows search box at bottom)
func (a *App) switchToSearchMode() {
	a.state.currentMode = modeSearch
	a.searchInput.SetText("")
	a.state.searchQuery = ""

	// Replace help with search input
	a.mainFlex.RemoveItem(a.listHelp)
	a.mainFlex.AddItem(a.searchInput, 3, 0, true)

	a.app.SetFocus(a.searchInput)
}

// exitSearchMode exits search mode and returns to list
func (a *App) exitSearchMode() {
	a.state.currentMode = modeList

	// Replace search input with help
	a.mainFlex.RemoveItem(a.searchInput)
	a.mainFlex.AddItem(a.listHelp, 3, 0, false)

	a.app.SetFocus(a.listWidget)
}

// handleCopyAction copies the selected item and quits
func (a *App) handleCopyAction() {
	if len(a.state.filteredItems) == 0 {
		return
	}

	item := a.state.filteredItems[a.state.cursor]
	clipboard.SetClipboard(item.Content)
	a.app.Stop()
}

// handleDeleteAction deletes the selected item
func (a *App) handleDeleteAction() {
	if len(a.state.filteredItems) == 0 {
		return
	}

	item := a.state.filteredItems[a.state.cursor]
	a.state.storage.Delete(item.ID)

	a.reloadItems()

	if a.state.cursor >= len(a.state.filteredItems) && a.state.cursor > 0 {
		a.state.cursor--
	}

	a.updateListDisplay()
}

// handleClearAllAction clears all clipboard history
func (a *App) handleClearAllAction() {
	a.state.storage.Clear()
	a.state.items = []types.ClipboardItem{}
	a.state.filteredItems = []types.ClipboardItem{}
	a.state.cursor = 0
	a.updateListDisplay()
}

// reloadItems reloads items from storage
func (a *App) reloadItems() {
	items, _ := a.state.storage.GetRecent(100)
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

	if len(a.state.filteredItems) > 0 {
		title := fmt.Sprintf(" Clipboard History (%d/%d) ",
			a.state.cursor+1, len(a.state.filteredItems))
		a.listContainer.SetTitle(title)
	} else {
		a.listContainer.SetTitle(" Clipboard History (0) ")
	}

	if len(a.state.filteredItems) == 0 {
		var message string
		if a.state.searchQuery != "" {
			message = "No results found for '" + a.state.searchQuery + "'"
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

	a.listWidget.SetCell(0, 0, tview.NewTableCell("Content").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignLeft).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold))
	a.listWidget.SetCell(0, 1, tview.NewTableCell("Date").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight).
		SetSelectable(false).
		SetAttributes(tcell.AttrBold))

	for i, item := range a.state.filteredItems {
		row := i + 1 // +1 because row 0 is the header
		preview := truncate(item.Preview, 80)
		timestamp := formatTimestamp(item.Timestamp)

		a.listWidget.SetCell(row, 0, tview.NewTableCell(preview).
			SetAlign(tview.AlignLeft).
			SetTextColor(tcell.ColorDefault))

		a.listWidget.SetCell(row, 1, tview.NewTableCell(timestamp).
			SetAlign(tview.AlignRight).
			SetTextColor(tcell.ColorDefault).
			SetAttributes(tcell.AttrDim))
	}

	if len(a.state.filteredItems) > 0 && a.state.currentMode != modeSearch {
		a.listWidget.Select(a.state.cursor+1, 0) // +1 for header row
	}
}

// updatePreviewContent updates the preview view with current item
func (a *App) updatePreviewContent() {
	if len(a.state.filteredItems) == 0 || a.state.cursor >= len(a.state.filteredItems) {
		a.previewView.SetText("No item selected")
		return
	}

	item := a.state.filteredItems[a.state.cursor]

	timestamp := formatTimestamp(item.Timestamp)
	title := fmt.Sprintf(" Preview - %s • %d bytes • %s ",
		item.Type, len(item.Content), timestamp)
	a.previewView.SetTitle(title)

	content := FormatPreview(item.Content, item.Type, 1000)
	a.previewView.SetText(content)
	a.previewView.ScrollToBeginning()
}

