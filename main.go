package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karthik/presentshell/slides"
	"github.com/karthik/presentshell/terminal"
	"github.com/karthik/presentshell/ui"
)

// tickMsg triggers periodic terminal output refresh.
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// model is the main application model.
type model struct {
	presentation *slides.Presentation
	renderer     *slides.Renderer
	term         *terminal.Terminal
	layout       *ui.Layout
	currentSlide int
	totalSlides  int
	termOutput   string
	ready        bool
	quitting     bool
	err          error
}

func newModel(filepath string) model {
	pres, err := slides.Load(filepath)
	if err != nil {
		return model{err: err}
	}

	return model{
		presentation: pres,
		renderer:     slides.NewRenderer(40),
		layout:       ui.NewLayout(80, 24),
		totalSlides:  len(pres.Slides),
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.layout.Width = msg.Width
		m.layout.Height = msg.Height
		m.renderer.SetWidth(m.layout.LeftWidth() - 4)

		if !m.ready {
			cols := uint16(m.layout.RightWidth() - 4)
			rows := uint16(m.layout.ContentHeight() - 2)
			t, err := terminal.New(rows, cols)
			if err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.term = t
			m.ready = true
		} else if m.term != nil {
			cols := uint16(m.layout.RightWidth() - 4)
			rows := uint16(m.layout.ContentHeight() - 2)
			m.term.Resize(rows, cols)
		}
		return m, nil

	case tickMsg:
		if m.term != nil {
			raw := m.term.Read()
			m.termOutput = sanitizeOutput(string(raw))
		}
		return m, tickCmd()

	case tea.KeyMsg:
		// --- Terminal pane focused: forward everything except Tab ---
		if m.layout.FocusedPane == ui.TerminalPane {
			switch msg.String() {
			case "tab":
				m.layout.FocusedPane = ui.SlidePane
				return m, nil
			default:
				if m.term != nil {
					if b := keyToBytes(msg); len(b) > 0 {
						m.term.Write(b)
					}
				}
				return m, nil
			}
		}

		// --- Slide pane focused ---
		switch msg.String() {
		case "tab":
			m.layout.FocusedPane = ui.TerminalPane
			return m, nil

		case "right", "l", "n", " ":
			if m.currentSlide < m.totalSlides-1 {
				m.currentSlide++
			}
			return m, nil

		case "left", "h", "p":
			if m.currentSlide > 0 {
				m.currentSlide--
			}
			return m, nil

		case "q", "esc":
			m.quitting = true
			if m.term != nil {
				m.term.Close()
			}
			return m, tea.Quit

		case "ctrl+c":
			m.quitting = true
			if m.term != nil {
				m.term.Close()
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.err != nil {
		return fmt.Sprintf("\n  Error: %v\n\n  Press Ctrl+C to exit.\n", m.err)
	}
	if !m.ready {
		return "\n  Initializing PresentShell...\n"
	}

	// Render current slide markdown
	slideContent := ""
	if m.totalSlides > 0 {
		rendered, err := m.renderer.Render(m.presentation.Slides[m.currentSlide])
		if err != nil {
			slideContent = m.presentation.Slides[m.currentSlide]
		} else {
			slideContent = rendered
		}
	}

	// Terminal output (last N lines)
	termContent := getLastLines(m.termOutput, m.layout.ContentHeight()-2)

	// Split view
	view := m.layout.RenderSplitView(
		slideContent,
		termContent,
		fmt.Sprintf("📄 %s", m.presentation.Title),
		"🖥  Terminal",
	)

	// Status bar
	statusBar := ui.RenderStatusBar(
		m.currentSlide,
		m.totalSlides,
		m.layout.FocusedPane,
		m.layout.Width,
	)

	return view + "\n" + statusBar
}

// keyToBytes converts a key event to raw bytes for the PTY.
func keyToBytes(msg tea.KeyMsg) []byte {
	switch msg.Type {
	case tea.KeyEnter:
		return []byte{'\r'}
	case tea.KeyBackspace:
		return []byte{127}
	case tea.KeySpace:
		return []byte{' '}
	case tea.KeyTab:
		return []byte{'\t'}
	case tea.KeyEsc:
		return []byte{27}
	case tea.KeyUp:
		return []byte{27, '[', 'A'}
	case tea.KeyDown:
		return []byte{27, '[', 'B'}
	case tea.KeyRight:
		return []byte{27, '[', 'C'}
	case tea.KeyLeft:
		return []byte{27, '[', 'D'}
	case tea.KeyCtrlC:
		return []byte{3}
	case tea.KeyCtrlD:
		return []byte{4}
	case tea.KeyCtrlZ:
		return []byte{26}
	case tea.KeyCtrlL:
		return []byte{12}
	case tea.KeyRunes:
		return []byte(string(msg.Runes))
	default:
		return nil
	}
}

func sanitizeOutput(s string) string {
	s = strings.ReplaceAll(s, "\x1b]0;", "")
	return s
}

func getLastLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: presentshell <presentation.md>")
		fmt.Println("\nExample: presentshell examples/demo.md")
		os.Exit(1)
	}

	filepath := os.Args[1]

	p := tea.NewProgram(
		newModel(filepath),
		tea.WithAltScreen(),
		tea.WithInputTTY(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
