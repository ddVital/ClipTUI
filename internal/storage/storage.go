package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/dvd/cliptui/pkg/types"
)

// Storage handles clipboard history persistence
type Storage struct {
	db *sql.DB
}

// New creates a new storage instance
func New(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS clipboard_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			type TEXT NOT NULL,
			preview TEXT NOT NULL,
			timestamp DATETIME NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_timestamp ON clipboard_history(timestamp DESC);
	`)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

// Add inserts a new clipboard item
func (s *Storage) Add(content string) error {
	latest, err := s.GetLatest()
	if err == nil && latest != nil && latest.Content == content {
		return nil
	}

	itemType := types.DetectType(content)
	preview := types.TruncatePreview(content, 100)

	_, err = s.db.Exec(
		"INSERT INTO clipboard_history (content, type, preview, timestamp) VALUES (?, ?, ?, ?)",
		content, itemType, preview, time.Now(),
	)
	return err
}

// GetAll retrieves all clipboard items, newest first
func (s *Storage) GetAll() ([]types.ClipboardItem, error) {
	rows, err := s.db.Query(`
		SELECT id, content, type, preview, timestamp
		FROM clipboard_history
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.ClipboardItem
	for rows.Next() {
		var item types.ClipboardItem
		err := rows.Scan(&item.ID, &item.Content, &item.Type, &item.Preview, &item.Timestamp)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// GetRecent retrieves the N most recent items
func (s *Storage) GetRecent(limit int) ([]types.ClipboardItem, error) {
	rows, err := s.db.Query(`
		SELECT id, content, type, preview, timestamp
		FROM clipboard_history
		ORDER BY timestamp DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.ClipboardItem
	for rows.Next() {
		var item types.ClipboardItem
		err := rows.Scan(&item.ID, &item.Content, &item.Type, &item.Preview, &item.Timestamp)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Delete removes an item by ID
func (s *Storage) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM clipboard_history WHERE id = ?", id)
	return err
}

// Clear removes all items
func (s *Storage) Clear() error {
	_, err := s.db.Exec("DELETE FROM clipboard_history")
	return err
}

// GetLatest returns the most recent item
func (s *Storage) GetLatest() (*types.ClipboardItem, error) {
	var item types.ClipboardItem
	err := s.db.QueryRow(`
		SELECT id, content, type, preview, timestamp
		FROM clipboard_history
		ORDER BY timestamp DESC
		LIMIT 1
	`).Scan(&item.ID, &item.Content, &item.Type, &item.Preview, &item.Timestamp)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}
