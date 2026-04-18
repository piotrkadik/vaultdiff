package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultdiff/internal/vault"
)

// SchemaOptions configures the Schema command.
type SchemaOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Format  string
	Output  io.Writer
}

// DefaultSchemaOptions returns options populated from environment variables.
func DefaultSchemaOptions() SchemaOptions {
	return SchemaOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  formatFromEnv(),
	}
}

// SchemaField describes a single key in a secret.
type SchemaField struct {
	Key      string `json:"key"`
	Present  bool   `json:"present"`
	NonEmpty bool   `json:"non_empty"`
}

// SchemaResult holds the derived schema for a secret version.
type SchemaResult struct {
	Path    string        `json:"path"`
	Version int           `json:"version"`
	Fields  []SchemaField `json:"fields"`
}

// Schema fetches a secret version and prints its key schema.
func Schema(c vault.Client, opts SchemaOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	sv, err := c.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("schema: fetch %s@%d: %w", opts.Path, opts.Version, err)
	}

	keys := make([]string, 0, len(sv.Data))
	for k := range sv.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fields := make([]SchemaField, 0, len(keys))
	for _, k := range keys {
		v := sv.Data[k]
		fields = append(fields, SchemaField{
			Key:      k,
			Present:  true,
			NonEmpty: v != "",
		})
	}

	result := SchemaResult{
		Path:    opts.Path,
		Version: sv.Version,
		Fields:  fields,
	}

	if opts.Format == "json" {
		return json.NewEncoder(opts.Output).Encode(result)
	}

	fmt.Fprintf(opts.Output, "schema: %s @ v%d\n", result.Path, result.Version)
	for _, f := range result.Fields {
		empty := ""
		if !f.NonEmpty {
			empty = " (empty)"
		}
		fmt.Fprintf(opts.Output, "  %s%s\n", f.Key, empty)
	}
	return nil
}
