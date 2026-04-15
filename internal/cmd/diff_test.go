package cmd_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
	"github.com/example/vaultdiff/internal/config"
)

func newDiffVaultStub(t *testing.T, path string, dataA, dataB map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vaultStub(w, r, path, dataA, dataB)
	}))
}

func TestDiff_BasicDiff(t *testing.T) {
	dataA := map[string]interface{}{"key": "old"}
	dataB := map[string]interface{}{"key": "new"}
	srv := newDiffVaultStub(t, "secret/data/myapp", dataA, dataB)
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DiffOptions{
		Flags: cmd.Flags{Path: "myapp", VersionA: 1, VersionB: 2},
		Cfg:   config.Config{Address: srv.URL, Token: "test", Mount: "secret", Format: "text", Mask: false},
		Output: &buf,
		Audit:  nil,
	}

	code, err := cmd.Diff(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != cmd.ExitDrift {
		t.Errorf("expected ExitDrift(%d), got %d", cmd.ExitDrift, code)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestDiff_NoDrift(t *testing.T) {
	data := map[string]interface{}{"key": "same"}
	srv := newDiffVaultStub(t, "secret/data/myapp", data, data)
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DiffOptions{
		Flags: cmd.Flags{Path: "myapp", VersionA: 1, VersionB: 2},
		Cfg:   config.Config{Address: srv.URL, Token: "test", Mount: "secret", Format: "text", Mask: false},
		Output: &buf,
	}

	code, err := cmd.Diff(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != cmd.ExitOK {
		t.Errorf("expected ExitOK(%d), got %d", cmd.ExitOK, code)
	}
}

func TestDiff_InvalidAddress(t *testing.T) {
	opts := cmd.DiffOptions{
		Flags: cmd.Flags{Path: "myapp", VersionA: 1, VersionB: 2},
		Cfg:   config.Config{Address: "://bad", Token: "t", Mount: "secret"},
		Output: &bytes.Buffer{},
	}
	code, err := cmd.Diff(context.Background(), opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
	if code != cmd.ExitError {
		t.Errorf("expected ExitError, got %d", code)
	}
}
