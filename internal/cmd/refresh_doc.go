// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # Refresh
//
// Refresh re-writes the current version of a Vault KV-v2 secret back to itself,
// effectively bumping the secret's updated_time metadata without altering any
// values.  This is useful when an external system tracks staleness by
// modification time and you need to signal that a secret is still current.
//
// Usage:
//
//	result, err := cmd.Refresh(client, writer, cmd.RefreshOptions{
//		Address: "https://vault.example.com",
//		Path:    "services/api",
//		DryRun:  false,
//		Mask:    true,
//	})
//
// Dry-run mode reads the secret and reports what would happen without writing.
// Set Format to "json" to receive machine-readable output.
package cmd
