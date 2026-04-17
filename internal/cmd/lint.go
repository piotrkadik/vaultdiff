package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/vaultdiff/internal/vault"
)

// LintOptions configures the Lint command.
type LintOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Format  string
	Output  io.Writer
}

// DefaultLintOptions returns options populated from environment variables.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  formatFromEnv(),
	}
}

// LintResult holds the outcome of a lint check for a single key.
type LintResult struct {
	Key     string `json:"key"`
	Issue   string `json:"issue"`
	Severity string `json:"severity"`
}

// Lint fetches a secret version and reports keys with empty or suspicious values.
func Lint(opts LintOptions) ([]LintResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	c, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("lint: %w", err)
	}

	secret, err := c.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("lint: %w", err)
	}

	var results []LintResult
	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secret.Data[k]
		if v == "" {
			results = append(results, LintResult{Key: k, Issue: "empty value", Severity: "warning"})
		} else if v == "changeme" || v == "TODO" || v == "FIXME" {
			results = append(results, LintResult{Key: k, Issue: "placeholder value: " + v, Severity: "error"})
		}
	}

	if opts.Format == "json" {
		return results, json.NewEncoder(opts.Output).Encode(results)
	}

	if len(results) == 0 {
		fmt.Fprintln(opts.Output, "no issues found")
		return results, nil
	}
	for _, r := range results {
		fmt.Fprintf(opts.Output, "[%s] %s: %s\n", r.Severity, r.Key, r.Issue)
	}
	return results, nil
}
