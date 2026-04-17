package cmd_test

import (
	"os"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
)

func TestDefaultCloneOptions_AddressFromEnv(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://vault:8200")
	t.Cleanup(func() { os.Unsetenv("VAULT_ADDR") })
	opts := cmd.DefaultCloneOptions()
	if opts.Address != "http://vault:8200" {
		t.Errorf("got %q, want %q", opts.Address, "http://vault:8200")
	}
}

func TestDefaultCloneOptions_MaskFalseWhenDisabled(t *testing.T) {
	os.Setenv("VAULTDIFF_MASK", "false")
	t.Cleanup(func() { os.Unsetenv("VAULTDIFF_MASK") })
	opts := cmd.DefaultCloneOptions()
	if opts.Mask {
		t.Error("expected Mask=false when VAULTDIFF_MASK=false")
	}
}

func TestDefaultCloneOptions_DryRunDefaultsFalse(t *testing.T) {
	opts := cmd.DefaultCloneOptions()
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}

func TestCloneOptions_ZeroValue(t *testing.T) {
	var opts cmd.CloneOptions
	if opts.Address != "" || opts.DryRun || opts.Mask {
		t.Error("zero value should have empty fields")
	}
}
