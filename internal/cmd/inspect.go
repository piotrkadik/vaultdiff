package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// InspectOptions configures the Inspect command.
type InspectOptions struct {
	Address string
	Mount   string
	Token   string
	Path    string
	Version int
	Format  string
	Mask    bool
	Output  io.Writer
}

// DefaultInspectOptions returns options populated from environment variables.
func DefaultInspectOptions() InspectOptions {
	return InspectOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// InspectResult holds the resolved secret data for a single version.
type InspectResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Data    map[string]string `json:"data,omitempty"`
}

// Inspect fetches a single secret version and prints its key/value pairs.
func Inspect(opts InspectOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Mount, opts.Token)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("inspect: %w", err)
	}

	result := InspectResult{
		Path:    opts.Path,
		Version: opts.Version,
	}

	if !opts.Mask {
		result.Data = secret.Data
	}

	if opts.Format == "json" {
		enc := json.NewEncoder(opts.Output)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Fprintf(opts.Output, "Path:    %s\n", result.Path)
	fmt.Fprintf(opts.Output, "Version: %d\n", result.Version)
	fmt.Fprintln(opts.Output, "Data:")

	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secret.Data[k]
		if opts.Mask {
			v = "***"
		}
		fmt.Fprintf(opts.Output, "  %s = %s\n", k, v)
	}
	return nil
}
