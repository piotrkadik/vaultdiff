package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestTag_WritesJSON(t *testing.T) {
	stub := newDiffVaultStub() // reuse existing httptest stub
	var buf bytes.Buffer
	opts := TagOptions{
		Address: stub.URL,
		Token:   "test-token",
		Mount:   "secret",
		Output:  &buf,
	}
	tags := map[string]string{"env": "staging"}
	err := Tag("myapp/config", 1, tags, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result TagResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Path != "myapp/config" {
		t.Errorf("expected path myapp/config, got %q", result.Path)
	}
	if result.Version != 1 {
		t.Errorf("expected version 1, got %d", result.Version)
	}
	if result.Tags["env"] != "staging" {
		t.Errorf("expected tag env=staging, got %q", result.Tags["env"])
	}
	if result.TaggedAt.IsZero() {
		t.Error("expected non-zero TaggedAt")
	}

func TestTag_InvalidAddress_ReturnsError(t *testing.T) {
	opts := TagOptions{
		Address: "http://127.0.0.1:0",
		Token:   "x",
		Mount:   "secret",
	}
	err := Tag("myapp/config", 1, nil, opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestTag_NilOutputDefaultsToStdout(t *testing.T) {
	stub := newDiffVaultStub()
	opts := TagOptions{
		Address: stub.URL,
		Token:   "test-token",
		Mount:   "secret",
		Output:  nil,
	}
	// Should not panic; stdout will receive output
	_ = Tag("myapp/config", 1, map[string]string{}, opts)
}

func TestDefaultTagOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultTagOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestTagResult_JSONShape(t *testing.T) {
	stub := newDiffVaultStub()
	var buf bytes.Buffer
	opts := TagOptions{Address: stub.URL, Token: "t", Mount: "secret", Output: &buf}
	_ = Tag("myapp/config", 1, map[string]string{"k": "v"}, opts)
	if !strings.Contains(buf.String(), "tagged_at") {
		t.Error("expected tagged_at field in JSON output")
	}
}
