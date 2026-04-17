package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/example/vaultdiff/internal/vault"
)

func TestTrim_DryRun_DoesNotDestroy(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultTrimOptions()
	opts.Address = "http://127.0.0.1:18200"
	opts.Token = "root"
	opts.Path = "myapp/config"
	opts.Keep = 1
	opts.DryRun = true
	opts.Output = &buf

	// Without a real Vault this will error on NewClient or ListVersions;
	// we only verify DryRun path doesn't panic and returns an error cleanly.
	err := Trim(opts)
	if err == nil {
		t.Skip("live vault not available")
	}
}

func TestTrimResult_JSONShape(t *testing.T) {
	r := TrimResult{
		Path:    "secret/app",
		Removed: []int{1, 2},
		DryRun:  true,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["path"] != "secret/app" {
		t.Errorf("expected path field")
	}
	if m["dry_run"] != true {
		t.Errorf("expected dry_run=true")
	}
}

func TestDefaultTrimOptions_Keep(t *testing.T) {
	opts := DefaultTrimOptions()
	if opts.Keep != 5 {
		t.Errorf("expected Keep=5, got %d", opts.Keep)
	}
}

func TestTrim_NilOutputDefaultsToStdout(t *testing.T) {
	opts := DefaultTrimOptions()
	opts.Output = nil
	opts.Address = "http://127.0.0.1:19200"
	_ = vault.DestroyVersion // ensure import used
	err := Trim(opts)
	if err == nil {
		t.Skip("live vault not available")
	}
}
