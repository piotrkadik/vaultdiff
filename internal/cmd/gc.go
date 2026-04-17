package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/vault"
)

// DefaultGCOptions returns GCOptions populated from environment variables.
func DefaultGCOptions() GCOptions {
	return GCOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Keep:    5,
		DryRun:  false,
	}
}

// GCOptions controls the behaviour of the GC command.
type GCOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Keep    int
	DryRun  bool
	Output  io.Writer
}

// GCResult summarises what was (or would be) deleted.
type GCResult struct {
	Path    string `json:"path"`
	Deleted []int  `json:"deleted"`
	Kept    []int  `json:"kept"`
	DryRun  bool   `json:"dry_run"`
}

// GC destroys old secret versions beyond the keep threshold.
func GC(opts GCOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("gc: %w", err)
	}

	metas, err := vault.ListVersions(client, opts.Path)
	if err != nil {
		return fmt.Errorf("gc: list versions: %w", err)
	}

	var versions []int
	for _, m := range metas {
		if !m.Destroyed {
			versions = append(versions, m.Version)
		}
	}

	var kept, deleted []int
	for i, v := range versions {
		if i < len(versions)-opts.Keep {
			deleted = append(deleted, v)
		} else {
			kept = append(kept, v)
		}
	}

	result := GCResult{Path: opts.Path, Deleted: deleted, Kept: kept, DryRun: opts.DryRun}

	if !opts.DryRun {
		for _, v := range deleted {
			if err := vault.DestroyVersion(client, opts.Path, v); err != nil {
				return fmt.Errorf("gc: destroy v%d: %w", v, err)
			}
		}
	}

	return json.NewEncoder(opts.Output).Encode(result)
}
