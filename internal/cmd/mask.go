package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// MaskOptions controls how the mask command fetches and redacts secret values.
type MaskOptions struct {
	Address string
	Mount   string
	Token   string
	Version int
	Keys    []string // if non-empty, only these keys are masked
	Format  string
	DryRun  bool
	Mask    bool
	Output  io.Writer
}

// DefaultMaskOptions returns options populated from environment variables.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// MaskResult is the JSON-serialisable output of a Mask operation.
type MaskResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	DryRun  bool              `json:"dry_run"`
	Data    map[string]string `json:"data,omitempty"`
	Masked  []string          `json:"masked_keys"`
}

// Mask fetches a secret and returns it with the requested keys replaced by
// a redaction marker, optionally writing the sanitised secret back to Vault.
func Mask(path string, opts MaskOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Mount, opts.Token)
	if err != nil {
		return fmt.Errorf("mask: %w", err)
	}

	secret, err := client.GetSecretVersion(path, opts.Version)
	if err != nil {
		return fmt.Errorf("mask: fetch: %w", err)
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	masked := make([]string, 0)
	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		if len(keySet) == 0 || keySet[k] {
			result[k] = "***"
			masked = append(masked, k)
		} else {
			result[k] = v
		}
	}
	sort.Strings(masked)

	if !opts.DryRun {
		if err := client.WriteSecret(path, result); err != nil {
			return fmt.Errorf("mask: write: %w", err)
		}
	}

	record := MaskResult{
		Path:    path,
		Version: secret.Version,
		DryRun:  opts.DryRun,
		Masked:  masked,
	}
	if !opts.Mask {
		record.Data = result
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(record)
}
