// Package cmd provides the command-layer wiring for vaultdiff, including flag
// parsing, output formatting, interactive prompts, pager integration, and exit
// code conventions.
//
// # Watch mode
//
// Watch mode polls a Vault KV-v2 path on a configurable interval and emits a
// unified diff to the configured writer whenever the latest version number
// advances. It is intended for long-running audit pipelines and local
// development workflows where engineers want real-time visibility into secret
// rotation.
//
// Usage:
//
//	opts := cmd.DefaultWatchOptions()
//	opts.Interval = 10 * time.Second
//	err := cmd.Watch(ctx, client, "secret/data/myapp", opts)
//
// Watch respects context cancellation, making it straightforward to integrate
// with OS signal handlers:
//
//	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
//	defer stop()
//	err := cmd.Watch(ctx, client, path, opts)
//
// When MaxRounds is set to a positive integer Watch exits cleanly after that
// many polling cycles — useful in tests and one-shot CI checks.
package cmd
