package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Pane represents which pane is currently focused.
type Pane int

const (
	SlidePane    Pane = iota
	TerminalPane Pane = iota
)

// Layout manages the split-pane layout dimensions.
type Layout struct {
	Width       int
	Height      int
	SplitRatio  float64 // ratio of left pane (0.0 to 1.0)
	FocusedPane Pane
}

// NewLayout creates a layout with default 50/50 split.
func NewLayout(width, height int) *Layout {
	return &Layout{
		Width:       width,
		Height:      height,
		SplitRatio:  0.5,
		FocusedPane: SlidePane,
	}
}

// LeftWidth returns the width of the left (slides) pane.
func (l *Layout) LeftWidth() int {
	return int(float64(l.Width) * l.SplitRatio)
}

// RightWidth returns the width of the right (terminal) pane.
func (l *Layout) RightWidth() int {
	return l.Width - l.LeftWidth()
}

// ContentHeight returns the height available for pane content (minus borders and status bar).
func (l *Layout) ContentHeight() int {
	return l.Height - 4 // 2 for border top/bottom, 1 for title, 1 for status bar
}

// RenderSplitView renders the two panes side by side.
func (l *Layout) RenderSplitView(leftContent, rightContent, leftTitle, rightTitle string) string {
	leftW := l.LeftWidth() - 2  // account for border
	rightW := l.RightWidth() - 2 // account for border
	contentH := l.ContentHeight()

	// Build left pane (slides - centered both vertically and horizontally)
	leftTitleStr := SlideTitleStyle.Render(leftTitle)
	leftBody := lipgloss.Place(leftW, contentH-1, lipgloss.Center, lipgloss.Center, leftContent)
	leftFull := leftTitleStr + "\n" + leftBody

	// Build right pane (terminal - top aligned)
	rightTitleStr := TerminalTitleStyle.Render(rightTitle)
	rightBody := padOrTruncate(rightContent, rightW, contentH-1, false)
	rightFull := rightTitleStr + "\n" + rightBody

	// Apply border styles based on focus
	var leftStyle, rightStyle lipgloss.Style
	if l.FocusedPane == SlidePane {
		leftStyle = ActiveBorderStyle.Width(leftW).Height(contentH)
		rightStyle = InactiveBorderStyle.Width(rightW).Height(contentH)
	} else {
		leftStyle = InactiveBorderStyle.Width(leftW).Height(contentH)
		rightStyle = ActiveBorderStyle.Width(rightW).Height(contentH)
	}

	leftPane := leftStyle.Render(leftFull)
	rightPane := rightStyle.Render(rightFull)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

// padOrTruncate ensures content fits within width x height.
// If centerVertical is true, content is vertically centered.
func padOrTruncate(content string, width, height int, centerVertical bool) string {
	lines := strings.Split(content, "\n")

	// Remove trailing empty lines for accurate content height
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// Truncate to height if too tall
	if len(lines) > height {
		lines = lines[:height]
	}

	if centerVertical && len(lines) < height {
		// Center vertically: add padding above
		topPad := (height - len(lines)) / 2
		padded := make([]string, 0, height)
		for i := 0; i < topPad; i++ {
			padded = append(padded, "")
		}
		padded = append(padded, lines...)
		// Fill remaining below
		for len(padded) < height {
			padded = append(padded, "")
		}
		lines = padded
	} else {
		// Pad to height at bottom
		for len(lines) < height {
			lines = append(lines, "")
		}
	}

	// Truncate each line to width
	for i, line := range lines {
		if lipgloss.Width(line) > width {
			lines[i] = line[:width]
		}
	}

	return strings.Join(lines, "\n")
}
