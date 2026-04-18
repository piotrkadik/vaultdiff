package cmd_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

func newKeysStub(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestKeys_TextFormat_PrintsKeys(t *testing.T) {
	srv := newKeysStub(map[string]interface{}{"alpha": "1", "beta": "2", "gamma": "3"})
	defer srv.Close()

	var buf bytes.Buffer
	err := cmd.Keys(cmd.KeysOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Version: 1,
		Format:  "text",
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, k := range []string{"alpha", "beta", "gamma"} {
		if !strings.Contains(out, k) {
			t.Errorf("expected key %q in output", k)
		}
	}
}

func TestKeys_JSONFormat_ValidOutput(t *testing.T) {
	srv := newKeysStub(map[string]interface{}{"x": "1", "y": "2"})
	defer srv.Close()

	var buf bytes.Buffer
	err := cmd.Keys(cmd.KeysOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Version: 2,
		Format:  "json",
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result cmd.KeysResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}
	if result.Path != "myapp/config" {
		t.Errorf("expected path myapp/config, got %s", result.Path)
	}
}

func TestKeys_NilOutput_DefaultsToStdout(t *testing.T) {
	srv := newKeysStub(map[string]interface{}{"k": "v"})
	defer srv.Close()

	err := cmd.Keys(cmd.KeysOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Version: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeys_InvalidAddress_ReturnsError(t *testing.T) {
	err := cmd.Keys(cmd.KeysOptions{
		Address: "://bad",
		Token:   "test",
		Path:    "myapp/config",
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}
