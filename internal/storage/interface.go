package storage

import "github.com/dvd/cliptui/pkg/types"

// Store defines the interface for clipboard history storage
// This interface allows for easier testing with mock implementations
type Store interface {
	// Add inserts a new clipboard item
	Add(content string) error

	// GetAll retrieves all clipboard items, newest first
	GetAll() ([]types.ClipboardItem, error)

	// GetRecent retrieves the N most recent items
	GetRecent(limit int) ([]types.ClipboardItem, error)

	// Delete removes an item by ID
	Delete(id int64) error

	// Clear removes all items
	Clear() error

	// GetLatest returns the most recent item
	GetLatest() (*types.ClipboardItem, error)

	// Close closes the database connection
	Close() error
}
