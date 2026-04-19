package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/someone/vaultdiff/internal/vault"
)

// DefaultRevertOptions returns RevertOptions populated from environment variables.
func DefaultRevertOptions() RevertOptions {
	return RevertOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		DryRun:  false,
	}
}

// RevertOptions configures the Revert command.
type RevertOptions struct {
	Address    string
	Token      string
	Mount      string
	Path       string
	Version    int
	Mask       bool
	DryRun     bool
	Output     io.Writer
}

// RevertResult holds the outcome of a revert operation.
type RevertResult struct {
	Path    string         `json:"path"`
	Version int            `json:"version"`
	DryRun  bool           `json:"dry_run"`
	Data    map[string]any `json:"data,omitempty"`
}

// Revert fetches a specific secret version and re-writes it as the latest version.
func Revert(ctx context.Context, opts RevertOptions) (*RevertResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	cl, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("revert: new client: %w", err)
	}

	secret, err := vault.GetSecretVersion(ctx, cl, opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("revert: fetch version %d: %w", opts.Version, err)
	}

	result := &RevertResult{
		Path:    opts.Path,
		Version: opts.Version,
		DryRun:  opts.DryRun,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(ctx, cl, opts.Path, secret); err != nil {
			return nil, fmt.Errorf("revert: write: %w", err)
		}
		if !opts.Mask {
			result.Data = secret
		}
	}

	if err := json.NewEncoder(opts.Output).Encode(result); err != nil {
		return nil, fmt.Errorf("revert: encode result: %w", err)
	}

	return result, nil
}
