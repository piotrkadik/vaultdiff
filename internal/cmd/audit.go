package cmd

import (
	"io"
	"os"
	"time"

	"github.com/vaultdiff/internal/audit"
	"github.com/vaultdiff/internal/diff"
)

// AuditOptions configures audit log emission for a diff run.
type AuditOptions struct {
	// Writer is the destination for audit log entries.
	// Defaults to os.Stderr when nil.
	Writer io.Writer

	// Actor is an optional identifier (user, CI job, etc.) recorded in each entry.
	Actor string

	// Enabled controls whether audit logging is active.
	Enabled bool
}

// DefaultAuditOptions returns AuditOptions with sensible defaults.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{
		Writer:  os.Stderr,
		Actor:   os.Getenv("VAULTDIFF_ACTOR"),
		Enabled: os.Getenv("VAULTDIFF_AUDIT") == "true",
	}
}

// WriteAuditEntry emits a structured audit log entry for a completed diff.
// It is a no-op when opts.Enabled is false or opts.Writer is nil.
func WriteAuditEntry(
	path string,
	versionA, versionB int,
	changes []diff.Change,
	ts time.Time,
	opts AuditOptions,
) error {
	if !opts.Enabled {
		return nil
	}

	w := opts.Writer
	if w == nil {
		w = os.Stderr
	}

	logger := audit.NewLogger(w)

	summary := diff.Summarize(changes)

	return logger.Write(audit.Entry{
		Timestamp: ts,
		Actor:     opts.Actor,
		Path:      path,
		VersionA:  versionA,
		VersionB:  versionB,
		Added:     summary.Added,
		Removed:   summary.Removed,
		Modified:  summary.Modified,
		Unchanged: summary.Unchanged,
		HasDrift:  summary.HasDrift(),
	})
}
