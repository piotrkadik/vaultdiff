// Package cmd provides high-level CLI command implementations for vaultdiff.
//
// # Redact
//
// Redact replaces the values of one or more keys in a Vault KV-v2 secret with
// the literal string "REDACTED" and writes the result as a new version.
//
// When Keys is empty every key in the secret is redacted. Set DryRun to true
// to preview which keys would be affected without performing the write.
//
// Example usage:
//
//	result, err := cmd.Redact(cmd.RedactOptions{
//		Address: "https://vault.example.com",
//		Token:   "s.abc123",
//		Path:    "myapp/config",
//		Keys:    []string{"password", "api_key"},
//	})
package cmd
