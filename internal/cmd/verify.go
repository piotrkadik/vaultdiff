package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your/vaultdiff/internal/vault"
)

// VerifyOptions configures the Verify command.
type VerifyOptions struct {
	Address string
	Mount   string
	Token   string
	Mask    bool
	Format  string
}

// DefaultVerifyOptions returns options populated from environment variables.
func DefaultVerifyOptions() VerifyOptions {
	return VerifyOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// VerifyResult holds the outcome of a Verify run.
type VerifyResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Keys    int               `json:"keys"`
	Data    map[string]string `json:"data,omitempty"`
	OK      bool              `json:"ok"`
}

// Verify fetches a secret version and confirms it is readable and non-empty.
func Verify(path string, version int, opts VerifyOptions, out io.Writer) error {
	if out == nil {
		out = os.Stdout
	}

	c, err := vault.NewClient(opts.Address, opts.Mount, opts.Token)
	if err != nil {
		return fmt.Errorf("verify: client: %w", err)
	}

	data, err := vault.GetSecretVersion(c, path, version)
	if err != nil {
		return fmt.Errorf("verify: fetch: %w", err)
	}

	result := VerifyResult{
		Path:    path,
		Version: version,
		Keys:    len(data),
		OK:      len(data) > 0,
	}

	if !opts.Mask {
		result.Data = data
	}

	if opts.Format == "json" {
		return json.NewEncoder(out).Encode(result)
	}

	status := "OK"
	if !result.OK {
		status = "EMPTY"
	}
	fmt.Fprintf(out, "verify %s@v%d — %s (%d keys)\n", path, version, status, result.Keys)
	return nil
}
