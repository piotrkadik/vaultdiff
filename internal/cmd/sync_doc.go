// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Sync
//
// Sync copies all key-value pairs from a source Vault secret path to a
// destination path. It is useful for propagating secrets across environments
// (e.g. production → staging) in a single, auditable operation.
//
// Basic usage:
//
//	err := cmd.Sync(cmd.SyncOptions{
//		Address: "https://vault.example.com",
//		Token:   token,
//		Mount:   "secret",
//		SrcPath: "myapp/prod",
//		DstPath: "myappt	DryRun:  false,
//		Mask:    true,
//	})
//
// When DryRun is true the destination is never written; the result JSON will
// still reflect what would have been synced.
//
// When Mask is true secret values are omitted from the JSON output.
package cmd
