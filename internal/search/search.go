package search

import (
	"github.com/sahilm/fuzzy"
	"github.com/dvd/cliptui/pkg/types"
)

// Filter performs fuzzy search on clipboard items
func Filter(items []types.ClipboardItem, query string) []types.ClipboardItem {
	if query == "" {
		return items
	}

	// Build searchable strings
	searchStrings := make([]string, len(items))
	for i, item := range items {
		searchStrings[i] = item.Content
	}

	// Perform fuzzy search
	matches := fuzzy.Find(query, searchStrings)

	// Build result list maintaining order of matches
	results := make([]types.ClipboardItem, len(matches))
	for i, match := range matches {
		results[i] = items[match.Index]
	}

	return results
}
