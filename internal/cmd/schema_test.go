package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type schemaVaultStub struct {
	data    map[string]string
	version int
	err     error
}

func (s *schemaVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &vault.SecretVersion{Data: s.data, Version: s.version}, nil
}

func (s *schemaVaultStub) ListVersions(path string) ([]vault.VersionMeta, error) {
	return nil, nil
}

func (s *schemaVaultStub) WriteSecret(path string, data map[string]string) error {
	return nil
}

func (s *schemaVaultStub) DestroyVersion(path string, versions []int) error {
	return nil
}

func TestSchema_TextFormat_PrintsKeys(t *testing.T) {
	stub := &schemaVaultStub{data: map[string]string{"alpha": "val", "beta": ""}, version: 2}
	var buf bytes.Buffer
	err := Schema(stub, SchemaOptions{Path: "svc/cfg", Version: 2, Format: "text", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected key alpha in output")
	}
	if !strings.Contains(out, "beta (empty)") {
		t.Errorf("expected beta marked empty")
	}
}

func TestSchema_JSONFormat_ValidOutput(t *testing.T) {
	stub := &schemaVaultStub{data: map[string]string{"key1": "v"}, version: 1}
	var buf bytes.Buffer
	err := Schema(stub, SchemaOptions{Path: "svc/cfg", Version: 1, Format: "json", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result SchemaResult
	if err := json.NewDecoder(&buf).Decode(&result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result.Fields) != 1 || result.Fields[0].Key != "key1" {
		t.Errorf("unexpected fields: %+v", result.Fields)
	}
}

func TestSchema_NilOutput_DefaultsToStdout(t *testing.T) {
	stub := &schemaVaultStub{data: map[string]string{"x": "y"}, version: 1}
	err := Schema(stub, SchemaOptions{Path: "p", Version: 1, Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchema_FetchError_ReturnsError(t *testing.T) {
	stub := &schemaVaultStub{err: vault.ErrNotFound}
	var buf bytes.Buffer
	err := Schema(stub, SchemaOptions{Path: "missing", Version: 1, Output: &buf})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultSchemaOptions_FormatFromEnv(t *testing.T) {
	t.Setenv("VAULTDIFF_FORMAT", "json")
	opts := DefaultSchemaOptions()
	if opts.Format != "json" {
		t.Errorf("expected json, got %s", opts.Format)
	}
}
