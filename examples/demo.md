# PresentShell

### A Terminal-Based Presentation Tool

> Present slides and run live demos — all from your terminal.

**Built with Go + Charm ecosystem**

Press `→` or `n` to begin...

---

## 🎯 The Problem

Traditional presentation tools:

- ❌ Can't run live code demos
- ❌ Require switching between apps
- ❌ Break your flow as a developer
- ❌ Don't work well over SSH

**PresentShell** solves this with a split-pane layout:
slides on the left, live terminal on the right.

---

## ⌨️ Navigation

| Key | Action |
|-----|--------|
| `Tab` | Switch focus (slides ↔ terminal) |
| `→` / `n` | Next slide |
| `←` / `p` | Previous slide |
| `q` | Quit (when slides focused) |
| `Ctrl+C` | Force quit |

> 💡 When the terminal is focused, all keys go to the shell!

---

## 📝 Writing Slides

Slides are plain **Markdown** files separated by `---`:

```markdown
# First Slide

Content here...

---

## Second Slide

More content...
```

That's it! No special syntax needed.

---

## 🖥 Live Terminal Demo

Switch to the terminal with `Tab` and try:

```bash
# Check system info
uname -a

# List files
ls -la

# Run a script
echo "Hello from PresentShell! 🚀"
```

The right pane is a **real PTY** — run anything you want.

---

## 🐳 Docker Demo

You can demo containerized apps live:

```bash
# Pull and run a container
docker run --rm -it alpine sh

# Or show docker compose
docker compose up -d
docker compose ps
```

Your audience sees exactly what you see!

---

## 🔗 SSH Into Servers

Perfect for ops/infra presentations:

```bash
# Connect to a remote server
ssh user@production-server

# Show live metrics
htop

# Check logs
tail -f /var/log/app.log
```

No more pre-recorded terminal GIFs!

---

## 🐍 Python Example

```python
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/items/{item_id}")
async def read_item(item_id: int):
    return {"item_id": item_id}
```

Try running it in the terminal →

---

## 🦀 Rust Example

```rust
use actix_web::{web, App, HttpServer};

async fn hello() -> &'static str {
    "Hello from Rust!"
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new().route("/", web::get().to(hello))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
```

---

## 🏗 Architecture

```
┌─────────────────────────────────────────┐
│            PresentShell                  │
├───────────────────┬─────────────────────┤
│   Slide Pane      │   Terminal Pane      │
│                   │                      │
│  ┌─────────────┐  │  ┌───────────────┐  │
│  │  Glamour    │  │  │   PTY Shell   │  │
│  │  Rendered   │  │  │   (creack/pty)│  │
│  │  Markdown   │  │  │               │  │
│  └─────────────┘  │  └───────────────┘  │
│                   │                      │
├───────────────────┴─────────────────────┤
│  Status Bar: Slide 10/12 │ Tab:switch   │
└─────────────────────────────────────────┘
```

---

## 📦 Tech Stack

| Component | Library |
|-----------|---------|
| TUI Framework | `bubbletea` (Elm architecture) |
| Styling | `lipgloss` (CSS-like) |
| Markdown | `glamour` (terminal renderer) |
| Terminal | `creack/pty` (pseudo-terminal) |
| Language | **Go** (single binary) |

All from the [Charm](https://charm.sh) ecosystem 🪄

---

## 🚀 Getting Started

```bash
# Install
go install github.com/karthik/presentshell@latest

# Run with a presentation file
presentshell my-talk.md

# Or use the demo
presentshell examples/demo.md
```

Write your slides in any text editor,
present from any terminal. Simple.

---

# 🙏 Thank You!

**PresentShell** — Present like a developer.

- 📖 Slides are just Markdown
- 🖥 Live terminal built-in
- 🎨 Beautiful rendering
- 📦 Single binary, no dependencies

> "The best demo is a live demo."

*Questions? Switch to the terminal and let's explore!*
