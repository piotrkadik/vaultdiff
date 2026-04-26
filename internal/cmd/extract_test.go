package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newExtractStub(data map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]interface{}{
			"data": map[string]interface{}{
				"data":     data,
				"metadata": map[string]interface{}{"version": 3},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	}))
}

func TestExtract_AllKeys_TextFormat(t *testing.T) {
	srv := newExtractStub(map[string]interface{}{"alpha": "1", "beta": "2"})
	defer srv.Close()

	var buf bytes.Buffer
	_, err := Extract(ExtractOptions{
		Address: srv.URL,
		Token:   "test",
		Path:    "myapp/config",
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alpha=1") {
		t.Errorf("expected alpha=1 in output, got: %s", out)
	}
}

func TestExtract_SubsetKeys_ReturnsOnlyRequested(t *testing.T) {
	srv := newExtractStub(map[string]interface{}{"alpha": "1", "beta": "2", "gamma": "3"})
	defer srv.Close()

	var buf bytes.Buffer
	res, err := Extract(ExtractOptions{
		Address: srv.URL,
		Token:   "test",
		Path:    "myapp/config",
		Keys:    []string{"alpha", "gamma"},
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Data["beta"]; ok {
		t.Error("beta should not be present in extracted result")
	}
	if res.Data["alpha"] != "1" {
		t.Errorf("expected alpha=1, got %s", res.Data["alpha"])
	}
}

func TestExtract_MissingKey_ReturnsError(t *testing.T) {
	srv := newExtractStub(map[string]interface{}{"alpha": "1"})
	defer srv.Close()

	_, err := Extract(ExtractOptions{
		Address: srv.URL,
		Token:   "test",
		Path:    "myapp/config",
		Keys:    []string{"missing"},
		Mask:    false,
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestExtract_MasksValues(t *testing.T) {
	srv := newExtractStub(map[string]interface{}{"secret": "s3cr3t"})
	defer srv.Close()

	var buf bytes.Buffer
	_, err := Extract(ExtractOptions{
		Address: srv.URL,
		Token:   "test",
		Path:    "myapp/config",
		Mask:    true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "s3cr3t") {
		t.Error("raw value should be masked in output")
	}
}

func TestExtract_JSONFormat_ValidOutput(t *testing.T) {
	srv := newExtractStub(map[string]interface{}{"key": "val"})
	defer srv.Close()

	var buf bytes.Buffer
	_, err := Extract(ExtractOptions{
		Address: srv.URL,
		Token:   "test",
		Path:    "myapp/config",
		Format:  "json",
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result ExtractResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.Path != "myapp/config" {
		t.Errorf("expected path myapp/config, got %s", result.Path)
	}
}

func TestExtract_NilOutput_DefaultsToStdout(t *testing.T) {
	srv := newExtractStub(map[string]interface{}{"k": "v"})
	defer srv.Close()

	_, err := Extract(ExtractOptions{
		Address: srv.URL,
		Token:   "test",
		Path:    "myapp/config",
		Mask:    false,
		Output:  nil,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
