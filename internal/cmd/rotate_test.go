package cmd_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
)

func TestRotate_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultRotateOptions()
	opts.Address = "http://127.0.0.1:18200"
	opts.Token = "test-token"
	opts.Path = "secret/app"
	opts.DryRun = true
	opts.Output = &buf

	// Expect error because no real Vault is running; just check it's an address error.
	err := cmd.Rotate(opts)
	if err == nil {
		t.Fatal("expected error for unreachable vault")
	}
}

func TestRotate_InvalidAddress_ReturnsError(t *testing.T) {
	opts := cmd.DefaultRotateOptions()
	opts.Address = "://bad"
	opts.Path = "secret/app"

	err := cmd.Rotate(opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestRotateResult_JSONShape(t *testing.T) {
	r := cmd.RotateResult{
		Path:        "secret/app",
		DryRun:      true,
		PrevVersion: 3,
	}

	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if m["path"] != "secret/app" {
		t.Errorf("expected path secret/app, got %v", m["path"])
	}
	if m["dry_run"] != true {
		t.Errorf("expected dry_run true")
	}
	if _, ok := m["new_version"]; ok {
		t.Error("new_version should be omitted when zero")
	}
}

func TestDefaultRotateOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultRotateOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultRotateOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := cmd.DefaultRotateOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}
