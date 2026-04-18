package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// DefaultUnlockOptions returns UnlockOptions populated from environment variables.
func DefaultUnlockOptions() UnlockOptions {
	return UnlockOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		DryRun:  false,
	}
}

// UnlockOptions configures the Unlock operation.
type UnlockOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	DryRun  bool
	Output  io.Writer
}

// UnlockResult holds the outcome of an unlock operation.
type UnlockResult struct {
	Path    string         `json:"path"`
	Version int            `json:"version"`
	DryRun  bool           `json:"dry_run"`
	Data    map[string]any `json:"data,omitempty"`
}

// Unlock undeletes (unlocks) a previously soft-deleted secret version in Vault.
func Unlock(opts UnlockOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("unlock: create client: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("unlock: fetch version %d: %w", opts.Version, err)
	}

	result := UnlockResult{
		Path:    opts.Path,
		Version: opts.Version,
		DryRun:  opts.DryRun,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, opts.Path, secret.Data); err != nil {
			return fmt.Errorf("unlock: write secret: %w", err)
		}
		result.Data = secret.Data
	}

	return json.NewEncoder(out).Encode(result)
}
