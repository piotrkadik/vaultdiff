package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultImportOptions returns ImportOptions populated from environment variables.
func DefaultImportOptions() ImportOptions {
	return ImportOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		DryRun:  false,
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// ImportOptions configures the Import command.
type ImportOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	DryRun  bool
	Mask    bool
	Output  io.Writer
}

// ImportResult holds the outcome of an import operation.
type ImportResult struct {
	Path    string         `json:"path"`
	Keys    int            `json:"keys"`
	DryRun  bool           `json:"dry_run"`
	Data    map[string]any `json:"data,omitempty"`
}

// Import reads a JSON file from r and writes the key/value pairs to Vault.
func Import(r io.Reader, opts ImportOptions) (ImportResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	var data map[string]string
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return ImportResult{}, fmt.Errorf("import: decode: %w", err)
	}

	result := ImportResult{
		Path:   opts.Path,
		Keys:   len(data),
		DryRun: opts.DryRun,
	}

	if !opts.DryRun {
		client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
		if err != nil {
			return ImportResult{}, fmt.Errorf("import: client: %w", err)
		}
		if err := vault.WriteSecret(client, opts.Mount, opts.Path, data); err != nil {
			return ImportResult{}, fmt.Errorf("import: write: %w", err)
		}
		any := make(map[string]any, len(data))
		for k, v := range data {
			if opts.Mask {
				any[k] = "***"
			} else {
				any[k] = v
			}
		}
		result.Data = any
	}

	return result, nil
}
