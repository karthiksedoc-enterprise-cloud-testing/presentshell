package slides

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

// Renderer renders markdown slides to styled terminal output.
type Renderer struct {
	width    int
	renderer *glamour.TermRenderer
}

// NewRenderer creates a renderer with the given width constraint.
func NewRenderer(width int) *Renderer {
	r := &Renderer{width: width}
	r.buildRenderer()
	return r
}

// SetWidth updates the rendering width and rebuilds the internal renderer.
func (r *Renderer) SetWidth(width int) {
	if width != r.width {
		r.width = width
		r.buildRenderer()
	}
}

func (r *Renderer) buildRenderer() {
	tr, err := glamour.NewTermRenderer(
		glamour.WithStyles(styles.DarkStyleConfig),
		glamour.WithWordWrap(r.width-4),
	)
	if err == nil {
		r.renderer = tr
	}
}

// Render converts markdown content to styled terminal output.
func (r *Renderer) Render(markdown string) (string, error) {
	if r.renderer == nil {
		return markdown, nil
	}
	return r.renderer.Render(markdown)
}
