// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Rotate
//
// Rotate re-writes the current version of a Vault KV-v2 secret in-place,
// producing a new version with identical values. This is useful for triggering
// version-aware workflows (e.g. certificate renewal hooks) without actually
// changing any secret material.
//
// Usage:
//
//	err := cmd.Rotate(cmd.RotateOptions{
//		Address: "https://vault.example.com",
//		Token:   "s.xxxx",
//		Mount:   "secret",
//		Path:    "myapp/config",
//		Mask:    true,
//		DryRun:  false,
//		Output:  os.Stdout,
//	})
//
// When DryRun is true the secret is fetched but not re-written; the result
// JSON will contain dry_run:true and new_version will be omitted.
//
// When Mask is false the secret data is included in the JSON output.
package cmd
