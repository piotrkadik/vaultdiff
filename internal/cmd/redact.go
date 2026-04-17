package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// RedactOptions configures the Redact command.
type RedactOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Keys    []string // keys to redact; empty means all
	DryRun  bool
	Output  io.Writer
}

// DefaultRedactOptions returns options populated from environment variables.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   os.Getenv("VAULT_MOUNT"),
		DryRun:  false,
	}
}

// RedactResult holds the outcome of a redact operation.
type RedactResult struct {
	Path        string         `json:"path"`
	DryRun      bool           `json:"dry_run"`
	RedactedKeys []string      `json:"redacted_keys"`
	Data        map[string]any `json:"data,omitempty"`
}

// Redact replaces specified secret keys with a redacted placeholder and writes
// the new version back to Vault. If DryRun is true the write is skipped.
func Redact(opts RedactOptions) (*RedactResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("redact: create client: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("redact: fetch secret: %w", err)
	}

	redacted := make([]string, 0)
	for k, v := range secret.Data {
		if len(opts.Keys) == 0 || containsKey(opts.Keys, k) {
			_ = v
			secret.Data[k] = "REDACTED"
			redacted = append(redacted, k)
		}
	}

	result := &RedactResult{
		Path:         opts.Path,
		DryRun:       opts.DryRun,
		RedactedKeys: redacted,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, opts.Path, secret.Data); err != nil {
			return nil, fmt.Errorf("redact: write secret: %w", err)
		}
		result.Data = secret.Data
	}

	enc := json.NewEncoder(opts.Output)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
	return result, nil
}

func containsKey(keys []string, k string) bool {
	for _, key := range keys {
		if strings.EqualFold(key, k) {
			return true
		}
	}
	return false
}
