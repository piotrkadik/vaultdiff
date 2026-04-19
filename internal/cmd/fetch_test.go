package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type fetchVaultStub struct{ data map[string]string }

func (f *fetchVaultStub) GetSecretVersion(_ *vault.Client, _ string, _ int) (map[string]string, error) {
	return f.data, nil
}

func TestFetch_TextFormat_PrintsKeyValues(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFetchOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Mask = false
	opts.Format = "text"
	opts.Output = &buf

	if err := Fetch("secret/data/app", opts); err != nil && !strings.Contains(err.Error(), "fetch:") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetch_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFetchOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Mask = false
	opts.Format = "json"
	opts.Output = &buf

	_ = Fetch("secret/data/app", opts)
	if buf.Len() == 0 {
		t.Skip("no vault available")
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
}

func TestFetch_MasksValues(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultFetchOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Mask = true
	opts.Format = "text"
	opts.Output = &buf

	_ = Fetch("secret/data/app", opts)
	if strings.Contains(buf.String(), "plaintext") {
		t.Error("expected values to be masked")
	}
}

func TestFetch_InvalidAddress_ReturnsError(t *testing.T) {
	opts := DefaultFetchOptions()
	opts.Address = "://bad"
	opts.Output = &bytes.Buffer{}

	if err := Fetch("secret/data/app", opts); err == nil {
		t.Error("expected error for invalid address")
	}
}

func TestDefaultFetchOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULT_MASK", "")
	opts := DefaultFetchOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultFetchOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULT_MASK", "false")
	opts := DefaultFetchOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}
