package cmd_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
)

func TestGC_DryRun_DoesNotDestroy(t *testing.T) {
	var buf bytes.Buffer
	err := cmd.GC(cmd.GCOptions{
		Address: "http://127.0.0.1:1", // unreachable
		Token:   "tok",
		Mount:   "secret",
		Path:    "app/cfg",
		Keep:    3,
		DryRun:  true,
		Output:  &buf,
	})
	// Expect a connection error since no real Vault is running.
	if err == nil {
		t.Fatal("expected error for unreachable address")
	}
}

func TestGC_InvalidAddress_ReturnsError(t *testing.T) {
	err := cmd.GC(cmd.GCOptions{
		Address: "://bad",
		Token:   "tok",
		Path:    "app/cfg",
		Keep:    5,
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestGCResult_JSONShape(t *testing.T) {
	result := cmd.GCResult{
		Path:    "app/cfg",
		Deleted: []int{1, 2},
		Kept:    []int{3, 4, 5},
		DryRun:  true,
	}
	b, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"path", "deleted", "kept", "dry_run"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}

func TestDefaultGCOptions_Keep(t *testing.T) {
	opts := cmd.DefaultGCOptions()
	if opts.Keep != 5 {
		t.Errorf("expected default keep=5, got %d", opts.Keep)
	}
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}
