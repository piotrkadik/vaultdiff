// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # Extract
//
// Extract fetches a specific version of a Vault secret and returns only the
// requested subset of keys. When no keys are specified every key in the secret
// is returned, making it a lightweight alternative to Fetch when callers only
// need a handful of fields.
//
// # Options
//
//   - Address  – Vault server address (VAULT_ADDR)
//   - Token    – Vault token (VAULT_TOKEN)
//   - Mount    – KV v2 mount path (VAULT_MOUNT, default "secret")
//   - Path     – secret path inside the mount
//   - Version  – secret version (0 = latest)
//   - Keys     – whitelist of keys to extract; empty means all keys
//   - Format   – output format: "" / "text" or "json"
//   - Mask     – replace values with *** in output (default true)
//
// # Exit behaviour
//
// Extract returns an error when a requested key is absent from the secret,
// allowing callers to treat a missing key as a hard failure.
package cmd
