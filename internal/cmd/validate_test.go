package cmd_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

func newValidateStub(data map[string]string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		kv := ""
		for k, v := range data {
			if kv != "" {
				kv += ","
			}
			kv += `"` + k + `":"` + v + `"`
		}
		fmt.Fprintf(w, `{"data":{"data":{%s},"metadata":{"version":1}}}`, kv)
	}))
}

func TestValidate_AllKeysPresent(t *testing.T) {
	srv := newValidateStub(map[string]string{"DB_PASS": "secret", "API_KEY": "abc"}, http.StatusOK)
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DefaultValidateOptions()
	opts.Address = srv.URL
	opts.Token = "test"
	opts.Output = &buf
	opts.Mask = false

	err := cmd.Validate("myapp/config", 1, []string{"DB_PASS", "API_KEY"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "All required keys present") {
		t.Errorf("expected success message, got: %s", buf.String())
	}
}

func TestValidate_MissingKey_ReturnsError(t *testing.T) {
	srv := newValidateStub(map[string]string{"DB_PASS": "secret"}, http.StatusOK)
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DefaultValidateOptions()
	opts.Address = srv.URL
	opts.Token = "test"
	opts.Output = &buf

	err := cmd.Validate("myapp/config", 1, []string{"DB_PASS", "MISSING_KEY"}, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "MISSING_KEY") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestValidate_MasksValues(t *testing.T) {
	srv := newValidateStub(map[string]string{"TOKEN": "super-secret"}, http.StatusOK)
	defer srv.Close()

	var buf bytes.Buffer
	opts := cmd.DefaultValidateOptions()
	opts.Address = srv.URL
	opts.Token = "test"
	opts.Output = &buf
	opts.Mask = true

	_ = cmd.Validate("myapp/config", 1, []string{"TOKEN"}, opts)
	if strings.Contains(buf.String(), "super-secret") {
		t.Error("masked output should not contain plaintext value")
	}
}

func TestValidate_InvalidAddress_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultValidateOptions()
	opts.Address = "://bad-address"
	opts.Output = &buf

	err := cmd.Validate("myapp/config", 1, []string{"KEY"}, opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultValidateOptions_MaskDefaultsTrue(t *testing.T) {
	opts := cmd.DefaultValidateOptions()
	if !opts.Mask {
		t.Error("Mask should default to true")
	}
}
