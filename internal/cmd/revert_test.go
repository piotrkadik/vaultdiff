package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/someone/vaultdiff/internal/cmd"
)

func newRevertStub(data map[string]any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"data": data},
			})
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
}

func TestRevert_DryRun_SkipsWrite(t *testing.T) {
	srv := newRevertStub(map[string]any{"key": "value"})
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DefaultRevertOptions()
	opts.Address = srv.URL
	opts.Token = "test"
	opts.Path = "secret/app"
	opts.Version = 2
	opts.DryRun = true
	opts.Mask = false
	opts.Output = &buf

	result, err := cmd.Revert(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun=true")
	}
	if result.Data != nil {
		t.Error("expected no data on dry run")
	}
}

func TestRevertResult_JSONShape(t *testing.T) {
	srv := newRevertStub(map[string]any{"token": "abc"})
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DefaultRevertOptions()
	opts.Address = srv.URL
	opts.Token = "test"
	opts.Path = "secret/app"
	opts.Version = 1
	opts.Mask = false
	opts.Output = &buf

	_, err := cmd.Revert(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["path"]; !ok {
		t.Error("expected 'path' field in output")
	}
	if _, ok := out["version"]; !ok {
		t.Error("expected 'version' field in output")
	}
}

func TestDefaultRevertOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultRevertOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestRevert_InvalidAddress_ReturnsError(t *testing.T) {
	opts := cmd.DefaultRevertOptions()
	opts.Address = "://bad"
	opts.Token = "test"
	opts.Path = "secret/app"
	opts.Version = 1

	_, err := cmd.Revert(context.Background(), opts)
	if err == nil {
		t.Error("expected error for invalid address")
	}
}
