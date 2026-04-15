// Package cmd provides CLI command implementations for vaultdiff.
package cmd

import (
	"context"
	"fmt"

	"github.com/vaultdiff/internal/vault"
)

// ResolveOptions controls how version resolution behaves.
type ResolveOptions struct {
	// Mount is the KV v2 mount path.
	Mount string
	// SkipDestroyed causes destroyed versions to be excluded.
	SkipDestroyed bool
}

// DefaultResolveOptions returns sensible defaults for version resolution.
func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		Mount:         "secret",
		SkipDestroyed: true,
	}
}

// ResolvedVersions holds the two version numbers chosen for a diff.
type ResolvedVersions struct {
	From int
	To   int
}

// ResolveVersions determines the concrete from/to version numbers for the
// given secret path. If toVersion is 0 the latest non-destroyed version is
// used. If fromVersion is 0 the version immediately preceding toVersion is
// used. An error is returned when resolution is not possible (e.g. the path
// has fewer than two versions).
func ResolveVersions(
	ctx context.Context,
	client *vault.Client,
	path string,
	fromVersion, toVersion int,
	opts ResolveOptions,
) (ResolvedVersions, error) {
	if fromVersion > 0 && toVersion > 0 {
		return ResolvedVersions{From: fromVersion, To: toVersion}, nil
	}

	meta, err := vault.ListVersions(ctx, client, opts.Mount, path)
	if err != nil {
		return ResolvedVersions{}, fmt.Errorf("resolve versions: list %q: %w", path, err)
	}

	latest := vault.LatestVersion(meta, opts.SkipDestroyed)
	if latest <= 0 {
		return ResolvedVersions{}, fmt.Errorf("resolve versions: no usable versions found for %q", path)
	}

	to := toVersion
	if to == 0 {
		to = latest
	}

	from := fromVersion
	if from == 0 {
		from = to - 1
	}

	if from <= 0 {
		return ResolvedVersions{}, fmt.Errorf("resolve versions: cannot determine a prior version for %q (latest=%d)", path, latest)
	}

	return ResolvedVersions{From: from, To: to}, nil
}
