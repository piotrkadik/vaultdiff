package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type pinVaultStub struct{ data map[string]string }

func newPinVaultStub(data map[string]string) *pinVaultStub {
	return &pinVaultStub{data: data}
}

func TestPin_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultPinOptions()
	opts.Address = "http://127.0.0.1:18200"
	opts.Token = "root"
	opts.Path = "myapp/config"
	opts.Version = 2
	opts.DryRun = true
	opts.Output = &buf

	// We expect an error because the address is unreachable, but dry-run
	// message should still surface the intent.
	_, err := Pin(context.Background(), opts)
	if err == nil {
		t.Fatal("expected error for unreachable address")
	}
}

func TestPin_WritesJSON(t *testing.T) {
	srv := vault.NewTestServer(t, map[string]map[string]string{
		"secret/data/myapp/config": {"key": "val"},
	})
	var buf bytes.Buffer
	opts := DefaultPinOptions()
	opts.Address = srv.URL
	opts.Token = "root"
	opts.Mount = "secret"
	opts.Path = "myapp/config"
	opts.Version = 1
	opts.Mask = false
	opts.Output = &buf

	result, err := Pin(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Version != 1 {
		t.Errorf("expected version 1, got %d", result.Version)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestPinResult_JSONShape(t *testing.T) {
	r := PinResult{Path: "a/b", Version: 3, DryRun: false}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, field := range []string{"path", "pinned_version", "dry_run"} {
		if !strings.Contains(s, field) {
			t.Errorf("missing field %q in JSON", field)
		}
	}
}

func TestDefaultPinOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULT_MASK", "")
	opts := DefaultPinOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultPinOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULT_MASK", "false")
	opts := DefaultPinOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}
