package types

import "time"

// ClipboardItem represents a single clipboard entry
type ClipboardItem struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // text, code, markdown, url
	Timestamp time.Time `json:"timestamp"`
	Preview   string    `json:"preview"` // truncated version for list view
}

// ItemType constants
const (
	TypeText     = "text"
	TypeCode     = "code"
	TypeMarkdown = "markdown"
	TypeURL      = "url"
)

// DetectType attempts to determine the content type
func DetectType(content string) string {
	if len(content) == 0 {
		return TypeText
	}

	// Simple heuristics
	if isURL(content) {
		return TypeURL
	}
	if isCode(content) {
		return TypeCode
	}
	if isMarkdown(content) {
		return TypeMarkdown
	}
	return TypeText
}

func isURL(s string) bool {
	return len(s) < 2048 && (len(s) > 7 && (s[:7] == "http://" || s[:8] == "https://"))
}

func isCode(s string) bool {
	// Check for common code patterns
	codeIndicators := []string{"{", "}", "func ", "def ", "class ", "import ", "package ", "const ", "var ", "let "}
	for _, indicator := range codeIndicators {
		if contains(s, indicator) {
			return true
		}
	}
	return false
}

func isMarkdown(s string) bool {
	// Check for markdown patterns
	mdIndicators := []string{"# ", "## ", "- ", "* ", "```", "**", "__"}
	for _, indicator := range mdIndicators {
		if contains(s, indicator) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TruncatePreview creates a preview string
func TruncatePreview(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "..."
}
