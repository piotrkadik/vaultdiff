// Package cmd provides CLI command implementations for vaultdiff.
//
// # Copy
//
// Copy reads a specific version of a Vault secret and writes its data to a
// destination path, optionally performing a dry-run that skips the write.
//
// Usage:
//
//	err := cmd.Copy(ctx, "prod/db", "staging/db", 3, cmd.DefaultCopyOptions())
//
// Dry-run mode fetches the source secret and reports what would be written
// without modifying Vault:
//
//	opts := cmd.DefaultCopyOptions()
//	opts.DryRun = true
//	err := cmd.Copy(ctx, "prod/db", "staging/db", 3, opts)
//
// Environment variables:
//
//	VAULT_ADDR        – Vault server address
//	VAULT_TOKEN       – authentication token
//	VAULT_MOUNT       – KV v2 mount path (default: secret)
//	VAULTDIFF_MASK    – set to "false" to reveal secret values in output
package cmd
