package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	secondaryColor = lipgloss.Color("#A78BFA")
	accentColor    = lipgloss.Color("#EC4899") // Pink
	textColor      = lipgloss.Color("#E5E7EB")
	mutedColor     = lipgloss.Color("#9CA3AF")
	bgColor        = lipgloss.Color("#1F2937")
	selectedBg     = lipgloss.Color("#374151")

	// Type colors
	typeColors = map[string]lipgloss.Color{
		"text":     lipgloss.Color("#60A5FA"), // Blue
		"code":     lipgloss.Color("#34D399"), // Green
		"markdown": lipgloss.Color("#FBBF24"), // Yellow
		"url":      lipgloss.Color("#F472B6"), // Pink
	}

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(textColor)

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Item styles
	itemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(textColor)

	selectedItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(textColor).
				Background(selectedBg).
				Bold(true)

	// Type badge style
	typeBadgeStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)

	// Preview style
	previewStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(0, 2)

	// Timestamp style
	timestampStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Align(lipgloss.Right)

	// Search input style
	searchStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Padding(0, 1)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(1, 2)

	// Border style
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	// Status bar style
	statusStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(bgColor).
			Padding(0, 1)
)

func getTypeBadge(itemType string) string {
	color, ok := typeColors[itemType]
	if !ok {
		color = mutedColor
	}

	style := typeBadgeStyle.Copy().Foreground(color)
	return style.Render(itemType)
}
