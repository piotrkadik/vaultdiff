// Package cmd implements the top-level command logic for vaultdiff.
//
// It wires together the vault client, diff engine, report renderer, and audit
// logger into a single Run function that can be driven by the main binary or
// called directly in tests.
//
// Typical usage:
//
//	cfg := config.Default()
//	opts, err := cmd.ParseFlags(os.Args[1:], os.Stderr)
//	if err != nil {
//		os.Exit(2)
//	}
//	if err := cmd.Run(context.Background(), opts, cfg, os.Stdout); err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
package cmd
