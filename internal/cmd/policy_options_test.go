package cmd_test

import (
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

func TestDefaultPolicyOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com")
	opts := cmd.DefaultPolicyOptions()
	if opts.Address != "http://vault.example.com" {
		t.Errorf("expected address from env, got %s", opts.Address)
	}
}

func TestDefaultPolicyOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULT_MOUNT", "kv")
	opts := cmd.DefaultPolicyOptions()
	if opts.Mount != "kv" {
		t.Errorf("expected mount from env, got %s", opts.Mount)
	}
}

func TestDefaultPolicyOptions_FormatDefaultsEmpty(t *testing.T) {
	t.Setenv("VAULTDIFF_FORMAT", "")
	opts := cmd.DefaultPolicyOptions()
	if opts.Format != "" {
		t.Errorf("expected empty format, got %s", opts.Format)
	}
}

func TestPolicyOptions_ZeroValue(t *testing.T) {
	var opts cmd.PolicyOptions
	if opts.Mask {
		t.Error("zero value Mask should be false")
	}
	if opts.Path != "" {
		t.Error("zero value Path should be empty")
	}
}
