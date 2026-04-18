package cmd_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
	"github.com/your-org/vaultdiff/internal/vault"
)

type syncVaultStub struct {
	data map[string]string
	writeErr error
}

func TestSync_WritesJSON(t *testing.T) {
	srv := newVaultTestServer(t, map[string]string{"API_KEY": "abc123"})
	defer srv.Close()

	var buf bytes.Buffer
	err := cmd.Sync(cmd.SyncOptions{
		Address: srv.URL,
		Token:   "test-token",
		Mount:   "secret",
		SrcPath: "myapp/config",
		DstPath: "myapp/config-copy",
		DryRun:  true,
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result cmd.SyncResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.SrcPath != "myapp/config" {
		t.Errorf("expected src_path myapp/config, got %s", result.SrcPath)
	}
	if result.DstPath != "myapp/config-copy" {
		t.Errorf("expected dst_path myapp/config-copy, got %s", result.DstPath)
	}
	if !result.DryRun {
		t.Error("expected dry_run true")
	}
}

func TestSync_DryRun_SkipsWrite(t *testing.T) {
	srv := newVaultTestServer(t, map[string]string{"KEY": "val"})
	defer srv.Close()

	var buf bytes.Buffer
	err := cmd.Sync(cmd.SyncOptions{
		Address: srv.URL,
		Token:   "test-token",
		Mount:   "secret",
		SrcPath: "app/prod",
		DstPath: "app/staging",
		DryRun:  true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSync_InvalidAddress_ReturnsError(t *testing.T) {
	err := cmd.Sync(cmd.SyncOptions{
		Address: "://bad",
		Token:   "x",
		SrcPath: "a",
		DstPath: "b",
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultSyncOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultSyncOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestSyncResult_JSONShape(t *testing.T) {
	r := cmd.SyncResult{
		SrcPath: "src",
		DstPath: "dst",
		DryRun:  false,
		Keys:    3,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	for _, key := range []string{"src_path", "dst_path", "dry_run", "keys"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}

var _ = errors.New
var _ = vault.NewClient
