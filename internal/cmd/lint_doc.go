// Package cmd provides CLI command implementations for vaultdiff.
//
// # Lint
//
// Lint fetches a secret version from Vault and checks for common key/value
// issues such as empty values or well-known placeholder strings ("changeme",
// "TODO", "FIXME").
//
// Each issue is reported with a severity level:
//
//   - warning — the value is present but empty.
//   - error   — the value is a known placeholder that should be replaced.
//
// Output formats:
//
//   - text (default): one line per issue, prefixed with severity.
//   - json: a JSON array of LintResult objects.
//
// Environment variables:
//
//	VAULT_ADDR   — Vault server address
//	VAULT_TOKEN  — Vault token
//	VAULT_MOUNT  — KV v2 mount path (default: "secret")
//	VAULTDIFF_FORMAT — output format ("text" or "json")
package cmd
