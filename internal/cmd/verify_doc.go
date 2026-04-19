// Package cmd provides high-level command implementations for vaultdiff.
//
// # Verify
//
// Verify fetches a specific version of a Vault secret and confirms that it is
// both readable and contains at least one key. It is useful as a lightweight
// health-check after a write, promote, or rotate operation.
//
// Basic usage:
//
//	err := cmd.Verify("secret/myapp", 3, cmd.DefaultVerifyOptions(), os.Stdout)
//
// Output formats:
//
//	text (default) — human-readable single line summary
//	json           — machine-readable VerifyResult struct
//
// Environment variables:
//
//	VAULT_ADDR        Vault server address
//	VAULT_TOKEN       Authentication token
//	VAULTDIFF_MOUNT   KV mount path (default: secret)
//	VAULTDIFF_MASK    Set to "false" to include plaintext values in output
//	VAULTDIFF_FORMAT  Output format: text or json
package cmd
