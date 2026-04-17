package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/vault"
)

// DefaultTrimOptions returns TrimOptions populated from environment variables.
func DefaultTrimOptions() TrimOptions {
	return TrimOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Keep:    5,
		DryRun:  false,
	}
}

// TrimOptions configures the Trim operation.
type TrimOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Keep    int
	DryRun  bool
	Output  io.Writer
}

// TrimResult holds the outcome of a Trim operation.
type TrimResult struct {
	Path    string `json:"path"`
	Removed []int  `json:"removed_versions"`
	DryRun  bool   `json:"dry_run"`
}

// validate checks that required TrimOptions fields are set and values are sane.
func (o TrimOptions) validate() error {
	if o.Path == "" {
		return fmt.Errorf("trim: path must not be empty")
	}
	if o.Keep < 1 {
		return fmt.Errorf("trim: keep must be at least 1, got %d", o.Keep)
	}
	return nil
}

// Trim deletes old secret versions beyond the Keep threshold.
func Trim(opts TrimOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	if err := opts.validate(); err != nil {
		return err
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("trim: %w", err)
	}

	metas, err := vault.ListVersions(client, opts.Path)
	if err != nil {
		return fmt.Errorf("trim: list versions: %w", err)
	}

	if len(metas) <= opts.Keep {
		return json.NewEncoder(opts.Output).Encode(TrimResult{
			Path:    opts.Path,
			Removed: []int{},
			DryRun:  opts.DryRun,
		})
	}

	// Sort ascending by version; metas assumed ordered.
	toRemove := metas[:len(metas)-opts.Keep]
	removed := make([]int, 0, len(toRemove))
	for _, m := range toRemove {
		if m.Destroyed {
			continue
		}
		removed = append(removed, m.Version)
		if !opts.DryRun {
			if err := vault.DestroyVersion(client, opts.Path, m.Version); err != nil {
				return fmt.Errorf("trim: destroy v%d: %w", m.Version, err)
			}
		}
	}

	return json.NewEncoder(opts.Output).Encode(TrimResult{
		Path:    opts.Path,
		Removed: removed,
		DryRun:  opts.DryRun,
	})
}
