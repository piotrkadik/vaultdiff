package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// ScrubOptions controls the behaviour of the Scrub command.
type ScrubOptions struct {
	Address string
	Mount   string
	Token   string
	Path    string
	Version int
	Keys    []string // keys whose values should be zeroed out
	DryRun  bool
	Format  string
	Mask    bool
	Output  io.Writer
}

// DefaultScrubOptions returns options populated from environment variables.
func DefaultScrubOptions() ScrubOptions {
	mask := true
	if strings.EqualFold(os.Getenv("VAULTDIFF_MASK"), "false") {
		mask = false
	}
	return ScrubOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    mask,
		DryRun:  false,
	}
}

// ScrubResult is the JSON-serialisable result of a Scrub operation.
type ScrubResult struct {
	Path        string   `json:"path"`
	Version     int      `json:"version"`
	Scrubbed    []string `json:"scrubbed"`
	DryRun      bool     `json:"dry_run"`
}

// Scrub zeroes out the specified keys in a Vault secret, writing a new version.
func Scrub(c vault.Client, opts ScrubOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	secret, err := vault.GetSecretVersion(c, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("scrub: fetch %s@%d: %w", opts.Path, opts.Version, err)
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	scrubbed := make([]string, 0)
	for _, k := range opts.Keys {
		if _, exists := secret[k]; exists {
			scrubbed = append(scrubbed, k)
			if !opts.DryRun {
				secret[k] = ""
			}
		}
	}

	if !opts.DryRun && len(scrubbed) > 0 {
		if err := vault.WriteSecret(c, opts.Path, secret); err != nil {
			return fmt.Errorf("scrub: write %s: %w", opts.Path, err)
		}
	}

	result := ScrubResult{
		Path:     opts.Path,
		Version:  opts.Version,
		Scrubbed: scrubbed,
		DryRun:   opts.DryRun,
	}

	if opts.DryRun {
		fmt.Fprintln(out, "[dry-run] no changes written")
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
