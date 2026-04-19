// Package cmd provides CLI command implementations for vaultdiff.
//
// # Policy
//
// Policy inspects a Vault secret and classifies each key according to
// a built-in sensitivity heuristic:
//
//   - sensitive — keys whose names contain "password", "secret", or "token"
//   - config    — keys whose names contain "url", "host", or "addr"
//   - general   — all other keys
//
// # Usage
//
//	opts := cmd.DefaultPolicyOptions()
//	opts.Path = "myapp/config"
//	if err := cmd.Policy(client, opts); err != nil {
//		log.Fatal(err)
//	}
//
// # Output Formats
//
// Text (default) prints a human-readable table. JSON emits a PolicyResult
// object suitable for machine consumption or piping into jq.
package cmd
