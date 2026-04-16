package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vaultdiff/internal/diff"
	"github.com/vaultdiff/internal/vault"
)

// DefaultSearchOptions returns SearchOptions populated from environment.
func DefaultSearchOptions() SearchOptions {
	return SearchOptions{
		Format: os.Getenv("VAULTDIFF_FORMAT"),
		Mask:   os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// SearchOptions controls how Search behaves.
type SearchOptions struct {
	Query   string
	Format  string
	Mask    bool
	Output  io.Writer
}

// Search fetches a secret version and returns changes whose keys match query.
func Search(client *vault.Client, path string, version int, opts SearchOptions) ([]diff.Change, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	secret, err := client.GetSecretVersion(path, version)
	if err != nil {
		return nil, fmt.Errorf("search: fetch %s@%d: %w", path, version, err)
	}

	var changes []diff.Change
	for k, v := range secret.Data {
		if opts.Query == "" || strings.Contains(strings.ToLower(k), strings.ToLower(opts.Query)) {
			val := fmt.Sprintf("%v", v)
			changes = append(changes, diff.Change{
				Key:  k,
				Path: path,
				Type: diff.Unchanged,
				Old:  val,
				New:  val,
			})
		}
	}

	if opts.Format == "json" {
		return changes, json.NewEncoder(opts.Output).Encode(changes)
	}

	for _, c := range changes {
		v := c.New
		if opts.Mask {
			v = "***"
		}
		fmt.Fprintf(opts.Output, "%s\t%s\n", c.Key, v)
	}
	return changes, nil
}
