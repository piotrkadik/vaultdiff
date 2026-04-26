package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultTreeOptions returns TreeOptions populated from environment variables.
func DefaultTreeOptions() TreeOptions {
	return TreeOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// TreeOptions configures the Tree command.
type TreeOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Format  string
	Output  io.Writer
}

// TreeResult holds the hierarchical key listing.
type TreeResult struct {
	Path  string   `json:"path"`
	Keys  []string `json:"keys"`
	Count int      `json:"count"`
}

// Tree fetches all keys under a path and renders them as a tree.
func Tree(client vault.SecretGetter, opts TreeOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	data, err := client.GetSecretVersion(opts.Path, 0)
	if err != nil {
		return fmt.Errorf("tree: fetch %q: %w", opts.Path, err)
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := TreeResult{
		Path:  opts.Path,
		Keys:  keys,
		Count: len(keys),
	}

	if opts.Format == "json" {
		return json.NewEncoder(out).Encode(result)
	}

	fmt.Fprintf(out, "%s\n", opts.Path)
	for i, k := range keys {
		prefix := "├── "
		if i == len(keys)-1 {
			prefix = "└── "
		}
		fmt.Fprintf(out, "%s%s\n", prefix, strings.TrimSpace(k))
	}
	fmt.Fprintf(out, "\n%d key(s)\n", result.Count)
	return nil
}
