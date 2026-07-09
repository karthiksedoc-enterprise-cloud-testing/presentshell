# Welcome to PresentShell

A terminal-based presentation tool with live demos!

- **Left pane**: Your slides in beautiful markdown
- **Right pane**: A live terminal for demonstrations

Press `→` or `n` to go to the next slide.

---

## Navigation

| Key | Action |
|-----|--------|
| `Tab` | Switch focus between slides and terminal |
| `→` / `n` | Next slide |
| `←` / `p` | Previous slide |
| `q` | Quit (when slides are focused) |

---

## Live Demo

Switch to the terminal pane with `Tab` and try:

```bash
echo "Hello from PresentShell!"
ls -la
```

The terminal is a **real shell** — you can SSH into servers,
run scripts, or demonstrate anything!

---

## Code Highlighting

PresentShell renders code with syntax highlighting:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, PresentShell!")
}
```

---

## Features

- ✅ Markdown rendering with Glamour
- ✅ Syntax highlighting for code blocks
- ✅ Split-pane layout
- ✅ Embedded terminal (PTY)
- ✅ Keyboard navigation
- ✅ Responsive to window resize

---

# Thank You!

Built with ❤️ using:
- **Bubbletea** — TUI framework
- **Lipgloss** — Terminal styling
- **Glamour** — Markdown rendering
- **creack/pty** — Pseudo-terminal

> "The best demo is a live demo."
