package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultPruneOptions returns PruneOptions populated from environment variables.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Keep:    5,
		DryRun:  false,
	}
}

// PruneOptions configures the Prune command.
type PruneOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Keep    int
	DryRun  bool
	Output  io.Writer
}

// PruneResult holds the outcome of a prune operation.
type PruneResult struct {
	Path     string `json:"path"`
	Destroyed []int  `json:"destroyed"`
	Kept      int    `json:"kept"`
	DryRun   bool   `json:"dry_run"`
}

// Prune removes old secret versions beyond the keep threshold.
func Prune(ctx context.Context, opts PruneOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("prune: new client: %w", err)
	}

	metas, err := vault.ListVersions(ctx, client, opts.Path)
	if err != nil {
		return fmt.Errorf("prune: list versions: %w", err)
	}

	var toDestroy []int
	active := 0
	for i := len(metas) - 1; i >= 0; i-- {
		m := metas[i]
		if m.Destroyed {
			continue
		}
		active++
		if active > opts.Keep {
			toDestroy = append(toDestroy, m.Version)
		}
	}

	if !opts.DryRun {
		for _, v := range toDestroy {
			if err := vault.DestroyVersion(ctx, client, opts.Path, v); err != nil {
				return fmt.Errorf("prune: destroy version %d: %w", v, err)
			}
		}
	}

	result := PruneResult{
		Path:      opts.Path,
		Destroyed: toDestroy,
		Kept:      opts.Keep,
		DryRun:    opts.DryRun,
	}
	return json.NewEncoder(out).Encode(result)
}
