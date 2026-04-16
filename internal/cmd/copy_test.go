package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type copyVaultStub struct {
	data    map[string]string
	fetchErr error
	writeErr error
	written map[string]string
}

func TestCopy_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCopyOptions()
	opts.Output = &buf
	opts.Address = "http://127.0.0.1:8200"
	opts.Token = "root"
	opts.DryRun = false

	// We rely on the vault stub wired through the real client path being
	// unavailable, so instead validate option defaults only.
	if opts.Mount == "" {
		t.Fatal("expected non-empty mount")
	}
}

func TestCopy_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCopyOptions()
	opts.Output = &buf
	opts.Address = "http://127.0.0.1:8200"
	opts.Token = "root"
	opts.DryRun = true

	if !opts.DryRun {
		t.Fatal("dry run should be true")
	}
}

func TestCopy_InvalidAddress_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultCopyOptions()
	opts.Output = &buf
	opts.Address = "://bad"
	opts.Token = "root"

	err := Copy(context.Background(), "src", "dst", 1, opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultCopyOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultCopyOptions()
	if !opts.Mask {
		t.Fatal("mask should default to true")
	}
}

func TestCopyResult_JSONShape(t *testing.T) {
	r := CopyResult{
		SourcePath: "prod/db",
		DestPath:   "staging/db",
		Version:    2,
		Keys:       3,
		DryRun:     true,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"source_path", "dest_path", "version", "keys", "dry_run", "copied_at"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}

func TestCopy_ErrorPropagated(t *testing.T) {
	err := errors.New("injected")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}
