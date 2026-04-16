package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

// envVaultStub satisfies vault.SecretGetter for env tests.
type envVaultStub struct {
	data map[string]interface{}
	err  error
}

func (s *envVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &vault.SecretVersion{Path: path, Version: version, Data: s.data}, nil
}

func TestPrintEnv_KeyValueFormat(t *testing.T) {
	stub := &envVaultStub{data: map[string]interface{}{"DB_HOST": "localhost", "PORT": "5432"}}
	var buf bytes.Buffer
	opts := DefaultEnvOptions()
	opts.Output = &buf
	opts.Mask = false

	if err := PrintEnv(stub, "secret/app", 1, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=5432") {
		t.Errorf("expected PORT=5432 in output, got:\n%s", out)
	}
}

func TestPrintEnv_MasksValues(t *testing.T) {
	stub := &envVaultStub{data: map[string]interface{}{"SECRET_KEY": "super-secret"}}
	var buf bytes.Buffer
	opts := DefaultEnvOptions()
	opts.Output = &buf
	opts.Mask = true

	if err := PrintEnv(stub, "secret/app", 1, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "super-secret") {
		t.Errorf("expected value to be masked, got:\n%s", out)
	}
	if !strings.Contains(out, "SECRET_KEY=***") {
		t.Errorf("expected SECRET_KEY=*** in output, got:\n%s", out)
	}
}

func TestPrintEnv_SortedKeys(t *testing.T) {
	stub := &envVaultStub{data: map[string]interface{}{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}}
	var buf bytes.Buffer
	opts := DefaultEnvOptions()
	opts.Output = &buf
	opts.Mask = false

	if err := PrintEnv(stub, "secret/app", 2, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to be A_KEY, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to be Z_KEY, got %s", lines[2])
	}
}

func TestPrintEnv_FetchError(t *testing.T) {
	stub := &envVaultStub{err: errors.New("vault unavailable")}
	var buf bytes.Buffer
	opts := DefaultEnvOptions()
	opts.Output = &buf

	err := PrintEnv(stub, "secret/app", 1, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPrintEnv_NilOutputDefaultsToStdout(t *testing.T) {
	stub := &envVaultStub{data: map[string]interface{}{"K": "v"}}
	opts := DefaultEnvOptions()
	opts.Output = nil // should not panic

	if err := PrintEnv(stub, "secret/app", 1, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
