package terminal

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/hinshun/vt10x"
)

// Terminal manages a pseudo-terminal session with VT100 emulation.
type Terminal struct {
	ptmx    *os.File
	cmd     *exec.Cmd
	vt      vt10x.Terminal
	mu      sync.Mutex
	rows    int
	cols    int
	running bool
}

// New creates and starts a new terminal with the given dimensions.
func New(rows, cols uint16) (*Terminal, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
	if err != nil {
		return nil, err
	}

	// Create virtual terminal emulator
	vt := vt10x.New(
		vt10x.WithWriter(ptmx),
		vt10x.WithSize(int(cols), int(rows)),
	)

	t := &Terminal{
		ptmx:    ptmx,
		cmd:     cmd,
		vt:      vt,
		rows:    int(rows),
		cols:    int(cols),
		running: true,
	}

	// Read PTY output and feed to VT emulator
	go t.readLoop()

	return t, nil
}

// readLoop continuously reads from the PTY and feeds data to the VT emulator.
func (t *Terminal) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := t.ptmx.Read(buf)
		if n > 0 {
			t.mu.Lock()
			t.vt.Write(buf[:n])
			t.mu.Unlock()
		}
		if err != nil {
			t.mu.Lock()
			t.running = false
			t.mu.Unlock()
			return
		}
	}
}

// Write sends input to the terminal.
func (t *Terminal) Write(data []byte) (int, error) {
	return t.ptmx.Write(data)
}

// Content returns the current terminal screen as a string.
func (t *Terminal) Content() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.vt.String()
}

// Resize updates the terminal dimensions.
func (t *Terminal) Resize(rows, cols uint16) error {
	t.mu.Lock()
	t.rows = int(rows)
	t.cols = int(cols)
	t.vt.Resize(int(cols), int(rows))
	t.mu.Unlock()

	return pty.Setsize(t.ptmx, &pty.Winsize{
		Rows: rows,
		Cols: cols,
	})
}

// IsRunning returns whether the terminal process is still active.
func (t *Terminal) IsRunning() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.running
}

// Close terminates the terminal session.
func (t *Terminal) Close() error {
	if t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
	return t.ptmx.Close()
}

// Fd returns the file descriptor of the PTY master.
func (t *Terminal) Fd() uintptr {
	return t.ptmx.Fd()
}

// File returns the PTY master file.
func (t *Terminal) File() io.Reader {
	return t.ptmx
}
