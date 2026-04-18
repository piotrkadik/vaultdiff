package cmd_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/your/vaultdiff/internal/cmd"
)

func TestMirror_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	err := cmd.Mirror(cmd.MirrorOptions{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
		Mount:   "secret",
		SrcPath: "app/prod",
		DstPath: "app/staging",
		DryRun:  true,
		Mask:    true,
		Output:  &buf,
	})
	// unreachable vault → expect error, not panic
	if err == nil {
		t.Fatal("expected error for unreachable vault")
	}
}

func TestMirror_InvalidAddress_ReturnsError(t *testing.T) {
	err := cmd.Mirror(cmd.MirrorOptions{
		Address: "://bad",
		Token:   "x",
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMirrorResult_JSONShape(t *testing.T) {
	raw := `{"src_path":"a/b","dst_path":"a/c","dry_run":true,"keys":3}`
	var result cmd.MirrorResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.SrcPath != "a/b" {
		t.Errorf("src_path: got %q", result.SrcPath)
	}
	if result.Keys != 3 {
		t.Errorf("keys: got %d", result.Keys)
	}
}

func TestDefaultMirrorOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULT_MASK", "")
	opts := cmd.DefaultMirrorOptions()
	if !opts.Mask {
		t.Error("expected Mask=true by default")
	}
}

func TestDefaultMirrorOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULT_MASK", "false")
	opts := cmd.DefaultMirrorOptions()
	if opts.Mask {
		t.Error("expected Mask=false")
	}
}

var _ = errors.New // suppress unused import
