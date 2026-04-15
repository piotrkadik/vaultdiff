// Package config loads and validates vaultdiff runtime configuration
// from environment variables and CLI flag overrides.
package config

import (
	"errors"
	"os"
)

// Config holds all runtime settings for vaultdiff.
type Config struct {
	// VaultAddr is the base URL of the Vault server.
	VaultAddr string

	// VaultToken is the authentication token.
	VaultToken string

	// Mount is the KV v2 secrets engine mount path.
	Mount string

	// OutputFormat controls diff output: "text" or "json".
	OutputFormat string

	// MaskValues hides secret values in output when true.
	MaskValues bool

	// ShowUnchanged includes unchanged keys in diff output.
	ShowUnchanged bool

	// Color enables ANSI colour in text output.
	Color bool
}

// Default returns a Config populated from environment variables.
// Callers may override individual fields after calling Default.
func Default() Config {
	mount := os.Getenv("VAULT_MOUNT")
	if mount == "" {
		mount = "secret"
	}

	fmt := os.Getenv("VAULTDIFF_FORMAT")
	if fmt == "" {
		fmt = "text"
	}

	return Config{
		VaultAddr:     os.Getenv("VAULT_ADDR"),
		VaultToken:    os.Getenv("VAULT_TOKEN"),
		Mount:         mount,
		OutputFormat:  fmt,
		MaskValues:    os.Getenv("VAULTDIFF_MASK") != "false",
		ShowUnchanged: os.Getenv("VAULTDIFF_SHOW_UNCHANGED") == "true",
		Color:         os.Getenv("NO_COLOR") == "",
	}
}

// Validate returns an error if required fields are missing or invalid.
func (c Config) Validate() error {
	if c.VaultAddr == "" {
		return errors.New("config: VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		return errors.New("config: VAULT_TOKEN is required")
	}
	if c.OutputFormat != "text" && c.OutputFormat != "json" {
		return errors.New("config: output format must be \"text\" or \"json\"")
	}
	return nil
}
