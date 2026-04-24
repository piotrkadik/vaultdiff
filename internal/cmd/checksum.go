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

// DefaultChecksumOptions returns ChecksumOptions populated from environment variables.
func DefaultChecksumOptions() ChecksumOptions {
	return ChecksumOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// ChecksumOptions configures the Checksum command.
type ChecksumOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Format  string
	Mask    bool
	Output  io.Writer
}

// ChecksumResult holds the computed checksum for a secret version.
type ChecksumResult struct {
	Path     string `json:"path"`
	Version  int    `json:"version"`
	Checksum string `json:"checksum"`
	Keys     int    `json:"keys"`
}

// Checksum fetches a secret version and computes a deterministic SHA-256
// checksum over its sorted key=value pairs.
func Checksum(c vault.Client, opts ChecksumOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	sv, err := c.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("checksum: fetch %s@%d: %w", opts.Path, opts.Version, err)
	}

	keys := make([]string, 0, len(sv.Data))
	for k := range sv.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, sv.Data[k])
	}
	sum := hex.EncodeToString(h.Sum(nil))

	res := ChecksumResult{
		Path:     opts.Path,
		Version:  sv.Version,
		Checksum: sum,
		Keys:     len(sv.Data),
	}

	if opts.Format == "json" {
		return json.NewEncoder(out).Encode(res)
	}

	fmt.Fprintf(out, "path:     %s\n", res.Path)
	fmt.Fprintf(out, "version:  %d\n", res.Version)
	fmt.Fprintf(out, "keys:     %d\n", res.Keys)
	fmt.Fprintf(out, "checksum: %s\n", res.Checksum)
	return nil
}
