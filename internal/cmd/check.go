package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultCheckOptions returns CheckOptions populated from environment variables.
func DefaultCheckOptions() CheckOptions {
	return CheckOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// CheckOptions configures the Check command.
type CheckOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Keys    []string
	Format  string
	Output  io.Writer
}

// CheckResult holds the outcome of a key-existence check.
type CheckResult struct {
	Path    string         `json:"path"`
	Version int            `json:"version"`
	Present map[string]bool `json:"present"`
	AllOK   bool           `json:"all_ok"`
}

// Check verifies that the given keys exist (and are non-empty) in the secret
// at the specified path and version.
func Check(opts CheckOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	c, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}

	data, err := vault.GetSecretVersion(c, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("check: %w", err)
	}

	present := make(map[string]bool, len(opts.Keys))
	allOK := true
	for _, k := range opts.Keys {
		v, exists := data[k]
		ok := exists && v != ""
		present[k] = ok
		if !ok {
			allOK = false
		}
	}

	res := CheckResult{
		Path:    opts.Path,
		Version: opts.Version,
		Present: present,
		AllOK:   allOK,
	}

	if opts.Format == "json" {
		return json.NewEncoder(out).Encode(res)
	}

	fmt.Fprintf(out, "path: %s  version: %d\n", res.Path, res.Version)
	for _, k := range opts.Keys {
		status := "OK"
		if !res.Present[k] {
			status = "MISSING"
		}
		fmt.Fprintf(out, "  %-30s %s\n", k, status)
	}
	if !res.AllOK {
		return fmt.Errorf("check: one or more required keys are missing or empty")
	}
	return nil
}
