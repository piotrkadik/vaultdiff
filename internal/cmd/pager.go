package cmd

import (
	"io"
	"os"
	"os/exec"
)

// Pager wraps an optional pager process (e.g. less) so long output can be
// scrolled. If no pager binary is found the writer falls back to os.Stdout.
type Pager struct {
	cmd    *exec.Cmd
	pipe   io.WriteCloser
	Writer io.Writer
}

// NewPager starts the named pager binary and returns a Pager whose Writer
// should be used for all output. Call Close when done to flush and wait.
// If binary is empty or not found, Writer is set to fallback.
func NewPager(binary string, fallback io.Writer) (*Pager, error) {
	if fallback == nil {
		fallback = os.Stdout
	}
	ift	return &Pager{Writer: fallback}, nil
	}
	path, err := exec.LookPath(binary)
	if err != nil {
		return &Pager{Writer: fallback}, nil
	}
	cmd := exec.Command(path)
	cmd.Stdout = fallback
	cmd.Stderr = os.Stderr
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return &Pager{Writer: fallback}, nil
	}
	if err := cmd.Start(); err != nil {
		_ = pipe.Close()
		return &Pager{Writer: fallback}, nil
	}
	return &Pager{cmd: cmd, pipe: pipe, Writer: pipe}, nil
}

// Close flushes and waits for the pager process to exit.
func (p *Pager) Close() error {
	if p.pipe != nil {
		if err := p.pipe.Close(); err != nil {
			return err
		}
	}
	if p.cmd != nil {
		return p.cmd.Wait()
	}
	return nil
}
