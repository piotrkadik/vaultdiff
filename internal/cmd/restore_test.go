package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

func newRestoreStub(data map[string]string) *vault.StubClient {
	return &vault.StubClient{
		SecretData: data,
	}
}

func TestRestore_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := RestoreOptions{
		Address: "https://127.0.0.1:8200",
		Path:    "app/config",
		Version: 2,
		DryRun:  true,
		Output:  &buf,
	}

	result, err := Restore(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun=true")
	}
	if result.Data != nil {
		t.Error("expected Data to be nil on dry run")
	}
}

func TestRestoreResult_JSONShape(t *testing.T) {
	r := RestoreResult{
		Path:    "app/config",
		Version: 3,
		DryRun:  false,
		Data:    map[string]string{"key": "val"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, field := range []string{"path", "version", "dry_run", "data"} {
		if _, ok := m[field]; !ok {
			t.Errorf("missing field %q", field)
		}
	}
}

func TestDefaultRestoreOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultRestoreOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestRestore_InvalidAddress_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	opts := RestoreOptions{
		Address: "://bad",
		Path:    "app/config",
		Version: 1,
		Output:  &buf,
	}
	_, err := Restore(context.Background(), opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}
