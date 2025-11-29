package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	DBPath      string
	MaxItems    int
	PollInterval int // milliseconds
}

// Default returns default configuration
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".local", "share", "cliptui")

	// Ensure directory exists
	os.MkdirAll(dataDir, 0755)

	return &Config{
		DBPath:      filepath.Join(dataDir, "clipboard.db"),
		MaxItems:    1000,
		PollInterval: 500, // 500ms
	}
}
