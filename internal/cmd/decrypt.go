package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DecryptOptions controls the behaviour of Decrypt.
type DecryptOptions struct {
	Address string
	Mount   string
	Token   string
	Path    string
	Version int
	Format  string
	Mask    bool
	Output  io.Writer
}

// DefaultDecryptOptions returns options populated from environment variables.
func DefaultDecryptOptions() DecryptOptions {
	return DecryptOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// DecryptResult holds the decrypted key/value pairs for a secret version.
type DecryptResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Data    map[string]string `json:"data,omitempty"`
}

// Decrypt fetches a secret version and writes the plaintext key/value pairs.
func Decrypt(client vault.SecretGetter, opts DecryptOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	sv, err := client.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("decrypt: fetch %s@%d: %w", opts.Path, opts.Version, err)
	}

	data := make(map[string]string, len(sv.Data))
	for k, v := range sv.Data {
		if opts.Mask {
			data[k] = "***"
		} else {
			data[k] = v
		}
	}

	res := DecryptResult{Path: opts.Path, Version: sv.Version, Data: data}

	if opts.Format == "json" {
		return json.NewEncoder(opts.Output).Encode(res)
	}

	fmt.Fprintf(opts.Output, "# %s version %d\n", res.Path, res.Version)
	for k, v := range res.Data {
		fmt.Fprintf(opts.Output, "%s=%s\n", k, v)
	}
	return nil
}
