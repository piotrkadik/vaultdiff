// Package cmd provides CLI command implementations for vaultdiff.
//
// # Restore
//
// Restore re-writes a specific historical version of a Vault secret as a new
// latest version. This is useful when a secret has been accidentally overwritten
// and you need to recover a known-good state without manually copying values.
//
// The operation fetches the requested version from the KV v2 metadata and
// writes its data back to the same path, producing a new version number.
//
// When DryRun is true the target secret is fetched and the result is printed
// but no write is performed.
//
// Example:
//
//	result, err := cmd.Restore(ctx, cmd.RestoreOptions{
//		Address: "https://vault.example.com",
//		Path:    "myapp/config",
//		Version: 3,
//	})
package cmd
