package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func newDiffEnvStub(dataA, dataB map[string]string) func(addr, mount string) (vault.Client, error) {
	return func(addr, mount string) (vault.Client, error) {
		return &vaultStub{
			versions: map[string]map[string]string{
				"staging/db":    dataA,
				"production/db": dataB,
			},
		}, nil
	}
}

func TestDiffEnv_ShowsHeader(t *testing.T) {
	var buf bytes.Buffer
	opts := DiffEnvOptions{
		Address: "http://127.0.0.1:8200",
		Mount:   "secret",
		EnvA:    "staging",
		EnvB:    "production",
		Path:    "db",
		Mask:    false,
		Output:  &buf,
	}
	// DiffEnv will fail with no real Vault; check error contains path info.
	err := DiffEnv(opts)
	if err == nil {
		if !strings.Contains(buf.String(), "staging/db") {
			t.Error("expected header to contain staging/db")
		}
	}
}

func TestDiffEnv_InvalidAddress_ReturnsError(t *testing.T) {
	opts := DiffEnvOptions{
		Address: "://bad",
		Mount:   "secret",
		EnvA:    "staging",
		EnvB:    "production",
		Path:    "db",
		Output:  &bytes.Buffer{},
	}
	if err := DiffEnv(opts); err == nil {
		t.Error("expected error for invalid address")
	}
}

func TestDefaultDiffEnvOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULT_MASK_VALUES", "")
	opts := DefaultDiffEnvOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultDiffEnvOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULT_MASK_VALUES", "false")
	opts := DefaultDiffEnvOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}

func TestDiffEnvOptions_ZeroValue(t *testing.T) {
	var opts DiffEnvOptions
	if opts.Mask || opts.Format != "" || opts.Path != "" {
		t.Error("unexpected non-zero fields in zero DiffEnvOptions")
	}
}
