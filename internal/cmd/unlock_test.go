package cmd_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultdiff/internal/cmd"
)

func newUnlockStub(data map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"data": data, "metadata": map[string]any{"version": 2}},
		})
	}))
}

func TestUnlock_DryRun_SkipsWrite(t *testing.T) {
	srv := newUnlockStub(map[string]any{"key": "val"})
	defer srv.Close()

	var buf bytes.Buffer
	err := cmd.Unlock(cmd.UnlockOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Version: 2,
		DryRun:  true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result cmd.UnlockResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !result.DryRun {
		t.Error("expected dry_run=true")
	}
	if result.Data != nil {
		t.Error("expected data omitted on dry run")
	}
}

func TestUnlock_InvalidAddress_ReturnsError(t *testing.T) {
	err := cmd.Unlock(cmd.UnlockOptions{
		Address: "://bad",
		Token:   "test",
		Path:    "myapp/config",
		Version: 1,
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestUnlockResult_JSONShape(t *testing.T) {
	srv := newUnlockStub(map[string]any{"TOKEN": "abc"})
	defer srv.Close()

	var buf bytes.Buffer
	_ = cmd.Unlock(cmd.UnlockOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Version: 1,
		DryRun:  true,
		Output:  &buf,
	})

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, field := range []string{"path", "version", "dry_run"} {
		if _, ok := out[field]; !ok {
			t.Errorf("missing field %q", field)
		}
	}
}

func TestDefaultUnlockOptions_Defaults(t *testing.T) {
	opts := cmd.DefaultUnlockOptions()
	if opts.DryRun {
		t.Error("dry_run should default to false")
	}
}
