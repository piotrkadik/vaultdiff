package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/vault"
)

// DefaultRotateOptions returns RotateOptions populated from environment variables.
func DefaultRotateOptions() RotateOptions {
	return RotateOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		DryRun:  false,
	}
}

// RotateOptions configures the Rotate command.
type RotateOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Mask    bool
	DryRun  bool
	Output  io.Writer
}

// RotateResult holds the outcome of a rotate operation.
type RotateResult struct {
	Path        string         `json:"path"`
	DryRun      bool           `json:"dry_run"`
	PrevVersion int            `json:"prev_version"`
	NewVersion  int            `json:"new_version,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
}

// Rotate re-writes the latest version of a secret to create a new version,
// effectively rotating it in-place without changing values.
func Rotate(opts RotateOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, 0)
	if err != nil {
		return fmt.Errorf("rotate: fetch: %w", err)
	}

	result := RotateResult{
		Path:        opts.Path,
		DryRun:      opts.DryRun,
		PrevVersion: secret.Version,
	}

	if !opts.DryRun {
		newVersion, err := vault.WriteSecret(client, opts.Path, secret.Data)
		if err != nil {
			return fmt.Errorf("rotate: write: %w", err)
		}
		result.NewVersion = newVersion
		if !opts.Mask {
			result.Data = secret.Data
		}
	}

	return json.NewEncoder(out).Encode(result)
}
