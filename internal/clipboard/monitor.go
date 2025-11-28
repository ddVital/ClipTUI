package clipboard

import (
	"context"
	"time"

	"github.com/atotto/clipboard"
	"github.com/dvd/cliptui/internal/storage"
)

// Monitor watches the clipboard for changes
type Monitor struct {
	storage      *storage.Storage
	pollInterval time.Duration
	lastContent  string
}

// NewMonitor creates a new clipboard monitor
func NewMonitor(store *storage.Storage, pollInterval time.Duration) *Monitor {
	return &Monitor{
		storage:      store,
		pollInterval: pollInterval,
		lastContent:  "",
	}
}

// Start begins monitoring the clipboard
func (m *Monitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(m.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			content, err := clipboard.ReadAll()
			if err != nil {
				// Clipboard read error, continue
				continue
			}

			// Only store if content changed and is not empty
			if content != "" && content != m.lastContent {
				if err := m.storage.Add(content); err != nil {
					// Log error but continue monitoring
					continue
				}
				m.lastContent = content
			}
		}
	}
}

// SetClipboard sets the system clipboard content
func SetClipboard(content string) error {
	return clipboard.WriteAll(content)
}
