// Package cmd provides CLI command implementations for vaultdiff.
package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/vaultdiff/internal/diff"
	"github.com/vaultdiff/internal/vault"
)

// CompareOptions controls the behaviour of the Compare command.
type CompareOptions struct {
	// PathA is the Vault secret path for the left-hand side.
	PathA string
	// PathB is the Vault secret path for the right-hand side.
	PathB string
	// VersionA is the version number to fetch for PathA (0 = latest).
	VersionA int
	// VersionB is the version number to fetch for PathB (0 = latest).
	VersionB int
	// Mask controls whether secret values are redacted in output.
	Mask bool
	// Output is where rendered results are written; defaults to os.Stdout.
	Output io.Writer
}

// DefaultCompareOptions returns a CompareOptions with sensible defaults.
func DefaultCompareOptions() CompareOptions {
	return CompareOptions{
		Mask:   true,
		Output: os.Stdout,
	}
}

// Compare fetches two secrets from Vault (potentially different paths and/or
// versions) and writes a diff to opts.Output.
func Compare(ctx context.Context, client *vault.Client, opts CompareOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	secretA, err := client.GetSecretVersion(ctx, opts.PathA, opts.VersionA)
	if err != nil {
		return fmt.Errorf("fetching %s@%d: %w", opts.PathA, opts.VersionA, err)
	}

	secretB, err := client.GetSecretVersion(ctx, opts.PathB, opts.VersionB)
	if err != nil {
		return fmt.Errorf("fetching %s@%d: %w", opts.PathB, opts.VersionB, err)
	}

	changes := diff.Compare(secretA.Data, secretB.Data)

	renderOpts := diff.RenderOptions{
		Mask:          opts.Mask,
		Color:         true,
		ShowUnchanged: false,
	}

	header := fmt.Sprintf("--- %s (v%d)\n+++ %s (v%d)\n",
		opts.PathA, secretA.Version,
		opts.PathB, secretB.Version,
	)
	if _, err := fmt.Fprint(opts.Output, header); err != nil {
		return err
	}

	return diff.Render(opts.Output, changes, renderOpts)
}
