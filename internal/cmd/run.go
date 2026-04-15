// Package cmd wires together configuration, vault access, diffing, and
// reporting into the top-level vaultdiff command.
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

// RunOptions holds the parsed CLI flags for a single diff run.
type RunOptions struct {
	Path       string
	VersionA   int
	VersionB   int
	ShowAll    bool
	Mask       bool
	Format     string
	OutputFile string
}

// Run executes a full diff between two secret versions and writes the result
// to out (or a file when OutputFile is set).
func Run(ctx context.Context, opts RunOptions, cfg config.Config, out io.Writer) error {
	client, err := vault.NewClient(cfg.Address, cfg.Token, cfg.Mount)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	secA, err := vault.GetSecretVersion(ctx, client, opts.Path, opts.VersionA)
	if err != nil {
		return fmt.Errorf("fetch version %d: %w", opts.VersionA, err)
	}

	secB, err := vault.GetSecretVersion(ctx, client, opts.Path, opts.VersionB)
	if err != nil {
		return fmt.Errorf("fetch version %d: %w", opts.VersionB, err)
	}

	changes := diff.Compare(secA.Data, secB.Data)
	summary := diff.Summarize(changes)

	report := diff.Report{
		Path:    opts.Path,
		Changes: changes,
		Summary: summary,
	}

	writer := out
	if opts.OutputFile != "" {
		f, err := os.Create(opts.OutputFile)
		if err != nil {
			return fmt.Errorf("open output file: %w", err)
		}
		defer f.Close()
		writer = f
	}

	if err := diff.Render(report, diff.RenderOptions{
		ShowUnchanged: opts.ShowAll,
		Mask:          opts.Mask,
		Color:         opts.OutputFile == "",
	}, writer); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	logger := audit.NewLogger(nil)
	_ = logger.Write(audit.Entry{
		Path:     opts.Path,
		VersionA: opts.VersionA,
		VersionB: opts.VersionB,
		Summary:  summary,
	})

	return nil
}
