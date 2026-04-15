package config_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultdiff/internal/config"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestDefault_UsesEnvVars(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "https://vault.example.com")
	setEnv(t, "VAULT_TOKEN", "s.abc123")
	setEnv(t, "VAULT_MOUNT", "kv")

	cfg := config.Default()

	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("VaultAddr = %q, want %q", cfg.VaultAddr, "https://vault.example.com")
	}
	if cfg.VaultToken != "s.abc123" {
		t.Errorf("VaultToken = %q, want %q", cfg.VaultToken, "s.abc123")
	}
	if cfg.Mount != "kv" {
		t.Errorf("Mount = %q, want %q", cfg.Mount, "kv")
	}
}

func TestDefault_MountFallback(t *testing.T) {
	os.Unsetenv("VAULT_MOUNT")
	cfg := config.Default()
	if cfg.Mount != "secret" {
		t.Errorf("Mount = %q, want default %q", cfg.Mount, "secret")
	}
}

func TestDefault_FormatFallback(t *testing.T) {
	os.Unsetenv("VAULTDIFF_FORMAT")
	cfg := config.Default()
	if cfg.OutputFormat != "text" {
		t.Errorf("OutputFormat = %q, want %q", cfg.OutputFormat, "text")
	}
}

func TestDefault_MaskDefaultsTrue(t *testing.T) {
	os.Unsetenv("VAULTDIFF_MASK")
	cfg := config.Default()
	if !cfg.MaskValues {
		t.Error("MaskValues should default to true")
	}
}

func TestValidate_MissingAddr(t *testing.T) {
	cfg := config.Config{VaultToken: "tok", OutputFormat: "text"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing VaultAddr")
	}
}

func TestValidate_MissingToken(t *testing.T) {
	cfg := config.Config{VaultAddr: "http://localhost:8200", OutputFormat: "text"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing VaultToken")
	}
}

func TestValidate_BadFormat(t *testing.T) {
	cfg := config.Config{
		VaultAddr:    "http://localhost:8200",
		VaultToken:   "tok",
		OutputFormat: "yaml",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unsupported output format")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := config.Config{
		VaultAddr:    "http://localhost:8200",
		VaultToken:   "tok",
		OutputFormat: "json",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
