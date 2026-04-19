// Package cmd provides CLI command implementations for vaultdiff.
//
// # DiffEnv
//
// DiffEnv compares the same secret path across two environment prefixes
// within the same Vault cluster. This is useful for auditing drift between
// staging and production, or any pair of named environments that share a
// common secret layout.
//
// # Usage
//
//	 opts := DefaultDiffEnvOptions()
//	 opts.EnvA = "staging"
//	 opts.EnvB = "production"
//	 opts.Path = "myapp/database"
//	 if err := DiffEnv(opts); err != nil {
//	     log.Fatal(err)
//	 }
//
// # Environment Variables
//
//   - VAULT_ADDR        — Vault server address
//   - VAULT_MOUNT       — KV v2 mount path (default: secret)
//   - VAULTDIFF_FORMAT  — output format: text | json | csv
//   - VAULT_MASK_VALUES — set to "false" to reveal secret values
package cmd
