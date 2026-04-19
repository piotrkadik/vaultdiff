package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// FetchOptions configures the Fetch command.
type FetchOptions struct {
	Address string
	Mount   string
	Token   string
	Version int
	Mask    bool
	Format  string
	Output  io.Writer
}

// DefaultFetchOptions returns options populated from environment variables.
func DefaultFetchOptions() FetchOptions {
	mask := true
	if os.Getenv("VAULT_MASK") == "false" {
		mask = false
	}
	return FetchOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   os.Getenv("VAULT_MOUNT"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    mask,
		Format:  os.Getenv("VAULT_FORMAT"),
	}
}

// Fetch retrieves a single secret version and writes it to the output.
func Fetch(path string, opts FetchOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Mount, opts.Token)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	data, err := vault.GetSecretVersion(client, path, opts.Version)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	if opts.Mask {
		for k := range data {
			data[k] = "***"
		}
	}

	if opts.Format == "json" {
		return json.NewEncoder(opts.Output).Encode(map[string]any{
			"path":    path,
			"version": opts.Version,
			"data":    data,
		})
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(opts.Output, "%s=%s\n", k, data[k])
	}
	return nil
}
