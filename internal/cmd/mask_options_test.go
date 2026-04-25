package cmd

import (
	"testing"
)

func TestDefaultMaskOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com")
	opts := DefaultMaskOptions()
	if opts.Address != "http://vault.example.com" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultMaskOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULT_MOUNT", "kv")
	opts := DefaultMaskOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount from env, got %q", opts.Mount)
	}
}

func TestDefaultMaskOptions_DryRunDefaultsFalse(t *testing.T) {
	opts := DefaultMaskOptions()
	if opts.DryRun {
		t.Error("expected DryRun=false by default")
	}
}

func TestMaskOptions_ZeroValue(t *testing.T) {
	var opts MaskOptions
	if opts.DryRun {
		t.Error("zero-value DryRun should be false")
	}
	if opts.Mask {
		t.Error("zero-value Mask should be false")
	}
	if opts.Output != nil {
		t.Error("zero-value Output should be nil")
	}
}
