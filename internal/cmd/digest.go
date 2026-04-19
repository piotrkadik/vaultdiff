package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DigestOptions configures the Digest command.
type DigestOptions struct {
	Address string
	Mount   string
	Token   string
	Version int
	Mask    bool
	Format  string
	Output  io.Writer
}

// DefaultDigestOptions returns options populated from environment variables.
func DefaultDigestOptions() DigestOptions {
	return DigestOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// DigestResult holds the computed digest for a secret version.
type DigestResult struct {
	Path    string `json:"path"`
	Version int    `json:"version"`
	Digest  string `json:"digest"`
	Keys    int    `json:"keys"`
}

// Digest fetches a secret version and computes a deterministic SHA-256 digest
// over its key-value pairs, useful for integrity checks without exposing values.
func Digest(path string, opts DigestOptions) (*DigestResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	c, err := vault.NewClient(opts.Address, opts.Mount, opts.Token)
	if err != nil {
		return nil, fmt.Errorf("digest: %w", err)
	}

	secret, err := c.GetSecretVersion(path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("digest: %w", err)
	}

	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secret.Data[k])
	}
	digest := hex.EncodeToString(h.Sum(nil))

	result := &DigestResult{
		Path:    path,
		Version: secret.Version,
		Digest:  digest,
		Keys:    len(keys),
	}

	if opts.Format == "json" {
		return result, json.NewEncoder(opts.Output).Encode(result)
	}
	fmt.Fprintf(opts.Output, "path:    %s\nversion: %d\nkeys:    %d\ndigest:  %s\n",
		result.Path, result.Version, result.Keys, result.Digest)
	return result, nil
}
