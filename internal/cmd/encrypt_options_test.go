package cmd

import (
	"testing"
)

func TestDefaultEncryptOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com:8200")
	opts := DefaultEncryptOptions()
	if opts.Address != "http://vault.example.com:8200" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultEncryptOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULT_MOUNT", "kv")
	opts := DefaultEncryptOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", opts.Mount)
	}
}

func TestDefaultEncryptOptions_DryRunDefaultsFalse(t *testing.T) {
	opts := DefaultEncryptOptions()
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}

func TestEncryptOptions_ZeroValue(t *testing.T) {
	var opts EncryptOptions
	if opts.Address != "" {
		t.Error("zero-value Address should be empty")
	}
	if opts.DryRun {
		t.Error("zero-value DryRun should be false")
	}
	if opts.Mask {
		t.Error("zero-value Mask should be false")
	}
}
