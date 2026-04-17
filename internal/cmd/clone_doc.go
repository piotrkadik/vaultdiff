// Package cmd provides high-level command operations for vaultdiff.
//
// # Clone
//
// Clone copies all key/value pairs from one Vault secret path to another.
// It reads the latest version of the source secret and writes those values
// to the destination path as a new version.
//
// Use DryRun to preview the operation without writing to Vault.
//
// Example:
//
//	result, err := cmd.Clone(ctx, cmd.CloneOptions{
//		Address: "http://127.0.0.1:8200",
//		Token:   "root",
//		Mount:   "secret",
//		SrcPath: "app/production",
//		DstPath: "app/staging",
//		DryRun:  false,
//	})
package cmd
