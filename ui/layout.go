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

	// Build left pane
	leftTitleStr := SlideTitleStyle.Render(leftTitle)
	leftBody := padOrTruncate(leftContent, leftW, contentH-1)
	leftFull := leftTitleStr + "\n" + leftBody

	// Build right pane
	rightTitleStr := TerminalTitleStyle.Render(rightTitle)
	rightBody := padOrTruncate(rightContent, rightW, contentH-1)
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
func padOrTruncate(content string, width, height int) string {
	lines := strings.Split(content, "\n")

	// Truncate to height
	if len(lines) > height {
		lines = lines[len(lines)-height:]
	}

	// Pad to height
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Truncate each line to width
	for i, line := range lines {
		if lipgloss.Width(line) > width {
			lines[i] = line[:width]
		}
	}

	return strings.Join(lines, "\n")
}
