// Package cmd provides CLI command implementations for vaultdiff.
//
// # Normalize
//
// Normalize reads a KV v2 secret at a given path and version, applies
// cosmetic cleanup rules, and writes the result back to Vault.
//
// Cleanup rules applied:
//   - All key names are lowercased.
//   - Leading and trailing whitespace is stripped from both keys and values.
//
// The operation is idempotent: running Normalize twice on the same secret
// produces the same output.
//
// # Dry-run mode
//
// When DryRun is true the cleaned key/value pairs are printed to the output
// writer but no write is issued to Vault. This is useful for previewing
// changes before committing them.
//
// # Output formats
//
// Normalize supports two output formats controlled by NormalizeOptions.Format:
//   - "" or "text" — human-readable key = value lines (default).
//   - "json"       — a JSON-encoded NormalizeResult object.
package cmd
