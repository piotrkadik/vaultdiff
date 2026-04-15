package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/audit"
	"github.com/example/vaultdiff/internal/config"
	"github.com/example/vaultdiff/internal/diff"
	"github.com/example/vaultdiff/internal/vault"
)

// DiffOptions controls the behaviour of a single diff invocation.
type DiffOptions struct {
	Flags  Flags
	Cfg    config.Config
	Output io.Writer
	Audit  *audit.Logger
}

// DefaultDiffOptions returns DiffOptions populated from the environment.
func DefaultDiffOptions(f Flags) DiffOptions {
	return DiffOptions{
		Flags:  f,
		Cfg:    config.Default(),
		Output: os.Stdout,
		Audit:  audit.NewLogger(nil),
	}
}

// Diff fetches two secret versions from Vault, compares them, and writes
// the result to opts.Output using the configured format.
func Diff(ctx context.Context, opts DiffOptions) (int, error) {
	client, err := vault.NewClient(opts.Cfg.Address, opts.Cfg.Token, opts.Cfg.Mount)
	if err != nil {
		return ExitError, fmt.Errorf("vault client: %w", err)
	}

	a, err := vault.GetSecretVersion(ctx, client, opts.Flags.Path, opts.Flags.VersionA)
	if err != nil {
		return ExitError, fmt.Errorf("fetch version %d: %w", opts.Flags.VersionA, err)
	}

	b, err := vault.GetSecretVersion(ctx, client, opts.Flags.Path, opts.Flags.VersionB)
	if err != nil {
		return ExitError, fmt.Errorf("fetch version %d: %w", opts.Flags.VersionB, err)
	}

	changes := diff.Compare(a.Data, b.Data, opts.Flags.Path)

	report := diff.BuildReport(opts.Flags.Path, opts.Flags.VersionA, opts.Flags.VersionB, changes)

	outWriter := opts.Output
	if outWriter == nil {
		outWriter = os.Stdout
	}

	if err := WriteOutput(outWriter, changes, DefaultOutputOptions(opts.Cfg)); err != nil {
		return ExitError, fmt.Errorf("write output: %w", err)
	}

	if opts.Audit != nil {
		_ = opts.Audit.Write(audit.Entry{
			Path:     opts.Flags.Path,
			VersionA: opts.Flags.VersionA,
			VersionB: opts.Flags.VersionB,
			Summary:  diff.Summarize(changes),
		})
	}

	return ExitCodeForDrift(report), nil
}
