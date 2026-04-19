// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # Import
//
// Import reads a JSON file and writes its key/value pairs to a Vault KV v2
// secret at the specified path. It is the complement of the Fetch and Snapshot
// commands and is useful for seeding or restoring secrets from a local file.
//
// Dry-run mode parses and validates the input without contacting Vault, making
// it safe to use in CI pipelines for pre-flight checks.
//
// Example usage:
//
//	result, err := cmd.Import(file, cmd.ImportOptions{
//		Address: "https://vault.example.com",
//		Token:   token,
//		Mount:   "secret",
//		Path:    "myapp/config",
//		DryRun:  false,
//		Mask:    true,
//	})
package cmd
