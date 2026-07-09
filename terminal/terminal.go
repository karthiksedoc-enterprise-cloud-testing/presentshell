package terminal

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// Terminal manages a pseudo-terminal session.
type Terminal struct {
	ptmx    *os.File
	cmd     *exec.Cmd
	mu      sync.Mutex
	rows    uint16
	cols    uint16
	output  []byte
	maxBuf  int
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

	t := &Terminal{
		ptmx:    ptmx,
		cmd:     cmd,
		rows:    rows,
		cols:    cols,
		output:  make([]byte, 0, 4096),
		maxBuf:  32 * 1024, // 32KB buffer
		running: true,
	}

	// Read PTY output in background
	go t.readLoop()

	return t, nil
}

// readLoop continuously reads from the PTY.
func (t *Terminal) readLoop() {
	buf := make([]byte, 1024)
	for {
		n, err := t.ptmx.Read(buf)
		if n > 0 {
			t.mu.Lock()
			t.output = append(t.output, buf[:n]...)
			// Trim buffer if too large
			if len(t.output) > t.maxBuf {
				t.output = t.output[len(t.output)-t.maxBuf:]
			}
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

// Read returns the current terminal output buffer.
func (t *Terminal) Read() []byte {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]byte, len(t.output))
	copy(out, t.output)
	return out
}

// Resize updates the terminal dimensions.
func (t *Terminal) Resize(rows, cols uint16) error {
	t.mu.Lock()
	t.rows = rows
	t.cols = cols
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

// Fd returns the file descriptor of the PTY master for use with readers.
func (t *Terminal) Fd() uintptr {
	return t.ptmx.Fd()
}

// File returns the PTY master file for use with io.Reader.
func (t *Terminal) File() io.Reader {
	return t.ptmx
}
