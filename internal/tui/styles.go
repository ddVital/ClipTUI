package tui

import (
	"fmt"
	"strings"
	"time"
)

// Mode types
type mode int

const (
	modeList mode = iota
	modePreview
	modeSearch
)

// formatTimestamp converts a timestamp to relative time string
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

// truncate shortens a string and replaces newlines
func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

