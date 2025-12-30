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
				continue
			}

			if content != "" && content != m.lastContent {
				latest, err := m.storage.GetLatest()
				if err == nil && latest != nil && latest.Content == content {
					m.lastContent = content
					continue
				}

				if err := m.storage.Add(content); err != nil {
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
