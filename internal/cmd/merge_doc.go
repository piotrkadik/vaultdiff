// Package cmd provides CLI command implementations for vaultdiff.
//
// # Merge
//
// Merge combines two Vault secrets by copying keys from a source path into a
// destination path. Keys that already exist in the destination are preserved;
// only missing keys are added from the source.
//
// This is useful for propagating new keys added to a template or base secret
// into environment-specific secrets without clobbering existing overrides.
//
// Example usage:
//
//	err := cmd.Merge(cmd.MergeOptions{
//		Address: "https://vault.example.com",
//		Token:   token,
//		Mount:   "secret",
//		SrcPath: "base/config",
//		DstPath: "prod/config",
//		DryRun:  true,
//	})
//
// Set DryRun to true to preview the merged result without writing to Vault.
package cmd
