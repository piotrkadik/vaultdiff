package cmd

import (
	"testing"
)

func TestDefaultSchemaOptions_AddressFromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	opts := DefaultSchemaOptions()
	if opts.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected address from env, got %q", opts.Address)
	}
}

func TestDefaultSchemaOptions_MountFromEnv(t *testing.T) {
	t.Setenv("VAULT_MOUNT", "kv-eng")
	opts := DefaultSchemaOptions()
	if opts.Mount != "kv-eng" {
		t.Errorf("expected mount kv-eng, got %q", opts.Mount)
	}
}

func TestDefaultSchemaOptions_FormatDefaultsEmpty(t *testing.T) {
	t.Setenv("VAULTDIFF_FORMAT", "")
	opts := DefaultSchemaOptions()
	if opts.Format != "" {
		t.Errorf("expected empty format, got %q", opts.Format)
	}
}

func TestSchemaOptions_ZeroValue(t *testing.T) {
	var opts SchemaOptions
	if opts.Version != 0 {
		t.Errorf("expected zero version")
	}
	if opts.Output != nil {
		t.Errorf("expected nil output")
	}
}
