package search

import (
	"strings"
	"github.com/dvd/cliptui/pkg/types"
)

// Filter performs case-insensitive substring search on clipboard items
func Filter(items []types.ClipboardItem, query string) []types.ClipboardItem {
	if query == "" {
		return items
	}

	if len(items) == 0 {
		return items
	}

	query = strings.ToLower(query)
	results := make([]types.ClipboardItem, 0)

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Content), query) {
			results = append(results, item)
		}
	}

	return results
}
