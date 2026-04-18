package cmd_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

func newProtectStub(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"data": data},
			})
		case r.Method == http.MethodPost || r.Method == http.MethodPut:
			w.WriteHeader(http.StatusOK)
		}
	}))
}

func TestProtect_DryRun_SkipsWrite(t *testing.T) {
	srv := newProtectStub(map[string]interface{}{"API_KEY": "secret", "DB_PASS": "pass"})
	defer srv.Close()

	var buf bytes.Buffer
	err := cmd.Protect(cmd.ProtectOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Keys:    []string{"API_KEY"},
		DryRun:  true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["dry_run"] != true {
		t.Errorf("expected dry_run=true, got %v", result["dry_run"])
	}
}

func TestProtectResult_JSONShape(t *testing.T) {
	srv := newProtectStub(map[string]interface{}{"TOKEN": "abc"})
	defer srv.Close()

	var buf bytes.Buffer
	_ = cmd.Protect(cmd.ProtectOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Keys:    []string{"TOKEN"},
		DryRun:  false,
		Output:  &buf,
	})

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, field := range []string{"path", "protected_keys", "dry_run"} {
		if _, ok := result[field]; !ok {
			t.Errorf("missing field %q in result", field)
		}
	}
}

func TestProtect_MissingKey_ReturnsError(t *testing.T) {
	srv := newProtectStub(map[string]interface{}{"ONLY_KEY": "val"})
	defer srv.Close()

	err := cmd.Protect(cmd.ProtectOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "myapp/config",
		Keys:    []string{"MISSING_KEY"},
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestDefaultProtectOptions_Defaults(t *testing.T) {
	opts := cmd.DefaultProtectOptions()
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}
