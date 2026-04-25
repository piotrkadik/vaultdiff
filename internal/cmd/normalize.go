package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// NormalizeOptions controls the behaviour of Normalize.
type NormalizeOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	DryRun  bool
	Mask    bool
	Format  string
	Output  io.Writer
}

// DefaultNormalizeOptions returns options populated from the environment.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// NormalizeResult holds the outcome of a Normalize run.
type NormalizeResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	DryRun  bool              `json:"dry_run"`
	Keys    int               `json:"keys"`
	Data    map[string]string `json:"data,omitempty"`
}

// Normalize trims whitespace and lowercases all keys in a secret, then
// writes the cleaned version back to Vault (unless DryRun is set).
func Normalize(opts NormalizeOptions) (*NormalizeResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("normalize: %w", err)
	}

	secret, err := client.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("normalize: fetch: %w", err)
	}

	normalized := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		nk := strings.TrimSpace(strings.ToLower(k))
		nv := strings.TrimSpace(v)
		normalized[nk] = nv
	}

	result := &NormalizeResult{
		Path:    opts.Path,
		Version: secret.Version,
		DryRun:  opts.DryRun,
		Keys:    len(normalized),
	}

	if !opts.DryRun {
		if err := client.WriteSecret(opts.Path, normalized); err != nil {
			return nil, fmt.Errorf("normalize: write: %w", err)
		}
		result.Data = normalized
	}

	if opts.Format == "json" {
		return result, json.NewEncoder(opts.Output).Encode(result)
	}

	keys := make([]string, 0, len(normalized))
	for k := range normalized {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if opts.DryRun {
		fmt.Fprintln(opts.Output, "[dry-run] normalize would rewrite the following keys:")
	}
	for _, k := range keys {
		v := normalized[k]
		if opts.Mask {
			v = "***"
		}
		fmt.Fprintf(opts.Output, "  %s = %s\n", k, v)
	}
	return result, nil
}
