// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # Scrub
//
// Scrub zeroes out one or more keys within a Vault KV-v2 secret, writing the
// result as a new secret version.  This is useful for removing sensitive values
// (e.g. credentials that have been rotated) while preserving the rest of the
// secret's structure and history.
//
// # Usage
//
//	result, err := cmd.Scrub(client, cmd.ScrubOptions{
//		Path:    "secret/data/myapp",
//		Version: 4,
//		Keys:    []string{"api_key", "db_password"},
//		DryRun:  false,
//	})
//
// When DryRun is true, no write is performed and the output is prefixed with a
// "[dry-run]" notice so callers can preview which keys would be scrubbed.
//
// # Output
//
// Scrub always emits a JSON object containing:
//
//	{
//	  "path":     "secret/data/myapp",
//	  "version":  5,
//	  "scrubbed": ["api_key", "db_password"],
//	  "dry_run":  false
//	}
package cmd
