package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultEncryptOptions returns EncryptOptions populated from environment variables.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
		DryRun:  false,
	}
}

// EncryptOptions configures the Encrypt command.
type EncryptOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Keys    []string
	Mask    bool
	Format  string
	DryRun  bool
	Output  io.Writer
}

// EncryptResult holds the outcome of an Encrypt operation.
type EncryptResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Keys    []string          `json:"keys"`
	DryRun  bool              `json:"dry_run"`
	Data    map[string]string `json:"data,omitempty"`
}

// Encrypt re-writes selected keys of a secret with base64-encoded values,
// simulating an at-rest encryption layer for audit purposes.
func Encrypt(opts EncryptOptions) (*EncryptResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, 0)
	if err != nil {
		return nil, fmt.Errorf("encrypt: fetch %q: %w", opts.Path, err)
	}

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range secret.Data {
			targetKeys = append(targetKeys, k)
		}
	}

	encrypted := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		encrypted[k] = v
	}
	for _, k := range targetKeys {
		if v, ok := encrypted[k]; ok {
			encrypted[k] = fmt.Sprintf("enc:%x", []byte(v))
		}
	}

	result := &EncryptResult{
		Path:   opts.Path,
		Keys:   targetKeys,
		DryRun: opts.DryRun,
	}

	if !opts.DryRun {
		ver, err := vault.WriteSecret(client, opts.Path, encrypted)
		if err != nil {
			return nil, fmt.Errorf("encrypt: write %q: %w", opts.Path, err)
		}
		result.Version = ver
		if !opts.Mask {
			result.Data = encrypted
		}
	}

	if opts.Format == "json" {
		return result, json.NewEncoder(opts.Output).Encode(result)
	}
	fmt.Fprintf(opts.Output, "encrypted %d key(s) at %s (dry_run=%v)\n", len(targetKeys), opts.Path, opts.DryRun)
	return result, nil
}
