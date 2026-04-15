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

// DefaultPromoteOptions returns PromoteOptions populated from environment/defaults.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		DryRun: false,
		Output: os.Stdout,
	}
}

// PromoteOptions controls the behaviour of Promote.
type PromoteOptions struct {
	// DryRun prints what would be written without performing the write.
	DryRun bool
	// Output is where the JSON confirmation record is written.
	Output io.Writer
}

// PromoteRecord is the JSON record emitted after a successful promotion.
type PromoteRecord struct {
	Path      string         `json:"path"`
	FromEnv   string         `json:"from_env"`
	ToEnv     string         `json:"to_env"`
	Version   int            `json:"version"`
	Data      map[string]any `json:"data"`
	DryRun    bool           `json:"dry_run"`
	Timestamp time.Time      `json:"timestamp"`
}

// Promote copies a secret version from one Vault path to another.
// When opts.DryRun is true the write is skipped and the record is still emitted.
func Promote(
	ctx context.Context,
	src vault.SecretGetter,
	dst vault.SecretWriter,
	path string,
	version int,
	fromEnv, toEnv string,
	opts PromoteOptions,
) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	secret, err := src.GetSecretVersion(ctx, path, version)
	if err != nil {
		return fmt.Errorf("promote: fetch %s@v%d: %w", path, version, err)
	}

	if !opts.DryRun {
		if err := dst.WriteSecret(ctx, path, secret.Data); err != nil {
			return fmt.Errorf("promote: write %s: %w", path, err)
		}
	}

	rec := PromoteRecord{
		Path:      path,
		FromEnv:   fromEnv,
		ToEnv:     toEnv,
		Version:   version,
		Data:      secret.Data,
		DryRun:    opts.DryRun,
		Timestamp: time.Now().UTC(),
	}

	enc := json.NewEncoder(opts.Output)
	enc.SetIndent("", "  ")
	return enc.Encode(rec)
}
