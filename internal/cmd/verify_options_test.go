package cmd

import (
	"testing"
)

func TestDefaultVerifyOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com")
	opts := DefaultVerifyOptions()
	if opts.Address != "http://vault.example.com" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultVerifyOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULTDIFF_MOUNT", "kv")
	opts := DefaultVerifyOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount kv, got %q", opts.Mount)
	}
}

func TestDefaultVerifyOptions_FormatDefaultsEmpty(t *testing.T) {
	t.Setenv("VAULTDIFF_FORMAT", "")
	opts := DefaultVerifyOptions()
	if opts.Format != "" {
		t.Errorf("expected empty format, got %q", opts.Format)
	}
}

func TestVerifyOptions_ZeroValue(t *testing.T) {
	var opts VerifyOptions
	if opts.Address != "" || opts.Mount != "" || opts.Mask || opts.Format != "" {
		t.Error("zero value should have empty fields and Mask=false")
	}
}
