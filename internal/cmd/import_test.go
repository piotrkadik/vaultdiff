package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestImport_DryRun_SkipsWrite(t *testing.T) {
	payload := `{"KEY":"value","OTHER":"data"}`
	opts := DefaultImportOptions()
	opts.DryRun = true
	opts.Path = "app/config"
	opts.Output = &bytes.Buffer{}

	result, err := Import(strings.NewReader(payload), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.DryRun {
		t.Error("expected DryRun=true")
	}
	if result.Keys != 2 {
		t.Errorf("expected 2 keys, got %d", result.Keys)
	}
	if result.Data != nil {
		t.Error("expected Data to be nil in dry-run mode")
	}
}

func TestImport_InvalidJSON_ReturnsError(t *testing.T) {
	opts := DefaultImportOptions()
	opts.DryRun = true
	_, err := Import(strings.NewReader(`not json`), opts)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestImportResult_JSONShape(t *testing.T) {
	r := ImportResult{
		Path:   "app/cfg",
		Keys:   3,
		DryRun: false,
		Data:   map[string]any{"A": "***"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{`"path"`, `"keys"`, `"dry_run"`, `"data"`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing field %s in JSON output", want)
		}
	}
}

func TestDefaultImportOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultImportOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultImportOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultImportOptions()
	if opts.Mask {
		t.Error("expected Mask=false when VAULTDIFF_MASK=false")
	}
}

func TestImport_NilOutput_DefaultsToStdout(t *testing.T) {
	opts := DefaultImportOptions()
	opts.DryRun = true
	opts.Output = nil
	_, err := Import(strings.NewReader(`{"K":"v"}`), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
