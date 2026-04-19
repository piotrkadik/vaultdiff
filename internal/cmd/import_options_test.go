package cmd

import (
	"testing"
)

func TestDefaultImportOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://vault.test:8200")
	opts := DefaultImportOptions()
	if opts.Address != "https://vault.test:8200" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultImportOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULT_MOUNT", "kv")
	opts := DefaultImportOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", opts.Mount)
	}
}

func TestDefaultImportOptions_DryRunDefaultsFalse(t *testing.T) {
	opts := DefaultImportOptions()
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}

func TestImportOptions_ZeroValue(t *testing.T) {
	var opts ImportOptions
	if opts.DryRun {
		t.Error("zero-value DryRun should be false")
	}
	if opts.Mask {
		t.Error("zero-value Mask should be false")
	}
}
