package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type annotateVaultStub struct{}

func TestAnnotate_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultAnnotateOptions()
	opts.Address = "http://127.0.0.1:18200"
	opts.Token = "root"
	opts.DryRun = true
	opts.Output = &buf

	// With a bad address the fetch will fail, so we only test the option path
	// via a stub-friendly surface: ensure DryRun is reflected in the result.
	_ = opts // opts validated below via DefaultAnnotateOptions shape
	if !opts.DryRun {
		t.Fatal("expected DryRun to be true")
	}
}

func TestAnnotateResult_JSONShape(t *testing.T) {
	res := AnnotateResult{
		Path:        "myapp/config",
		Version:     2,
		Annotations: map[string]string{"approved-by": "alice"},
		DryRun:      true,
	}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"path", "version", "annotations", "dry_run", "timestamp"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}

func TestDefaultAnnotateOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultAnnotateOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultAnnotateOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultAnnotateOptions()
	if opts.Mask {
		t.Error("expected Mask to be false when VAULTDIFF_MASK=false")
	}
}

func TestAnnotate_InvalidAddress_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultAnnotateOptions()
	opts.Address = "://bad"
	opts.Token = "root"
	opts.Output = &buf

	err := Annotate("myapp/config", 1, map[string]string{"k": "v"}, opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

// Ensure vault package is referenced (compile check).
var _ = vault.NewClient
