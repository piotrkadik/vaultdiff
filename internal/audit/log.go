package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record for a diff operation.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor,omitempty"`
	Path      string    `json:"path"`
	VersionA  int       `json:"version_a"`
	VersionB  int       `json:"version_b"`
	Added     int       `json:"added"`
	Removed   int       `json:"removed"`
	Modified  int       `json:"modified"`
	Unchanged int       `json:"unchanged"`
	HasDrift  bool      `json:"has_drift"`
}

// Logger writes newline-delimited JSON audit entries to an io.Writer.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger that writes to w.
// When w is nil, os.Stderr is used.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w}
}

// Write serialises e as a JSON line. If e.Timestamp is zero it is set to
// time.Now() before encoding.
func (l *Logger) Write(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	b = append(b, '\n')
	_, err = l.w.Write(b)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}
