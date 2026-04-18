package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// PinOptions configures the Pin command.
type PinOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	DryRun  bool
	Mask    bool
	Output  io.Writer
}

// DefaultPinOptions returns PinOptions populated from environment variables.
func DefaultPinOptions() PinOptions {
	mask := true
	if os.Getenv("VAULT_MASK") == "false" {
		mask = false
	}
	return PinOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   os.Getenv("VAULT_MOUNT"),
		Mask:    mask,
	}
}

// PinResult holds the outcome of a Pin operation.
type PinResult struct {
	Path    string         `json:"path"`
	Version int            `json:"pinned_version"`
	Data    map[string]any `json:"data,omitempty"`
	DryRun  bool           `json:"dry_run"`
}

// Pin records the current or specified version of a secret as the pinned reference.
func Pin(ctx context.Context, opts PinOptions) (*PinResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	if opts.Mount == "" {
		opts.Mount = "secret"
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("pin: %w", err)
	}

	secret, err := vault.GetSecretVersion(ctx, client, opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("pin: fetch version %d: %w", opts.Version, err)
	}

	result := &PinResult{
		Path:    opts.Path,
		Version: opts.Version,
		DryRun:  opts.DryRun,
	}

	if !opts.Mask {
		result.Data = secret.Data
	}

	if !opts.DryRun {
		enc := json.NewEncoder(opts.Output)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return nil, fmt.Errorf("pin: encode result: %w", err)
		}
	} else {
		fmt.Fprintf(opts.Output, "[dry-run] would pin %s at version %d\n", opts.Path, opts.Version)
	}

	return result, nil
}
