package cmd

import (
	"os"
	"testing"
)

func TestDefaultExtractOptions_AddressFromEnv(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://vault.example.com")
	defer os.Unsetenv("VAULT_ADDR")

	opts := DefaultExtractOptions()
	if opts.Address != "http://vault.example.com" {
		t.Errorf("expected address from env, got %s", opts.Address)
	}
}

func TestDefaultExtractOptions_MountFromEnv(t *testing.T) {
	os.Setenv("VAULT_MOUNT", "kv")
	defer os.Unsetenv("VAULT_MOUNT")

	opts := DefaultExtractOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount kv, got %s", opts.Mount)
	}
}

func TestDefaultExtractOptions_MaskDefaultsTrue(t *testing.T) {
	os.Unsetenv("VAULTDIFF_MASK")

	opts := DefaultExtractOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultExtractOptions_MaskFalseWhenDisabled(t *testing.T) {
	os.Setenv("VAULTDIFF_MASK", "false")
	defer os.Unsetenv("VAULTDIFF_MASK")

	opts := DefaultExtractOptions()
	if opts.Mask {
		t.Error("expected Mask to be false when VAULTDIFF_MASK=false")
	}
}

func TestExtractOptions_ZeroValue(t *testing.T) {
	var opts ExtractOptions
	if opts.Mask {
		t.Error("zero-value Mask should be false")
	}
	if opts.Version != 0 {
		t.Error("zero-value Version should be 0")
	}
}
