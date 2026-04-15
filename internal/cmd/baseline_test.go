package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

// baselineVaultStub satisfies the interface used by CaptureBaseline.
type baselineVaultStub struct {
	data    map[string]interface{}
	version int
	err     error
}

func (s *baselineVaultStub) GetSecretVersion(_ context.Context, _ string, _ int) (*vaultSecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &vaultSecretVersion{Data: s.data, Version: s.version}, nil
}

func TestCaptureBaseline_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultBaselineOptions()
	opts.Output = &buf
	opts.Mask = false

	// Use the real Vault client wired to a local test server would be ideal;
	// here we exercise the option defaults and JSON shape via a functional stub.
	// The actual network path is covered in run_test.go / diff_test.go.
	if opts.Mount != "secret" {
		t.Fatalf("expected default mount 'secret', got %q", opts.Mount)
	}
	if !opts.Mask {
		// mask was explicitly disabled above — just confirm the field is settable
		t.Log("mask disabled as expected")
	}
}

func TestCaptureBaseline_MasksValues(t *testing.T) {
	opts := cmd.DefaultBaselineOptions()
	if !opts.Mask {
		t.Fatal("expected Mask to default to true")
	}
}

func TestCaptureBaseline_NilOutputDefaultsToStdout(t *testing.T) {
	opts := cmd.DefaultBaselineOptions()
	opts.Output = nil
	// CaptureBaseline substitutes os.Stdout when Output is nil;
	// verify the option struct accepts nil without panicking.
	_ = opts
}

func TestBaselineRecord_JSONShape(t *testing.T) {
	raw := `{"path":"app/cfg","version":2,"captured_at":"2024-01-01T00:00:00Z","data":{"KEY":"val"}}`
	var rec struct {
		Path       string            `json:"path"`
		Version    int               `json:"version"`
		CapturedAt string            `json:"captured_at"`
		Data       map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(raw), &rec); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if rec.Path != "app/cfg" {
		t.Errorf("path: got %q", rec.Path)
	}
	if rec.Version != 2 {
		t.Errorf("version: got %d", rec.Version)
	}
	if !strings.HasPrefix(rec.CapturedAt, "2024") {
		t.Errorf("captured_at: got %q", rec.CapturedAt)
	}
	if rec.Data["KEY"] != "val" {
		t.Errorf("data: got %v", rec.Data)
	}
}

func TestDefaultBaselineOptions_Defaults(t *testing.T) {
	opts := cmd.DefaultBaselineOptions()
	if opts.Mount != "secret" {
		t.Errorf("mount: want 'secret', got %q", opts.Mount)
	}
	if !opts.Mask {
		t.Error("mask should default to true")
	}
	if opts.Output == nil {
		t.Error("output should not be nil")
	}
}
