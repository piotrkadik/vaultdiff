// Package cmd provides CLI command implementations for vaultdiff.
//
// # Rename
//
// Rename reads a secret at the given path, renames a key within the secret's
// data map, and writes the updated map back to Vault as a new version.
//
// The old key is removed and the value is preserved under the new key name.
// All other keys in the secret are left unchanged.
//
// Dry-run mode (DryRun: true) performs the fetch and validation steps but
// skips the write, printing the intended result as JSON.
//
// Environment variables:
//
//	VAULT_ADDR   – Vault server address
//	VAULT_TOKEN  – Vault token
//	VAULT_MOUNT  – KV v2 mount path (default: "secret")
package cmd
