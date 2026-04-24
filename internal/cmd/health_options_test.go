package cmd

import (
	"testing"
)

func TestDefaultHealthOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://vault.example.com")
	opts := DefaultHealthOptions()
	if opts.Address != "https://vault.example.com" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultHealthOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULT_KV_MOUNT", "kv")
	opts := DefaultHealthOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount kv, got %q", opts.Mount)
	}
}

func TestDefaultHealthOptions_FormatDefaultsEmpty(t *testing.T) {
	t.Setenv("VAULTDIFF_FORMAT", "")
	opts := DefaultHealthOptions()
	if opts.Format != "" {
		t.Errorf("expected empty format, got %q", opts.Format)
	}
}

func TestHealthOptions_ZeroValue(t *testing.T) {
	var opts HealthOptions
	if opts.Address != "" {
		t.Error("zero-value Address should be empty")
	}
	if opts.Mount != "" {
		t.Error("zero-value Mount should be empty")
	}
	if opts.Output != nil {
		t.Error("zero-value Output should be nil")
	}
}
