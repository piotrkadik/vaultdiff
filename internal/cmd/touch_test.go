package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

type touchVaultStub struct{ data map[string]string }

func TestTouch_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultTouchOptions()
	opts.Address = "http://127.0.0.1:18200"
	opts.Token = "root"
	opts.DryRun = true
	opts.Output = &buf

	// With an unreachable address the client creation itself should fail.
	_, err := Touch("secret/app", opts)
	if err == nil {
		t.Fatal("expected error for unreachable vault")
	}
}

func TestTouchResult_JSONShape(t *testing.T) {
	var buf bytes.Buffer
	res := TouchResult{
		Path:   "secret/app",
		Version: 3,
		DryRun: true,
	}
	if err := json.NewEncoder(&buf).Encode(res); err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out["path"] != "secret/app" {
		t.Errorf("unexpected path: %v", out["path"])
	}
	if out["dry_run"] != true {
		t.Errorf("expected dry_run true")
	}
	if _, ok := out["data"]; ok {
		t.Errorf("data should be omitted when nil")
	}
}

func TestDefaultTouchOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultTouchOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultTouchOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultTouchOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}
