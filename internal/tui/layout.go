package tui

import (
	"github.com/dvd/cliptui/internal/clipboard"
	"github.com/dvd/cliptui/internal/search"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// moveCursorUp moves the cursor up by one item
func (a *App) moveCursorUp() {
	row, _ := a.listWidget.GetSelection()
	if row > 1 {
		a.listWidget.Select(row-1, 1)
		a.state.mu.Lock()
		a.state.cursor = row - 2
		a.state.mu.Unlock()
		a.updateListDisplay()
	}
}

// moveCursorDown moves the cursor down by one item
func (a *App) moveCursorDown() {
	row, _ := a.listWidget.GetSelection()
	if row < a.listWidget.GetRowCount()-1 {
		a.listWidget.Select(row+1, 1)
		a.state.mu.Lock()
		a.state.cursor = row
		a.state.mu.Unlock()
		a.updateListDisplay()
	}
}

// buildListPage creates the list mode layout
func (a *App) buildListPage() tview.Primitive {
	// Table widget - use terminal default colors
	a.listWidget = tview.NewTable().
		SetFixed(1, 0). // Fix the header row
		SetSelectable(true, false).
		SetSelectedStyle(tcell.StyleDefault.
			Background(tcell.ColorDefault).
			Foreground(tcell.ColorDefault).
			Reverse(true)). // Use reverse video for selection
		SetSeparator(' ')
	a.listWidget.SetBorder(false)

	a.listWidget.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseScrollUp {
			a.moveCursorUp()
			return action, nil
		}
		if action == tview.MouseScrollDown {
			a.moveCursorDown()
			return action, nil
		}
		return action, event
	})

	a.listWidget.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			a.moveCursorDown()
			return nil
		case 'k':
			a.moveCursorUp()
			return nil
		case 'p':
			a.switchToPreviewMode()
			return nil
		case '/':
			a.switchToSearchMode()
			return nil
		case 'd':
			a.handleDeleteAction()
			return nil
		case 'D':
			a.handleClearAllAction()
			return nil
		case 'y':
			a.handleCopyAction()
			return nil
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// Quick copy by number (0 = first item, 9 = tenth item)
			num := int(event.Rune() - '0')
			a.state.mu.RLock()
			if num < len(a.state.filteredItems) {
				item := a.state.filteredItems[num]
				a.state.mu.RUnlock()
				clipboard.SetClipboard(item.Content)
				a.app.Stop()
			} else {
				a.state.mu.RUnlock()
			}
			return nil
		}

		if event.Key() == tcell.KeyEnter {
			a.handleCopyAction()
			return nil
		}

		if event.Key() == tcell.KeyDown {
			a.moveCursorDown()
			return nil
		}
		if event.Key() == tcell.KeyUp {
			a.moveCursorUp()
			return nil
		}

		return event
	})

	a.listHelp = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	a.listHelp.SetText("  0-9 quick copy • ↑/k up • ↓/j down • enter/y copy • p preview • / search • d delete • D clear • q quit")
	a.listHelp.SetBorder(true).
		SetTitle(" Shortcuts ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorBlue)

	a.listContainer = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.listWidget, 0, 1, true)
	a.listContainer.SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen).
		SetBorderPadding(0, 0, 1, 1).
		SetTitle(" Clipboard History ").
		SetTitleAlign(tview.AlignLeft)

	a.listHelp.SetBorderPadding(0, 0, 1, 1)

	a.searchInput = tview.NewInputField().
		SetLabel("").
		SetPlaceholder("Type to search...").
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorWhite).
		SetPlaceholderTextColor(tcell.ColorGray)

	a.searchInput.SetBorder(true).
		SetBorderColor(tcell.ColorYellow).
		SetTitleColor(tcell.ColorYellow).
		SetTitle(" Search (ESC to cancel, Enter to confirm) ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	a.searchInput.SetChangedFunc(func(text string) {
		a.state.mu.Lock()
		a.state.searchQuery = text
		if text == "" {
			a.state.filteredItems = a.state.items
		} else {
			a.state.filteredItems = search.Filter(a.state.items, text)
		}
		a.state.cursor = 0
		a.state.mu.Unlock()
		a.updateListDisplay()
	})

	a.searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.state.mu.Lock()
			a.state.filteredItems = a.state.items
			a.state.searchQuery = ""
			a.state.cursor = 0
			a.state.mu.Unlock()
			a.exitSearchMode()
			a.updateListDisplay()
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			a.exitSearchMode()
			return nil
		}
		return event
	})

	a.mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.listContainer, 0, 1, true).
		AddItem(a.listHelp, searchInputHeight, 0, false)

	outer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 2, 0, false).
			AddItem(a.mainFlex, 0, 1, true).
			AddItem(nil, 2, 0, false),
			0, 1, true).
		AddItem(nil, 1, 0, false)

	return outer
}

// buildPreviewPage creates the preview mode layout
func (a *App) buildPreviewPage() tview.Primitive {
	a.previewHeader = tview.NewTextView()

	a.previewView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetTextColor(tcell.ColorDefault)
	a.previewView.SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetTitleColor(tcell.ColorGreen).
		SetBorderPadding(1, 0, 1, 1).
		SetTitle(" Preview ").
		SetTitleAlign(tview.AlignLeft)

	a.previewView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			a.switchToListMode()
			return nil
		}
		if event.Key() == tcell.KeyEnter || event.Rune() == 'y' {
			a.handleCopyAction()
			return nil
		}
		return event
	})

	a.previewHelp = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	a.previewHelp.SetText("  enter/y copy • esc/q back • ↑↓ scroll")
	a.previewHelp.SetBorder(true).
		SetTitle(" Shortcuts ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.ColorBlue).
		SetTitleColor(tcell.ColorBlue).
		SetBorderPadding(0, 0, 1, 1)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.previewView, 0, 1, true).
		AddItem(a.previewHelp, searchInputHeight, 0, false)

	outer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 2, 0, false).
			AddItem(flex, 0, 1, true).
			AddItem(nil, 2, 0, false),
			0, 1, true).
		AddItem(nil, 1, 0, false)

	return outer
}
