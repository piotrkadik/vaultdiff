// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # Encrypt
//
// Encrypt re-encodes selected keys of a Vault KV-v2 secret using a simple
// hex-prefix scheme, allowing operators to verify that sensitive values are
// obfuscated before being stored or exported.
//
// Typical usage:
//
//	result, err := cmd.Encrypt(cmd.EncryptOptions{
//		Address: "http://127.0.0.1:8200",
//		Token:   "root",
//		Mount:   "secret",
//		Path:    "myapp/config",
//		Keys:    []string{"password", "api_key"},
//		DryRun:  true,
//	})
//
// When DryRun is true the secret is read but never written back to Vault.
// Set Mask to false to include the encoded values in the returned EncryptResult.
package cmd
