// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Digest
//
// Digest fetches a specific version of a Vault KV secret and computes a
// deterministic SHA-256 fingerprint over its key-value pairs. The digest
// is computed by sorting keys lexicographically and hashing each
// "key=value\n" line, making it suitable for integrity checks, change
// detection, and audit trails without exposing secret values directly.
//
// Usage:
//
//	result, err := cmd.Digest("secret/myapp", cmd.DefaultDigestOptions())
//
// Output formats:
//   - text (default): human-readable label/value lines
//   - json: structured DigestResult object
//
// Environment variables:
//   - VAULT_ADDR:       Vault server address
//   - VAULT_TOKEN:      authentication token
//   - VAULTDIFF_MOUNT:  KV mount path (default: "secret")
//   - VAULTDIFF_FORMAT: output format ("" or "json")
//   - VAULTDIFF_MASK:   set to "false" to disable value masking
package cmd
