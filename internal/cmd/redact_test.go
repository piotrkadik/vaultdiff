package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type redactVaultStub struct{}

func TestRedact_DryRun_SkipsWrite(t *testing.T) {
	srv := newVaultTestServer(t, map[string]any{
		"password": "s3cr3t",
		"user":     "admin",
	})

	var buf bytes.Buffer
	result, err := Redact(RedactOptions{
		Address: srv.URL,
		Token:   "test-token",
		Path:    "secret/data/app",
		Keys:    []string{"password"},
		DryRun:  true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun=true in result")
	}
	if len(result.RedactedKeys) != 1 || result.RedactedKeys[0] != "password" {
		t.Errorf("unexpected redacted keys: %v", result.RedactedKeys)
	}
	if result.Data != nil {
		t.Error("expected Data to be nil on dry run")
	}
}

func TestRedact_WritesJSON(t *testing.T) {
	srv := newVaultTestServer(t, map[string]any{
		"api_key": "abc",
		"region":  "us-east-1",
	})

	var buf bytes.Buffer
	_, err := Redact(RedactOptions{
		Address: srv.URL,
		Token:   "test-token",
		Path:    "secret/data/app",
		Keys:    []string{"api_key"},
		DryRun:  false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if out["path"] != "secret/data/app" {
		t.Errorf("unexpected path in output: %v", out["path"])
	}
}

func TestRedact_InvalidAddress_ReturnsError(t *testing.T) {
	_, err := Redact(RedactOptions{
		Address: "://bad",
		Token:   "x",
		Path:    "secret/data/app",
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultRedactOptions_Defaults(t *testing.T) {
	opts := DefaultRedactOptions()
	if opts.DryRun {
		t.Error("DryRun should default to false")
	}
	_ = vault.SecretVersion{} // ensure vault package is used
}

func TestRedactResult_JSONShape(t *testing.T) {
	r := RedactResult{
		Path:         "a/b",
		DryRun:       true,
		RedactedKeys: []string{"x"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	if _, ok := m["redacted_keys"]; !ok {
		t.Error("expected redacted_keys field in JSON")
	}
}
