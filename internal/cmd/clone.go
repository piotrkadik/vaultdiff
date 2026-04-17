package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/vault"
)

// DefaultCloneOptions returns CloneOptions populated from environment.
func DefaultCloneOptions() CloneOptions {
	return CloneOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		DryRun:  false,
	}
}

// CloneOptions configures the Clone operation.
type CloneOptions struct {
	Address  string
	Token    string
	Mount    string
	SrcPath  string
	DstPath  string
	Mask     bool
	DryRun   bool
	Output   io.Writer
}

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SrcPath string         `json:"src_path"`
	DstPath string         `json:"dst_path"`
	Keys    int            `json:"keys_cloned"`
	DryRun  bool           `json:"dry_run"`
	Data    map[string]any `json:"data,omitempty"`
}

// Clone reads secrets from SrcPath and writes them to DstPath.
func Clone(ctx context.Context, opts CloneOptions) (*CloneResult, error) {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("clone: new client: %w", err)
	}

	secret, err := vault.GetSecretVersion(ctx, client, opts.SrcPath, 0)
	if err != nil {
		return nil, fmt.Errorf("clone: fetch src: %w", err)
	}

	result := &CloneResult{
		SrcPath: opts.SrcPath,
		DstPath: opts.DstPath,
		Keys:    len(secret.Data),
		DryRun:  opts.DryRun,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(ctx, client, opts.DstPath, secret.Data); err != nil {
			return nil, fmt.Errorf("clone: write dst: %w", err)
		}
		result.Data = secret.Data
	}

	if err := json.NewEncoder(out).Encode(result); err != nil {
		return nil, fmt.Errorf("clone: encode result: %w", err)
	}
	return result, nil
}
