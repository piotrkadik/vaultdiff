package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func newFormatStub(data map[string]string) *stubVaultClient {
	return &stubVaultClient{data: data}
}

func TestFormat_TextFormat_PrintsKeyValues(t *testing.T) {
	client := newFormatStub(map[string]string{
		"API_KEY": "secret123",
		"DB_PASS": "hunter2",
	})
	var buf bytes.Buffer
	err := Format(client, FormatOptions{
		Path:    "myapp/config",
		Version: 1,
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY=secret123") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "myapp/config") {
		t.Errorf("expected path header in output, got: %s", out)
	}
}

func TestFormat_MasksValues(t *testing.T) {
	client := newFormatStub(map[string]string{"TOKEN": "topsecret"})
	var buf bytes.Buffer
	_ = Format(client, FormatOptions{
		Path:   "svc/token",
		Mask:   true,
		Output: &buf,
	})
	out := buf.String()
	if strings.Contains(out, "topsecret") {
		t.Errorf("masked output must not contain raw value")
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected mask placeholder in output")
	}
}

func TestFormat_JSONFormat_ValidOutput(t *testing.T) {
	client := newFormatStub(map[string]string{"FOO": "bar"})
	var buf bytes.Buffer
	err := Format(client, FormatOptions{
		Path:    "app/cfg",
		Version: 3,
		Mask:    false,
		Format:  "json",
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result FormatResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Path != "app/cfg" {
		t.Errorf("expected path app/cfg, got %s", result.Path)
	}
	if result.Version != 3 {
		t.Errorf("expected version 3, got %d", result.Version)
	}
}

func TestFormat_NilOutput_DefaultsToStdout(t *testing.T) {
	client := newFormatStub(map[string]string{"X": "y"})
	// Should not panic when Output is nil.
	err := Format(client, FormatOptions{Path: "p", Mask: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFormatResult_JSONShape(t *testing.T) {
	r := FormatResult{
		Path:    "x/y",
		Version: 2,
		Keys:    []string{"A", "B"},
		Data:    map[string]string{"A": "1", "B": "2"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if !strings.Contains(string(b), `"keys"`) {
		t.Errorf("expected keys field in JSON")
	}
}

func TestDefaultFormatOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultFormatOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}
