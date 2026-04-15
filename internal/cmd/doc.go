// Package cmd wires together the vaultdiff command-line interface.
//
// # Subpackages and responsibilities
//
//   - flags.go     – parses and validates CLI flags into a typed Args struct.
//   - run.go       – orchestrates the full diff workflow: fetch → compare → report → output.
//   - output.go    – selects the correct formatter (text / JSON / CSV) and writes results.
//   - pager.go     – optionally pipes output through a system pager (e.g. less).
//   - interactive.go – interactive drift-confirmation prompts for CI gate workflows.
//   - exit.go      – canonical exit-code constants and helpers.
//
// # Exit codes
//
//   0  No drift detected (or --no-exit-code flag supplied).
//   1  Drift detected between the two secret versions.
//   2  Fatal error; the diff could not be completed.
//   3  User cancelled an interactive confirmation prompt.
package cmd
