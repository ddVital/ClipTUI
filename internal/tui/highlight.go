package tui

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/dvd/cliptui/pkg/types"
)

// HighlightContent applies syntax highlighting to content based on type
func HighlightContent(content string, itemType string) string {
	if itemType != types.TypeCode {
		return content
	}

	// Try to detect lexer from content
	lexer := lexers.Analyse(content)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Use a dark theme
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return content
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return content
	}

	return buf.String()
}

// FormatPreview formats the preview with line numbers and highlighting
func FormatPreview(content string, itemType string, maxLines int) string {
	lines := strings.Split(content, "\n")

	// Limit lines if needed
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, "...")
	}

	// For code, attempt highlighting
	if itemType == types.TypeCode {
		highlighted := HighlightContent(strings.Join(lines, "\n"), itemType)
		return highlighted
	}

	return strings.Join(lines, "\n")
}
