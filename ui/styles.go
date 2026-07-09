package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED")
	secondaryColor = lipgloss.Color("#06B6D4")
	mutedColor     = lipgloss.Color("#6B7280")
	bgColor        = lipgloss.Color("#1F2937")

	// Pane styles
	ActiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor)

	InactiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(mutedColor)

	// Title styles
	SlideTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			PaddingLeft(1)

	TerminalTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(secondaryColor).
				PaddingLeft(1)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(lipgloss.Color("#E5E7EB")).
			PaddingLeft(1).
			PaddingRight(1)

	StatusKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor)

	StatusTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	// Help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			PaddingLeft(1)
)
