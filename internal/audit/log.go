// Package audit provides structured audit logging for vaultdiff operations.
package audit

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record capturing a diff operation.
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	Environment string    `json:"environment"`
	Path        string    `json:"path"`
	VersionA    int       `json:"version_a"`
	VersionB    int       `json:"version_b"`
	Added       int       `json:"added"`
	Removed     int       `json:"removed"`
	Modified    int       `json:"modified"`
	HasDrift    bool      `json:"has_drift"`
	User        string    `json:"user,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// NewLogger returns a Logger that writes to w.
// Pass os.Stdout or an *os.File as appropriate.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{w: w}
}

// Write serialises entry as a single JSON line to the underlying writer.
// It stamps the entry with the current UTC time if Timestamp is zero.
func (l *Logger) Write(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = l.w.Write(data)
	return err
}
