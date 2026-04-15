// Package cmd provides high-level command implementations for the vaultdiff
// CLI tool.
//
// # Rollback
//
// The Rollback function retrieves a specific historical version of a Vault KV
// secret and emits a JSON record describing the intended or applied change.
//
// Typical usage:
//
//	opts := cmd.DefaultRollbackOptions()
//	opts.DryRun = true  // preview only
//	err := cmd.Rollback(ctx, client, "secret/myapp", 5, 3, opts)
//
// When DryRun is false the record is still written to Output so that callers
// can pipe it into an audit log or CI artefact store. The caller is responsible
// for actually applying the data back to Vault via the appropriate write API.
//
// Fields in RollbackRecord:
//
//	path          - KV path of the secret
//	from_version  - the version that was active before rollback
//	to_version    - the version being restored
//	data          - key/value pairs from the target version
//	timestamp     - UTC time the rollback record was produced
//	dry_run       - true when no write was performed
package cmd
