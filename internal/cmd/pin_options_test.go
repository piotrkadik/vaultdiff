package cmd

import (
	"os"
	"testing"
)

func TestDefaultPinOptions_AddressFromEnv(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://vault.example.com")
	t.Cleanup(func() { os.Unsetenv("VAULT_ADDR") })
	opts := DefaultPinOptions()
	if opts.Address != "http://vault.example.com" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultPinOptions_MountFromEnv(t *testing.T) {
	os.Setenv("VAULT_MOUNT", "kv")
	t.Cleanup(func() { os.Unsetenv("VAULT_MOUNT") })
	opts := DefaultPinOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", opts.Mount)
	}
}

func TestDefaultPinOptions_DryRunDefaultsFalse(t *testing.T) {
	opts := DefaultPinOptions()
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}

func TestPinOptions_ZeroValue(t *testing.T) {
	var opts PinOptions
	if opts.Version != 0 {
		t.Error("expected zero Version")
	}
	if opts.DryRun {
		t.Error("expected DryRun false")
	}
}
