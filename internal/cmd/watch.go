package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/vaultdiff/internal/diff"
	"github.com/vaultdiff/internal/vault"
)

// WatchOptions configures the watch loop behaviour.
type WatchOptions struct {
	Interval time.Duration
	MaxRounds int // 0 means run until cancelled
	Output   io.Writer
}

// DefaultWatchOptions returns sensible defaults for watch mode.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval:  30 * time.Second,
		MaxRounds: 0,
		Output:    os.Stdout,
	}
}

// Watch polls Vault on a fixed interval and prints a diff whenever the secret
// at path changes between consecutive versions. It returns after ctx is
// cancelled, or after MaxRounds iterations when MaxRounds > 0.
func Watch(ctx context.Context, client *vault.Client, path string, opts WatchOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	outOpts := DefaultOutputOptions()

	round := 0
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	var prevVersion int

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			latest, err := vault.LatestVersion(ctx, client, path)
			if err != nil {
				fmt.Fprintf(opts.Output, "watch: error fetching versions for %s: %v\n", path, err)
				continue
			}

			if prevVersion == 0 {
				prevVersion = latest
				fmt.Fprintf(opts.Output, "watch: tracking %s at version %d\n", path, latest)
			} else if latest != prevVersion {
				changes, diffErr := runDiff(ctx, client, path, prevVersion, latest)
				if diffErr != nil {
					fmt.Fprintf(opts.Output, "watch: diff error: %v\n", diffErr)
				} else {
					fmt.Fprintf(opts.Output, "\n=== %s v%d → v%d ===\n", path, prevVersion, latest)
					_ = WriteOutput(changes, outOpts, opts.Output)
				}
				prevVersion = latest
			}

			round++
			if opts.MaxRounds > 0 && round >= opts.MaxRounds {
				return nil
			}
		}
	}
}

// runDiff fetches both versions and returns the list of Changes.
func runDiff(ctx context.Context, client *vault.Client, path string, from, to int) ([]diff.Change, error) {
	a, err := vault.GetSecretVersion(ctx, client, path, from)
	if err != nil {
		return nil, fmt.Errorf("fetch v%d: %w", from, err)
	}
	b, err := vault.GetSecretVersion(ctx, client, path, to)
	if err != nil {
		return nil, fmt.Errorf("fetch v%d: %w", to, err)
	}
	return diff.Compare(a.Data, b.Data), nil
}
