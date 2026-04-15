// Package config provides runtime configuration for the vaultdiff CLI.
//
// Configuration is sourced from environment variables and can be
// overridden programmatically by callers (e.g. the CLI flag layer).
//
// Supported environment variables:
//
//	VAULT_ADDR            – Vault server URL (required)
//	VAULT_TOKEN           – Vault authentication token (required)
//	VAULT_MOUNT           – KV v2 mount path (default: "secret")
//	VAULTDIFF_FORMAT      – Output format: "text" or "json" (default: "text")
//	VAULTDIFF_MASK        – Set to "false" to reveal secret values (default: true)
//	VAULTDIFF_SHOW_UNCHANGED – Set to "true" to include unchanged keys (default: false)
//	NO_COLOR              – Set to any non-empty value to disable ANSI colour
//
// Usage:
//
//	cfg := config.Default()
//	cfg.Mount = flagMount // CLI flag override
//	if err := cfg.Validate(); err != nil {
//	    log.Fatal(err)
//	}
package config
