package slides

import (
	"os"
	"strings"
)

// Presentation holds all slides parsed from a markdown file.
type Presentation struct {
	Slides []string
	Title  string
}

// Load reads a markdown file and splits it into slides using "---" as delimiter.
func Load(filepath string) (*Presentation, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	rawSlides := strings.Split(content, "\n---\n")

	slides := make([]string, 0, len(rawSlides))
	for _, s := range rawSlides {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			slides = append(slides, trimmed)
		}
	}

	title := extractTitle(slides)

	return &Presentation{
		Slides: slides,
		Title:  title,
	}, nil
}

// extractTitle tries to find the first H1 heading in the first slide.
func extractTitle(slides []string) string {
	if len(slides) == 0 {
		return "Presentation"
	}
	lines := strings.Split(slides[0], "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimPrefix(trimmed, "# ")
		}
	}
	return "Presentation"
}
