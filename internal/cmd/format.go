package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// DefaultFormatOptions returns FormatOptions populated from environment variables.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// FormatOptions controls the behaviour of Format.
type FormatOptions struct {
	Address string
	Mount   string
	Token   string
	Path    string
	Version int
	Mask    bool
	Format  string
	Output  io.Writer
}

// FormatResult holds the rendered output of a Format call.
type FormatResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Keys    []string          `json:"keys"`
	Data    map[string]string `json:"data,omitempty"`
}

// Format fetches a secret version and prints its key/value pairs in the
// requested output format, optionally masking values.
func Format(vc VaultClient, opts FormatOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	data, err := vc.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("format: fetch %s@%d: %w", opts.Path, opts.Version, err)
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := FormatResult{
		Path:    opts.Path,
		Version: opts.Version,
		Keys:    keys,
	}

	if !opts.Mask {
		result.Data = data
	}

	if strings.EqualFold(opts.Format, "json") {
		return json.NewEncoder(opts.Output).Encode(result)
	}

	fmt.Fprintf(opts.Output, "# %s (v%d)\n", opts.Path, opts.Version)
	for _, k := range keys {
		v := data[k]
		if opts.Mask {
			v = "***"
		}
		fmt.Fprintf(opts.Output, "%s=%s\n", k, v)
	}
	return nil
}
