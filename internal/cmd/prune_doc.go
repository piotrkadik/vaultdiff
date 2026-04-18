// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Prune
//
// Prune removes old secret versions at a given path, keeping only the N most
// recent non-destroyed versions. This is useful for enforcing retention
// policies and reducing storage in busy paths.
//
// Basic usage:
//
//	err := cmd.Prune(ctx, cmd.PruneOptions{
//		Address: "https://vault.example.com",
//		Token:   token,
//		Mount:   "secret",
//		Path:    "secret/data/myapp/prod",
//		Keep:    5,
//		DryRun:  true,
//	})
//
// When DryRun is true the versions to be removed are reported but no
// destructive API calls are made. The result is written as a JSON object to
// the configured Output writer (defaulting to os.Stdout).
package cmd
