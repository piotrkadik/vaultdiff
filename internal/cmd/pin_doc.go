// Package cmd provides CLI command implementations for vaultdiff.
//
// # Pin
//
// Pin records a specific version of a Vault secret as the canonical pinned
// reference. This is useful for locking a known-good secret version before
// performing environment promotions or rollbacks.
//
// # Usage
//
//	opts := cmd.DefaultPinOptions()
//	opts.Path    = "myapp/config"
//	opts.Version = 3
//	result, err := cmd.Pin(ctx, opts)
//
// # Dry Run
//
// When DryRun is true, no output is written to Vault and only a notice is
// printed. The returned PinResult will have DryRun set to true.
package cmd
