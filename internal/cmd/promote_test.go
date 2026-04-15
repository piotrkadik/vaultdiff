package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

// promoteStub satisfies both vault.SecretGetter and vault.SecretWriter.
type promoteStub struct {
	data    map[string]any
	getErr  error
	writeErr error
	written map[string]any
}

func (s *promoteStub) GetSecretVersion(_ context.Context, _ string, _ int) (*vault.SecretVersion, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return &vault.SecretVersion{Data: s.data, Version: 2}, nil
}

func (s *promoteStub) WriteSecret(_ context.Context, _ string, data map[string]any) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	s.written = data
	return nil
}

func TestPromote_WritesJSON(t *testing.T) {
	stub := &promoteStub{data: map[string]any{"key": "val"}}
	var buf bytes.Buffer
	opts := DefaultPromoteOptions()
	opts.Output = &buf

	if err := Promote(context.Background(), stub, stub, "app/cfg", 2, "staging", "prod", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var rec PromoteRecord
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec.FromEnv != "staging" || rec.ToEnv != "prod" {
		t.Errorf("unexpected envs: %s -> %s", rec.FromEnv, rec.ToEnv)
	}
	if rec.DryRun {
		t.Error("expected DryRun=false")
	}
}

func TestPromote_DryRun_SkipsWrite(t *testing.T) {
	stub := &promoteStub{data: map[string]any{"x": "y"}}
	var buf bytes.Buffer
	opts := DefaultPromoteOptions()
	opts.DryRun = true
	opts.Output = &buf

	if err := Promote(context.Background(), stub, stub, "app/cfg", 1, "dev", "staging", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.written != nil {
		t.Error("expected no write on dry run")
	}
	var rec PromoteRecord
	_ = json.Unmarshal(buf.Bytes(), &rec)
	if !rec.DryRun {
		t.Error("expected DryRun=true in record")
	}
}

func TestPromote_FetchError_ReturnsError(t *testing.T) {
	stub := &promoteStub{getErr: errors.New("not found")}
	opts := DefaultPromoteOptions()
	opts.Output = &bytes.Buffer{}

	err := Promote(context.Background(), stub, stub, "app/cfg", 99, "a", "b", opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPromote_NilOutputDefaultsToStdout(t *testing.T) {
	stub := &promoteStub{data: map[string]any{}}
	opts := DefaultPromoteOptions()
	opts.Output = nil // should fall back to os.Stdout without panic

	// We can't easily capture stdout here; just ensure no panic.
	_ = Promote(context.Background(), stub, stub, "app/cfg", 1, "a", "b", opts)
}
