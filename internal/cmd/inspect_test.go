package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

type inspectVaultStub struct{}

func newInspectStub(t *testing.T) string {
	t.Helper()
	srv := newVaultTestServer(t, map[string]map[string]string{
		"secret/data/myapp": {"API_KEY": "abc123", "DB_PASS": "hunter2"},
	})
	return srv
}

func TestInspect_TextFormat_PrintsKeyValues(t *testing.T) {
	addr := newInspectStub(t)
	var buf bytes.Buffer
	err := cmd.Inspect(cmd.InspectOptions{
		Address: addr,
		Mount:   "secret",
		Path:    "myapp",
		Version: 1,
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got: %s", out)
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestInspect_MasksValues(t *testing.T) {
	addr := newInspectStub(t)
	var buf bytes.Buffer
	_ = cmd.Inspect(cmd.InspectOptions{
		Address: addr,
		Mount:   "secret",
		Path:    "myapp",
		Version: 1,
		Mask:    true,
		Output:  &buf,
	})
	out := buf.String()
	if strings.Contains(out, "abc123") {
		t.Errorf("expected value to be masked, got: %s", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected *** mask in output, got: %s", out)
	}
}

func TestInspect_JSONFormat_ValidOutput(t *testing.T) {
	addr := newInspectStub(t)
	var buf bytes.Buffer
	err := cmd.Inspect(cmd.InspectOptions{
		Address: addr,
		Mount:   "secret",
		Path:    "myapp",
		Version: 1,
		Format:  "json",
		Mask:    false,
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result cmd.InspectResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Path != "myapp" {
		t.Errorf("expected path myapp, got %s", result.Path)
	}
}

func TestInspect_InvalidAddress_ReturnsError(t *testing.T) {
	err := cmd.Inspect(cmd.InspectOptions{
		Address: "http://127.0.0.1:0",
		Mount:   "secret",
		Path:    "myapp",
		Version: 1,
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultInspectOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultInspectOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}
