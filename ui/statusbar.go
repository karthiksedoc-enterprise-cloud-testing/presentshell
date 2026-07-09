package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar creates the bottom status bar.
func RenderStatusBar(current, total int, focused Pane, width int) string {
	slideInfo := StatusKeyStyle.Render(fmt.Sprintf(" Slide %d/%d ", current+1, total))

	var focusInfo string
	if focused == SlidePane {
		focusInfo = StatusTextStyle.Render("│ 📝 Slides focused ")
	} else {
		focusInfo = StatusTextStyle.Render("│ 💻 Terminal focused ")
	}

	help := HelpStyle.Render("Tab:switch  ←→:navigate  q:quit")

	left := slideInfo + focusInfo
	right := help

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	spacer := StatusBarStyle.Render(fmt.Sprintf("%*s", gap, ""))

	bar := StatusBarStyle.Width(width).Render(left + spacer + right)
	return bar
}
