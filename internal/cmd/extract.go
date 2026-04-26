package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultExtractOptions returns ExtractOptions populated from environment variables.
func DefaultExtractOptions() ExtractOptions {
	return ExtractOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// ExtractOptions configures the Extract command.
type ExtractOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Keys    []string // if non-empty, only these keys are extracted
	Format  string
	Mask    bool
	Output  io.Writer
}

// ExtractResult holds the subset of key/value pairs pulled from a secret.
type ExtractResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Data    map[string]string `json:"data,omitempty"`
}

// Extract fetches a secret version and returns only the requested keys.
func Extract(opts ExtractOptions) (*ExtractResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("extract: %w", err)
	}

	filtered := make(map[string]string)
	if len(opts.Keys) == 0 {
		for k, v := range secret.Data {
			filtered[k] = maskIfNeeded(v, opts.Mask)
		}
	} else {
		for _, k := range opts.Keys {
			v, ok := secret.Data[k]
			if !ok {
				return nil, fmt.Errorf("extract: key %q not found in %s@%d", k, opts.Path, secret.Version)
			}
			filtered[k] = maskIfNeeded(v, opts.Mask)
		}
	}

	result := &ExtractResult{
		Path:    opts.Path,
		Version: secret.Version,
		Data:    filtered,
	}

	if err := writeExtractResult(opts.Output, result, opts.Format); err != nil {
		return nil, err
	}
	return result, nil
}

func maskIfNeeded(v string, mask bool) string {
	if mask {
		return "***"
	}
	return v
}

func writeExtractResult(w io.Writer, r *ExtractResult, format string) error {
	if format == "json" {
		return json.NewEncoder(w).Encode(r)
	}
	keys := make([]string, 0, len(r.Data))
	for k := range r.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, r.Data[k])
	}
	return nil
}
