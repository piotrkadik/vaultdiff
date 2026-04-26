// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Format
//
// Format fetches a specific version of a Vault KV-v2 secret and renders its
// key/value pairs in a human-readable or machine-readable form.
//
// Supported output formats:
//
//   - text (default): prints "key=value" lines prefixed with a path header.
//   - json: emits a JSON object containing path, version, keys, and
//     optionally the raw data map when masking is disabled.
//
// Value masking
//
// When Mask is true (the default), secret values are replaced with "***" in
// the output. Set the environment variable VAULTDIFF_MASK=false to disable
// masking globally.
//
// Environment variables
//
//	VAULT_ADDR        Vault server address
//	VAULT_MOUNT       KV mount path (default: secret)
//	VAULT_TOKEN       Vault authentication token
//	VAULTDIFF_MASK    Set to "false" to reveal secret values
//	VAULTDIFF_FORMAT  Output format: "" (text) or "json"
package cmd
