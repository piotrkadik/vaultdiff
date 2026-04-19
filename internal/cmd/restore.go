package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// RestoreOptions configures the Restore command.
type RestoreOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	DryRun  bool
	Mask    bool
	Output  io.Writer
}

// DefaultRestoreOptions returns options populated from environment variables.
func DefaultRestoreOptions() RestoreOptions {
	return RestoreOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// RestoreResult holds the outcome of a restore operation.
type RestoreResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	DryRun  bool              `json:"dry_run"`
	Data    map[string]string `json:"data,omitempty"`
}

// Restore re-writes a specific historical version of a secret as the latest version.
func Restore(ctx context.Context, opts RestoreOptions) (*RestoreResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("restore: new client: %w", err)
	}

	secret, err := vault.GetSecretVersion(ctx, client, opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("restore: fetch version %d: %w", opts.Version, err)
	}

	result := &RestoreResult{
		Path:    opts.Path,
		Version: opts.Version,
		DryRun:  opts.DryRun,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(ctx, client, opts.Path, secret.Data); err != nil {
			return nil, fmt.Errorf("restore: write: %w", err)
		}
		result.Data = secret.Data
	}

	if err := json.NewEncoder(opts.Output).Encode(result); err != nil {
		return nil, fmt.Errorf("restore: encode result: %w", err)
	}

	return result, nil
}
