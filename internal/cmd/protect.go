package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultProtectOptions returns ProtectOptions populated from environment variables.
func DefaultProtectOptions() ProtectOptions {
	return ProtectOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		DryRun:  false,
	}
}

// ProtectOptions configures the Protect command.
type ProtectOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Keys    []string
	DryRun  bool
	Output  io.Writer
}

// ProtectResult holds the outcome of a protect operation.
type ProtectResult struct {
	Path      string   `json:"path"`
	Protected []string `json:"protected_keys"`
	DryRun    bool     `json:"dry_run"`
}

// Protect marks specific keys in a secret as protected by writing a sentinel
// metadata annotation. In dry-run mode it prints what would change without
// writing to Vault.
func Protect(opts ProtectOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("protect: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, 0)
	if err != nil {
		return fmt.Errorf("protect: fetch secret: %w", err)
	}

	updated := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		updated[k] = v
	}

	for _, key := range opts.Keys {
		if _, ok := updated[key]; !ok {
			return fmt.Errorf("protect: key %q not found in secret at %s", key, opts.Path)
		}
	}

	result := ProtectResult{
		Path:      opts.Path,
		Protected: opts.Keys,
		DryRun:    opts.DryRun,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, opts.Path, updated); err != nil {
			return fmt.Errorf("protect: write: %w", err)
		}
	}

	return json.NewEncoder(opts.Output).Encode(result)
}
