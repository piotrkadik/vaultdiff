package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newEncryptStub(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "metadata") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"versions": map[string]interface{}{}},
			})
			return
		}
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{"version": 2},
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{"data": data},
		})
	}))
}

func TestEncrypt_DryRun_SkipsWrite(t *testing.T) {
	srv := newEncryptStub(map[string]interface{}{"password": "secret"})
	defer srv.Close()

	var buf bytes.Buffer
	res, err := Encrypt(EncryptOptions{
		Address: srv.URL,
		Token:   "test",
		Mount:   "secret",
		Path:    "app/config",
		Keys:    []string{"password"},
		DryRun:  true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if res.Version != 0 {
		t.Errorf("expected version 0 for dry run, got %d", res.Version)
	}
}

func TestEncryptResult_JSONShape(t *testing.T) {
	res := EncryptResult{Path: "app/cfg", Version: 3, Keys: []string{"k"}, DryRun: false}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"path"`) {
		t.Error("expected 'path' in JSON output")
	}
	if !strings.Contains(string(b), `"version"`) {
		t.Error("expected 'version' in JSON output")
	}
}

func TestEncrypt_InvalidAddress_ReturnsError(t *testing.T) {
	_, err := Encrypt(EncryptOptions{
		Address: "://bad",
		Token:   "x",
		Path:    "app/cfg",
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultEncryptOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://localhost:8200")
	t.Setenv("VAULT_TOKEN", "tok")
	t.Unsetenv("VAULTDIFF_MASK")
	opts := DefaultEncryptOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultEncryptOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultEncryptOptions()
	if opts.Mask {
		t.Error("expected Mask=false when env var is 'false'")
	}
}
