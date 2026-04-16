// Package cmd provides CLI command implementations for vaultdiff.
package cmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// EnvOptions controls the behaviour of the Env command.
type EnvOptions struct {
	// Output is the writer used for printing environment variable mappings.
	// Defaults to os.Stdout when nil.
	Output io.Writer

	// Mask replaces secret values with "***" in the output.
	Mask bool
}

// DefaultEnvOptions returns an EnvOptions populated with sensible defaults.
func DefaultEnvOptions() EnvOptions {
	return EnvOptions{
		Output: os.Stdout,
		Mask:   true,
	}
}

// PrintEnv fetches the given version of a secret and prints its key/value
// pairs in KEY=VALUE shell-export format, suitable for sourcing or inspection.
func PrintEnv(client vault.SecretGetter, path string, version int, opts EnvOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	sv, err := client.GetSecretVersion(path, version)
	if err != nil {
		return fmt.Errorf("env: fetch %s@v%d: %w", path, version, err)
	}

	// Collect and sort keys for deterministic output.
	keys := make([]string, 0, len(sv.Data))
	for k := range sv.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val := fmt.Sprintf("%v", sv.Data[k])
		if opts.Mask {
			val = "***"
		}
		// Escape embedded newlines so the output stays sourceable.
		val = strings.ReplaceAll(val, "\n", "\\n")
		fmt.Fprintf(opts.Output, "%s=%s\n", k, val)
	}

	return nil
}
