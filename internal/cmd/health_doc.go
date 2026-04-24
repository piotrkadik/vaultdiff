// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Health
//
// Health performs a lightweight connectivity and mount-accessibility check
// against a HashiCorp Vault instance.  It attempts to read the supplied
// secret path at version 0 (latest) and reports whether the operation
// succeeds.
//
// # Output formats
//
// Text (default) – human-readable status lines:
//
//	vault health: OK
//	  address : https://vault.example.com
//	  mount   : secret
//
// JSON – machine-readable HealthResult object:
//
//	{
//	  "reachable": true,
//	  "address":   "https://vault.example.com",
//	  "mount":     "secret",
//	  "checked_at": "2024-01-15T12:00:00Z"
//	}
//
// # Environment variables
//
//	VAULT_ADDR        – Vault server address
//	VAULT_KV_MOUNT    – KV v2 mount path (default: "secret")
//	VAULTDIFF_FORMAT  – Output format: "" or "json"
package cmd
