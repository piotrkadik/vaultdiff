package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/vault"
)

// DefaultRenameOptions returns RenameOptions populated from environment.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		DryRun:  false,
	}
}

// RenameOptions configures a Rename operation.
type RenameOptions struct {
	Address string
	Token   string
	Mount   string
	DryRun  bool
	Output  io.Writer
}

// RenameResult describes the outcome of a rename.
type RenameResult struct {
	Path    string `json:"path"`
	OldKey  string `json:"old_key"`
	NewKey  string `json:"new_key"`
	DryRun  bool   `json:"dry_run"`
}

// Rename copies a secret key to a new name and removes the old key.
func Rename(ctx context.Context, path, oldKey, newKey string, opts RenameOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	secret, err := client.GetSecretVersion(ctx, path, 0)
	if err != nil {
		return fmt.Errorf("rename: fetch %s: %w", path, err)
	}

	if _, ok := secret.Data[oldKey]; !ok {
		return fmt.Errorf("rename: key %q not found in %s", oldKey, path)
	}

	result := RenameResult{Path: path, OldKey: oldKey, NewKey: newKey, DryRun: opts.DryRun}

	if !opts.DryRun {
		updated := make(map[string]string, len(secret.Data))
		for k, v := range secret.Data {
			if k == oldKey {
				updated[newKey] = v
			} else {
				updated[k] = v
			}
		}
		if err := vault.WriteSecret(ctx, opts.Address, opts.Token, opts.Mount, path, updated); err != nil {
			return fmt.Errorf("rename: write: %w", err)
		}
	}

	return json.NewEncoder(out).Encode(result)
}
