package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// TokenizeOptions configures the Tokenize command.
type TokenizeOptions struct {
	Address string
	Mount   string
	Path    string
	Version int
	Format  string
	Mask    bool
	Output  io.Writer
}

// DefaultTokenizeOptions returns options populated from the environment.
func DefaultTokenizeOptions() TokenizeOptions {
	mask := true
	if strings.EqualFold(os.Getenv("VAULTDIFF_MASK"), "false") {
		mask = false
	}
	return TokenizeOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   os.Getenv("VAULT_MOUNT"),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
		Mask:    mask,
	}
}

// TokenizeResult holds the output of a Tokenize run.
type TokenizeResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Tokens  map[string]string `json:"tokens"`
}

// Tokenize fetches a secret version and replaces each value with a
// deterministic placeholder token of the form <KEY_NAME>, useful for
// generating redacted templates or schema stubs.
func Tokenize(client vault.SecretGetter, opts TokenizeOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	secret, err := client.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("tokenize: fetch %q v%d: %w", opts.Path, opts.Version, err)
	}

	tokens := make(map[string]string, len(secret.Data))
	for k := range secret.Data {
		tokens[k] = "<" + strings.ToUpper(k) + ">"
	}

	result := TokenizeResult{
		Path:    opts.Path,
		Version: secret.Version,
		Tokens:  tokens,
	}

	if opts.Format == "json" {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	fmt.Fprintf(out, "path: %s (v%d)\n", result.Path, result.Version)
	for k, v := range result.Tokens {
		fmt.Fprintf(out, "  %-24s %s\n", k, v)
	}
	return nil
}
