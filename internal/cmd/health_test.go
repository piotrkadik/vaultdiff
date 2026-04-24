package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func newHealthStub(data map[string]string) *vaultClientStub {
	return &vaultClientStub{data: data}
}

func TestHealth_TextFormat_OK(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultHealthOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Mount = "secret"
	opts.Output = &buf

	// An invalid address produces an unreachable result; we test the text path.
	_ = Health("myapp/config", opts)

	out := buf.String()
	if !strings.Contains(out, "vault health:") {
		t.Errorf("expected header line, got: %s", out)
	}
	if !strings.Contains(out, "address") {
		t.Errorf("expected address line, got: %s", out)
	}
}

func TestHealth_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultHealthOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Mount = "secret"
	opts.Format = "json"
	opts.Output = &buf

	_ = Health("myapp/config", opts)

	var result HealthResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Address != opts.Address {
		t.Errorf("address mismatch: got %q want %q", result.Address, opts.Address)
	}
}

func TestHealth_InvalidAddress_ReturnsUnreachable(t *testing.T) {
	var buf bytes.Buffer
	opts := HealthOptions{
		Address: "http://%%bad-address",
		Mount:   "secret",
		Format:  "json",
		Output:  &buf,
	}

	_ = Health("myapp/config", opts)

	var result HealthResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Reachable {
		t.Error("expected reachable=false for bad address")
	}
	if result.Error == "" {
		t.Error("expected non-empty error field")
	}
}

func TestHealth_NilOutput_DefaultsToStdout(t *testing.T) {
	// Should not panic when Output is nil.
	opts := HealthOptions{
		Address: "http://127.0.0.1:8200",
		Mount:   "secret",
		Output:  nil,
	}
	// We cannot capture stdout here, but we verify no panic occurs.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	_ = Health("myapp/config", opts)
}

func TestHealthResult_JSONShape(t *testing.T) {
	r := HealthResult{
		Reachable: true,
		Address:   "https://vault.example.com",
		Mount:     "secret",
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	for _, key := range []string{"reachable", "address", "mount", "checked_at"} {
		if !strings.Contains(string(b), key) {
			t.Errorf("missing key %q in JSON: %s", key, b)
		}
	}
}

func TestDefaultHealthOptions_FormatFromEnv(t *testing.T) {
	t.Setenv("VAULTDIFF_FORMAT", "json")
	opts := DefaultHealthOptions()
	if opts.Format != "json" {
		t.Errorf("expected format json, got %q", opts.Format)
	}
}
