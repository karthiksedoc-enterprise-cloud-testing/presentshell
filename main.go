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

// Model is the main application model.
type Model struct {
	presentation *slides.Presentation
	renderer     *slides.Renderer
	term         *terminal.Terminal
	layout       *ui.Layout
	currentSlide int
	termOutput   string
	ready        bool
	err          error
}

func initialModel(filepath string) Model {
	pres, err := slides.Load(filepath)
	if err != nil {
		return Model{err: err}
	}

	return Model{
		presentation: pres,
		renderer:     slides.NewRenderer(40),
		layout:       ui.NewLayout(80, 24),
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout.Width = msg.Width
		m.layout.Height = msg.Height
		m.renderer.SetWidth(m.layout.LeftWidth() - 4)

		if !m.ready {
			// Start terminal with right pane dimensions
			cols := uint16(m.layout.RightWidth() - 4)
			rows := uint16(m.layout.ContentHeight() - 2)
			term, err := terminal.New(rows, cols)
			if err != nil {
				m.err = err
				return m, tea.Quit
			}
			m.term = term
			m.ready = true
		} else if m.term != nil {
			cols := uint16(m.layout.RightWidth() - 4)
			rows := uint16(m.layout.ContentHeight() - 2)
			m.term.Resize(rows, cols)
		}
		return m, nil

	case tickMsg:
		// Refresh terminal output
		if m.term != nil {
			raw := m.term.Read()
			m.termOutput = sanitizeOutput(string(raw))
		}
		return m, tickCmd()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// If terminal is focused, forward keys to PTY (except Tab)
	if m.layout.FocusedPane == ui.TerminalPane {
		if key == "tab" {
			m.layout.FocusedPane = ui.SlidePane
			return m, nil
		}
		// Forward all other keys to the terminal
		if m.term != nil {
			input := keyToBytes(msg)
			if len(input) > 0 {
				m.term.Write(input)
			}
		}
		return m, nil
	}

	// Slide pane is focused
	switch key {
	case "tab":
		m.layout.FocusedPane = ui.TerminalPane
		return m, nil
	case "right", "l", "n":
		m.nextSlide()
		return m, nil
	case "left", "h", "p":
		m.prevSlide()
		return m, nil
	case "q", "ctrl+c":
		if m.term != nil {
			m.term.Close()
		}
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) nextSlide() {
	if m.presentation != nil && m.currentSlide < len(m.presentation.Slides)-1 {
		m.currentSlide++
	}
}

func (m *Model) prevSlide() {
	if m.currentSlide > 0 {
		m.currentSlide--
	}
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress any key to exit.", m.err)
	}

	if !m.ready {
		return "Initializing..."
	}

	// Render current slide
	slideContent := ""
	if m.presentation != nil && len(m.presentation.Slides) > 0 {
		rendered, err := m.renderer.Render(m.presentation.Slides[m.currentSlide])
		if err != nil {
			slideContent = m.presentation.Slides[m.currentSlide]
		} else {
			slideContent = rendered
		}
	}

	// Get terminal content (last N lines)
	termContent := getLastLines(m.termOutput, m.layout.ContentHeight()-2)

	// Render split view
	view := m.layout.RenderSplitView(
		slideContent,
		termContent,
		fmt.Sprintf("📄 %s", m.presentation.Title),
		"🖥  Terminal",
	)

	// Add status bar
	statusBar := ui.RenderStatusBar(
		m.currentSlide,
		len(m.presentation.Slides),
		m.layout.FocusedPane,
		m.layout.Width,
	)

	return view + "\n" + statusBar
}

// keyToBytes converts a tea.KeyMsg to bytes for the PTY.
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

// sanitizeOutput strips problematic escape sequences for display.
func sanitizeOutput(s string) string {
	// Keep ANSI colors but strip title-setting sequences
	s = strings.ReplaceAll(s, "\x1b]0;", "")
	return s
}

// getLastLines returns the last n lines of a string.
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
		initialModel(filepath),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
