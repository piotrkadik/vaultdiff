// Package cmd provides the top-level command implementations for the
// vaultdiff CLI, wiring together Vault clients, diff logic, audit logging,
// and output formatting.
//
// # Baseline
//
// The baseline sub-feature captures a point-in-time snapshot of a single
// Vault secret version and persists it as a structured JSON record. The
// record can later be loaded and compared against a live secret version to
// detect drift introduced between two points in time (e.g. before and after
// a deployment).
//
// Usage:
//
//	opts := cmd.DefaultBaselineOptions()
//	opts.Version = 3          // pin to a specific KV version
//	opts.Mask = false         // emit plaintext values (use with care)
//	err := cmd.CaptureBaseline(ctx, client, "myapp/config", opts)
//
// The JSON record written to opts.Output contains:
//   - path       – the secret path
//   - version    – the KV version that was read
//   - captured_at – RFC-3339 timestamp
//   - data       – key/value pairs (optionally masked)
package cmd
