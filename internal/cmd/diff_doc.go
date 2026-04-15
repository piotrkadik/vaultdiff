// Package cmd provides the command-layer for vaultdiff.
//
// # Diff
//
// The Diff function is the primary entry-point for a single diff run.
// It wires together the Vault client, the diff engine, output formatting,
// and optional audit logging into one coherent operation.
//
// Usage:
//
//	flags, err := cmd.ParseFlags(os.Args[1:])
//	if err != nil {
//	    log.Fatal(err)
//	}
//	opts := cmd.DefaultDiffOptions(flags)
//	code, err := cmd.Diff(ctx, opts)
//	if err != nil {
//	    fmt.Fprintln(os.Stderr, err)
//	}
//	os.Exit(code)
//
// DiffOptions allows callers to inject a custom io.Writer for output and
// an *audit.Logger for structured audit trails. When either field is nil
// sensible defaults (os.Stdout / stderr) are used automatically.
package cmd
