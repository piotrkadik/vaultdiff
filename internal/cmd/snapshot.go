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

// Snapshot represents a point-in-time capture of a Vault secret version.
type Snapshot struct {
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	CapturedAt time.Time        `json:"captured_at"`
	Data      map[string]string `json:"data"`
}

// SnapshotOptions controls snapshot behaviour.
type SnapshotOptions struct {
	Mask   bool
	Output io.Writer
}

// DefaultSnapshotOptions returns sensible defaults.
func DefaultSnapshotOptions() SnapshotOptions {
	return SnapshotOptions{
		Mask:   true,
		Output: os.Stdout,
	}
}

// TakeSnapshot fetches the given secret version and writes it as JSON.
func TakeSnapshot(ctx context.Context, c *vault.Client, path string, version int, opts SnapshotOptions) (*Snapshot, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	sv, err := c.GetSecretVersion(ctx, path, version)
	if err != nil {
		return nil, fmt.Errorf("snapshot: fetch %s@%d: %w", path, version, err)
	}

	data := make(map[string]string, len(sv.Data))
	for k, v := range sv.Data {
		if opts.Mask {
			data[k] = "***"
		} else {
			data[k] = fmt.Sprintf("%v", v)
		}
	}

	snap := &Snapshot{
		Path:       path,
		Version:    version,
		CapturedAt: time.Now().UTC(),
		Data:       data,
	}

	enc := json.NewEncoder(opts.Output)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return nil, fmt.Errorf("snapshot: encode: %w", err)
	}

	return snap, nil
}
