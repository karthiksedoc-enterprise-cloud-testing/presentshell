# PresentShell 🖥️

A terminal-based presentation tool with a split-pane layout — render beautiful markdown slides on the left while running live demos in a real terminal on the right.

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

- 📄 **Markdown slides** — Write presentations in plain markdown, separated by `---`
- 🖥 **Live terminal** — Embedded PTY shell for real-time demos (SSH, Docker, scripts)
- 🎨 **Syntax highlighting** — Code blocks rendered with full color support
- ⌨️ **Keyboard driven** — Navigate slides and switch focus with simple key bindings
- 📐 **Responsive** — Adapts to terminal window resize
- 📦 **Single binary** — No runtime dependencies, runs anywhere

## Installation

### From source

```bash
git clone https://github.com/karthiksedoc-enterprise-cloud-testing/presentshell.git
cd presentshell
go build -o presentshell .
```

### With `go install`

```bash
go install github.com/karthiksedoc-enterprise-cloud-testing/presentshell@latest
```

## Usage

```bash
presentshell <path-to-presentation.md>
```

### Example

```bash
presentshell examples/demo.md
```

## Writing Presentations

Create a markdown file with slides separated by `---`:

```markdown
# My Talk Title

Welcome to my presentation!

---

## Slide Two

- Point one
- Point two
- Point three

---

## Live Demo

Switch to the terminal and run:

\```bash
curl https://api.example.com/health
\```

---

# Thank You!
```

## Key Bindings

| Key | Action |
|-----|--------|
| `→` / `n` / `l` / `Space` | Next slide |
| `←` / `p` / `h` | Previous slide |
| `Tab` | Switch focus between slides and terminal |
| `q` / `Esc` | Quit (when slides focused) |
| `Ctrl+C` | Force quit |

> 💡 When the terminal pane is focused, all keys except `Tab` are forwarded to the shell.

## Architecture

```
┌─────────────────────────────────────────────┐
│              PresentShell                    │
├─────────────────────┬───────────────────────┤
│    Slide Pane       │    Terminal Pane       │
│                     │                       │
│  Glamour-rendered   │  VT100-emulated PTY   │
│  markdown content   │  (your default shell) │
│                     │                       │
├─────────────────────┴───────────────────────┤
│  Status: Slide 3/12 │ Tab:switch ←→:nav     │
└─────────────────────────────────────────────┘
```

## Tech Stack

| Component | Library |
|-----------|---------|
| TUI Framework | [Bubbletea](https://github.com/charmbracelet/bubbletea) |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) |
| Markdown | [Glamour](https://github.com/charmbracelet/glamour) |
| Terminal Emulation | [vt10x](https://github.com/hinshun/vt10x) |
| PTY | [creack/pty](https://github.com/creack/pty) |

## Project Structure

```
presentshell/
├── main.go              # Entry point, Bubbletea model & update loop
├── slides/
│   ├── slides.go        # Markdown file parser (splits by ---)
│   └── renderer.go      # Glamour-based markdown renderer
├── terminal/
│   └── terminal.go      # PTY + VT100 emulator management
├── ui/
│   ├── layout.go        # Split-pane layout logic
│   ├── statusbar.go     # Bottom status bar
│   └── styles.go        # Lipgloss theme/styles
└── examples/
    └── demo.md          # Sample presentation
```

## Tips

- Use a large terminal window (120+ columns) for best results
- The terminal pane runs your default `$SHELL`
- You can SSH into servers, run Docker containers, or execute any command in the live terminal
- Slides support full GitHub Flavored Markdown (tables, code blocks, lists, etc.)

## License

MIT
