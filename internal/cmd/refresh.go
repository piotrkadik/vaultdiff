package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

// RefreshOptions configures the Refresh command.
type RefreshOptions struct {
	Address string
	Mount   string
	Token   string
	Path    string
	Mask    bool
	DryRun  bool
	Format  string
	Output  io.Writer
}

// DefaultRefreshOptions returns options populated from environment variables.
func DefaultRefreshOptions() RefreshOptions {
	mask := true
	if os.Getenv("VAULTDIFF_MASK") == "false" {
		mask = false
	}
	return RefreshOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   os.Getenv("VAULT_MOUNT"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    mask,
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// RefreshResult holds the outcome of a Refresh operation.
type RefreshResult struct {
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	Refreshed bool              `json:"refreshed"`
	DryRun    bool              `json:"dry_run"`
	Data      map[string]string `json:"data,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// Refresh re-writes the current version of a secret to bump its metadata
// timestamp without changing any values. In dry-run mode it only reads.
func Refresh(client vault.SecretGetter, writer vault.SecretWriter, opts RefreshOptions) (RefreshResult, error) {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	sv, err := client.GetSecretVersion(opts.Path, 0)
	if err != nil {
		return RefreshResult{}, fmt.Errorf("refresh: fetch %s: %w", opts.Path, err)
	}

	result := RefreshResult{
		Path:      opts.Path,
		Version:   sv.Version,
		Refreshed: !opts.DryRun,
		DryRun:    opts.DryRun,
		Timestamp: time.Now().UTC(),
	}

	if !opts.DryRun {
		if err := writer.WriteSecret(opts.Path, sv.Data); err != nil {
			return RefreshResult{}, fmt.Errorf("refresh: write %s: %w", opts.Path, err)
		}
		if !opts.Mask {
			result.Data = sv.Data
		}
	}

	if opts.Format == "json" {
		return result, json.NewEncoder(out).Encode(result)
	}

	status := "refreshed"
	if opts.DryRun {
		status = "dry-run (skipped write)"
	}
	fmt.Fprintf(out, "path=%s version=%d status=%s\n", result.Path, result.Version, status)
	return result, nil
}
