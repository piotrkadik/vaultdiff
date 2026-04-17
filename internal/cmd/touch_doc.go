// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Touch
//
// Touch re-writes the current version of a Vault secret at the given path
// with identical data, effectively creating a new version without changing
// any values. This is useful when you need to:
//
//   - Bump the version counter for audit or compliance purposes.
//   - Trigger downstream watchers that react to version changes.
//   - Record a "last confirmed" timestamp in the version metadata.
//
// Dry-run mode (DryRun: true) fetches and reports what would be written
// without performing the actual write.
//
// Example:
//
//	result, err := cmd.Touch("secret/myapp/config", cmd.DefaultTouchOptions())
package cmd
