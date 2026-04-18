package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// KeysOptions configures the Keys command.
type KeysOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Format  string
	Output  io.Writer
}

// DefaultKeysOptions returns options populated from environment variables.
func DefaultKeysOptions() KeysOptions {
	return KeysOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  formatFromEnv(),
	}
}

// KeysResult holds the list of keys returned from a secret version.
type KeysResult struct {
	Path    string   `json:"path"`
	Version int      `json:"version"`
	Keys    []string `json:"keys"`
	Count   int      `json:"count"`
}

// Keys fetches a secret version and prints its top-level keys.
func Keys(opts KeysOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("keys: %w", err)
	}

	data, err := vault.GetSecretVersion(client, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("keys: %w", err)
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := KeysResult{
		Path:    opts.Path,
		Version: opts.Version,
		Keys:    keys,
		Count:   len(keys),
	}

	if opts.Format == "json" {
		return json.NewEncoder(opts.Output).Encode(result)
	}

	for _, k := range keys {
		fmt.Fprintln(opts.Output, k)
	}
	return nil
}
