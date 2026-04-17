// Package cmd provides high-level command implementations for the vaultdiff
// CLI tool.
//
// # Annotate
//
// Annotate attaches arbitrary key/value metadata to a specific version of a
// Vault KV-v2 secret. Annotations are stored inline as reserved keys prefixed
// with "_annotation." so they remain visible in diffs and audit logs without
// requiring a separate metadata store.
//
// Annotations are written as a new secret version, preserving the existing
// data. Pass DryRun: true to preview the result without committing.
//
// Example:
//
//	err := cmd.Annotate("myapp/config", 3, map[string]string{
//		"approved-by": "alice",
//		"ticket":      "OPS-1234",
//	}, cmd.DefaultAnnotateOptions())
package cmd
