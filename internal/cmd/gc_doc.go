// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # GC
//
// GC (garbage-collect) destroys old KV-v2 secret versions that exceed a
// configurable retention threshold, keeping the N most-recent non-destroyed
// versions intact.
//
// Basic usage:
//
//	err := cmd.GC(cmd.GCOptions{
//		Address: "https://vault.example.com",
//		Token:   "s.token",
//		Mount:   "secret",
//		Path:    "myapp/config",
//		Keep:    5,
//		DryRun:  true,
//	})
//
// When DryRun is true no versions are destroyed; the result JSON will contain
// the versions that *would* have been removed.
package cmd
