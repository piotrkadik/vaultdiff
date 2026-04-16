// Package cmd provides CLI command implementations for vaultdiff.
//
// # Tag
//
// Tag attaches arbitrary string metadata (tags) to a specific version of a
// Vault secret path. Tags are recorded in the output JSON and can be used for
// auditing, promotion tracking, or environment labelling.
//
// Basic usage:
//
//	err := cmd.Tag("secret/myapp", 3, map[string]string{
//		"env":       "production",
//		"promoted-by": "ci",
//	}, cmd.DefaultTagOptions())
//
// The command verifies the requested version exists before recording the tag,
// and writes a JSON TagResult to the configured output writer.
package cmd
