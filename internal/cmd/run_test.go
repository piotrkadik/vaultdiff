package cmd_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
	"github.com/example/vaultdiff/internal/config"
)

func vaultStub(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/secret/data/myapp", func(w http.ResponseWriter, r *http.Request) {
		v := r.URL.Query().Get("version")
		switch v {
		case "1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"data":{"KEY":"old"},"metadata":{"version":1}}}`))
		case "2":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"data":{"KEY":"new"},"metadata":{"version":2}}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	return httptest.NewServer(mux)
}

func TestRun_BasicDiff(t *testing.T) {
	srv := vaultStub(t)
	defer srv.Close()

	cfg := config.Config{
		Address: srv.URL,
		Token:   "test-token",
		Mount:   "secret",
	}
	opts := cmd.RunOptions{
		Path:     "myapp",
		VersionA: 1,
		VersionB: 2,
		Mask:     false,
	}

	var buf bytes.Buffer
	err := cmd.Run(context.Background(), opts, cfg, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY") {
		t.Errorf("expected output to contain KEY, got: %s", buf.String())
	}
}

func TestRun_InvalidAddress(t *testing.T) {
	cfg := config.Config{
		Address: "://bad",
		Token:   "tok",
		Mount:   "secret",
	}
	err := cmd.Run(context.Background(), cmd.RunOptions{Path: "x", VersionA: 1, VersionB: 2}, cfg, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for bad address")
	}
}
