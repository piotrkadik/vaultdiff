// Package cmd provides high-level command implementations for the vaultdiff CLI.
//
// # Promote
//
// Promote copies a specific version of a Vault secret from one environment
// (source) to another (destination). It is intended for controlled promotion
// workflows where a secret that has been validated in staging needs to be
// propagated to production.
//
// Usage:
//
//	opts := cmd.DefaultPromoteOptions()
//	opts.DryRun = true   // preview without writing
//	err := cmd.Promote(ctx, srcClient, dstClient, "myapp/config", 3, "staging", "prod", opts)
//
// A JSON PromoteRecord is written to opts.Output on success, regardless of
// whether DryRun is enabled, so callers can log or audit the operation.
package cmd
