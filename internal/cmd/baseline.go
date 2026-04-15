package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/vault"
)

// BaselineOptions controls how a baseline snapshot is captured and stored.
type BaselineOptions struct {
	Mount   string
	Output  io.Writer
	Mask    bool
	Version int
}

// DefaultBaselineOptions returns sensible defaults for baseline capture.
func DefaultBaselineOptions() BaselineOptions {
	return BaselineOptions{
		Mount:  "secret",
		Output: os.Stdout,
		Mask:   true,
	}
}

// BaselineRecord holds a captured baseline for a single secret path.
type BaselineRecord struct {
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	CapturedAt time.Time        `json:"captured_at"`
	Data      map[string]string `json:"data"`
}

// CaptureBaseline fetches the specified secret version from Vault and writes
// a baseline JSON record to opts.Output. It is intended to be used before a
// deployment so that post-deployment drift can be detected via Compare.
func CaptureBaseline(ctx context.Context, client *vault.Client, path string, opts BaselineOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	sv, err := client.GetSecretVersion(ctx, path, opts.Version)
	if err != nil {
		return fmt.Errorf("baseline: fetch %q v%d: %w", path, opts.Version, err)
	}

	data := make(map[string]string, len(sv.Data))
	for k, v := range sv.Data {
		val := fmt.Sprintf("%v", v)
		if opts.Mask {
			val = diff.MaskValue(val)
		}
		data[k] = val
	}

	rec := BaselineRecord{
		Path:       path,
		Version:    sv.Version,
		CapturedAt: time.Now().UTC(),
		Data:       data,
	}

	enc := json.NewEncoder(opts.Output)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rec); err != nil {
		return fmt.Errorf("baseline: encode record: %w", err)
	}
	return nil
}
