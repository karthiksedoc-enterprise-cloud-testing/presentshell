package slides

import (
	"github.com/charmbracelet/glamour"
)

// Renderer renders markdown slides to styled terminal output.
type Renderer struct {
	width int
}

// NewRenderer creates a renderer with the given width constraint.
func NewRenderer(width int) *Renderer {
	return &Renderer{width: width}
}

// SetWidth updates the rendering width.
func (r *Renderer) SetWidth(width int) {
	r.width = width
}

// Render converts markdown content to styled terminal output.
func (r *Renderer) Render(markdown string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(r.width-4), // padding
	)
	if err != nil {
		return "", err
	}

	out, err := renderer.Render(markdown)
	if err != nil {
		return "", err
	}

	return out, nil
}
