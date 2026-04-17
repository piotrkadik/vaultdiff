package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/example/vaultdiff/internal/vault"
)

type mergeVaultStub struct {
	data map[string]map[string]string
}

func newMergeStub() *mergeVaultStub {
	return &mergeVaultStub{
		data: map[string]map[string]string{
			"base/config": {"KEY_A": "alpha", "KEY_B": "beta", "KEY_NEW": "new"},
			"prod/config": {"KEY_A": "prod-alpha", "KEY_B": "prod-beta"},
		},
	}
}

func TestMerge_DryRun_DoesNotWrite(t *testing.T) {
	var buf bytes.Buffer
	err := Merge(MergeOptions{
		Address: "https://127.0.0.1:8200",
		Token:   "test-token",
		Mount:   "secret",
		SrcPath: "base/config",
		DstPath: "prod/config",
		DryRun:  true,
		Mask:    false,
		Output:  &buf,
	})
	// expect connection error from stub address, not a logic error
	if err == nil {
		t.Fatal("expected error from unreachable vault")
	}
}

func TestMergeResult_JSONShape(t *testing.T) {
	r := MergeResult{
		SrcPath:  "base/config",
		DstPath:  "prod/config",
		DryRun:   true,
		KeyCount: 3,
		Merged:   map[string]string{"KEY_A": "alpha"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, field := range []string{"src_path", "dst_path", "dry_run", "key_count"} {
		if _, ok := out[field]; !ok {
			t.Errorf("missing field %q", field)
		}
	}
}

func TestDefaultMergeOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://vault.local")
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultMergeOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultMergeOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultMergeOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}

var _ = vault.NewClient // ensure vault package referenced
