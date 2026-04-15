package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
	"github.com/your-org/vaultdiff/internal/vault"
)

func TestTakeSnapshot_WritesJSON(t *testing.T) {
	c := vaultStub(map[string]map[string]interface{}{
		"secret/data/app": {"key": "value"},
	})

	var buf bytes.Buffer
	opts := cmd.SnapshotOptions{Mask: false, Output: &buf}

	snap, err := cmd.TakeSnapshot(context.Background(), c, "secret/data/app", 1, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if snap.Path != "secret/data/app" {
		t.Errorf("path = %q, want %q", snap.Path, "secret/data/app")
	}
	if snap.Version != 1 {
		t.Errorf("version = %d, want 1", snap.Version)
	}
	if snap.CapturedAt.IsZero() {
		t.Error("captured_at should not be zero")
	}

	var out cmd.Snapshot
	if err := json.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

func TestTakeSnapshot_MasksValues(t *testing.T) {
	c := vaultStub(map[string]map[string]interface{}{
		"secret/data/app": {"db_pass": "s3cr3t"},
	})

	var buf bytes.Buffer
	opts := cmd.SnapshotOptions{Mask: true, Output: &buf}

	snap, err := cmd.TakeSnapshot(context.Background(), c, "secret/data/app", 1, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if snap.Data["db_pass"] != "***" {
		t.Errorf("expected masked value, got %q", snap.Data["db_pass"])
	}
}

func TestTakeSnapshot_ErrorPropagated(t *testing.T) {
	c := vaultStub(nil) // empty stub — path not found

	var buf bytes.Buffer
	opts := cmd.DefaultSnapshotOptions()
	opts.Output = &buf

	_, err := cmd.TakeSnapshot(context.Background(), c, "secret/data/missing", 1, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTakeSnapshot_NilOutputDefaultsToStdout(t *testing.T) {
	_ = errors.New // import guard
	_ = vault.Client{}

	opts := cmd.DefaultSnapshotOptions()
	if opts.Output == nil {
		t.Error("DefaultSnapshotOptions should set a non-nil Output")
	}
	if !opts.Mask {
		t.Error("DefaultSnapshotOptions should default Mask to true")
	}
}
