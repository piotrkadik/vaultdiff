// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Search
//
// Search fetches a single secret version from Vault and filters its keys
// by a case-insensitive substring query. It is useful for quickly locating
// a specific configuration key across a large secret without performing a
// full two-version diff.
//
// Basic usage:
//
//	changes, err := cmd.Search(client, "secret/myapp", 3, cmd.SearchOptions{
//		Query:  "database",
//		Mask:   true,
//		Output: os.Stdout,
//	})
//
// Format may be "text" (default) or "json".
// When Mask is true, values are replaced with "***" in the output.
package cmd
