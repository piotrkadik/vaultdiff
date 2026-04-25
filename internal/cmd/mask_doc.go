// Package cmd – mask command
//
// # Mask
//
// Mask replaces the values of one or more secret keys with a redaction marker
// ("***") and, unless --dry-run is set, writes the sanitised secret back to
// Vault.  This is useful when a secret must remain in Vault but its plaintext
// value should no longer be accessible to operators who have read access.
//
// # Usage
//
//	vaultdiff mask <path> [--key <key>]... [--version <n>] [--dry-run]
//
// If no --key flags are provided every key in the secret is masked.
//
// # Output
//
// The command always emits a JSON object describing the operation:
//
//	{
//	  "path":        "secret/myapp",
//	  "version":     3,
//	  "dry_run":     false,
//	  "masked_keys": ["password", "token"]
//	}
//
// When the global --mask flag is disabled (VAULTDIFF_MASK=false) the redacted
// data map is included in the output for inspection purposes.
package cmd
