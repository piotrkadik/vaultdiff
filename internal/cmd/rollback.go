package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

// RollbackOptions configures the Rollback command.
type RollbackOptions struct {
	// DryRun prints what would be written without applying changes.
	DryRun bool
	// Output is the writer used for status messages.
	Output io.Writer
}

// DefaultRollbackOptions returns sensible defaults for RollbackOptions.
func DefaultRollbackOptions() RollbackOptions {
	return RollbackOptions{
		DryRun: false,
		Output: os.Stdout,
	}
}

// RollbackRecord is the JSON-serialisable record emitted when DryRun is true.
type RollbackRecord struct {
	Path      string         `json:"path"`
	FromVersion int          `json:"from_version"`
	ToVersion   int          `json:"to_version"`
	Data        map[string]any `json:"data,omitempty"`
	Timestamp   time.Time    `json:"timestamp"`
	DryRun      bool         `json:"dry_run"`
}

// Rollback fetches targetVersion of path from the Vault client and, unless
// DryRun is set, writes the record to opts.Output as newline-delimited JSON.
func Rollback(ctx context.Context, vc *vault.Client, path string, currentVersion, targetVersion int, opts RollbackOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	secret, err := vc.GetSecretVersion(ctx, path, targetVersion)
	if err != nil {
		return fmt.Errorf("rollback: fetch version %d of %q: %w", targetVersion, path, err)
	}

	rec := RollbackRecord{
		Path:        path,
		FromVersion: currentVersion,
		ToVersion:   targetVersion,
		Data:        secret.Data,
		Timestamp:   time.Now().UTC(),
		DryRun:      opts.DryRun,
	}

	if opts.DryRun {
		fmt.Fprintf(opts.Output, "[dry-run] would roll back %q from v%d to v%d\n", path, currentVersion, targetVersion)
	}

	enc := json.NewEncoder(opts.Output)
	if err := enc.Encode(rec); err != nil {
		return fmt.Errorf("rollback: encode record: %w", err)
	}
	return nil
}
