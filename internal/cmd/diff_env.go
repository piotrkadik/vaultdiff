package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultdiff/internal/diff"
	"github.com/yourusername/vaultdiff/internal/vault"
)

// DiffEnvOptions configures cross-environment secret diffing.
type DiffEnvOptions struct {
	Address  string
	Mount    string
	Path     string
	EnvA     string
	EnvB     string
	Mask     bool
	Format   string
	Output   io.Writer
}

// DefaultDiffEnvOptions returns options populated from environment variables.
func DefaultDiffEnvOptions() DiffEnvOptions {
	mask := true
	if os.Getenv("VAULT_MASK_VALUES") == "false" {
		mask = false
	}
	return DiffEnvOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Mask:    mask,
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// DiffEnv fetches the same secret path from two environment-prefixed paths
// and renders a diff between them.
func DiffEnv(opts DiffEnvOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Mount)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	pathA := opts.EnvA + "/" + opts.Path
	pathB := opts.EnvB + "/" + opts.Path

	secA, err := vault.GetSecretVersion(client, pathA, 0)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", pathA, err)
	}
	secB, err := vault.GetSecretVersion(client, pathB, 0)
	if err != nil {
		return fmt.Errorf("fetch %s: %w", pathB, err)
	}

	changes := diff.Compare(secA.Data, secB.Data)

	fmt.Fprintf(out, "--- %s\n+++ %s\n", pathA, pathB)
	return WriteOutput(out, changes, opts.Format, opts.Mask)
}
