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
	// Pre-allocate slice with capacity to avoid reallocations
	// In worst case, all items match, so use len(items) as capacity
	results := make([]types.ClipboardItem, 0, len(items))

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Content), query) {
			results = append(results, item)
		}
	}

	return results
}
