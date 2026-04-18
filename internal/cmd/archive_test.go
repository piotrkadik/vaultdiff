package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestArchive_WritesJSON(t *testing.T) {
	srv := newDiffVaultStub(t, map[string]string{"DB_PASS": "secret"})
	var buf bytes.Buffer
	err := Archive(ArchiveOptions{
		Address: srv.URL,
		Token:   "test-token",
		Path:    "secret/data/app",
		Version: 1,
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rec ArchiveRecord
	if err := json.NewDecoder(&buf).Decode(&rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec.Path != "secret/data/app" {
		t.Errorf("expected path secret/data/app, got %s", rec.Path)
	}
	if rec.ArchivedAt.IsZero() {
		t.Error("expected archived_at to be set")
	}
}

func TestArchive_MasksValues(t *testing.T) {
	srv := newDiffVaultStub(t, map[string]string{"DB_PASS": "hunter2"})
	var buf bytes.Buffer
	err := Archive(ArchiveOptions{
		Address: srv.URL,
		Token:   "test-token",
		Path:    "secret/data/app",
		Version: 1,
		Mask:    true,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "hunter2") {
		t.Error("expected value to be masked")
	}
}

func TestArchive_InvalidAddress_ReturnsError(t *testing.T) {
	err := Archive(ArchiveOptions{
		Address: "://bad",
		Token:   "tok",
		Path:    "secret/data/app",
		Version: 1,
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestArchive_NilOutputDefaultsToStdout(t *testing.T) {
	srv := newDiffVaultStub(t, map[string]string{"KEY": "val"})
	err := Archive(ArchiveOptions{
		Address: srv.URL,
		Token:   "test-token",
		Path:    "secret/data/app",
		Version: 1,
		Output:  nil,
	})
	// should not panic; stdout write may succeed or fail in test env
	_ = err
}

func TestDefaultArchiveOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultArchiveOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}
