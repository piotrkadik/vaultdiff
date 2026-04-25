package cmd_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
	"github.com/your-org/vaultdiff/internal/vault"
)

type refreshStub struct {
	data    map[string]string
	version int
	getErr  error
	writeErr error
	written map[string]string
}

func (s *refreshStub) GetSecretVersion(path string, version int) (vault.SecretVersion, error) {
	if s.getErr != nil {
		return vault.SecretVersion{}, s.getErr
	}
	return vault.SecretVersion{Version: s.version, Data: s.data}, nil
}

func (s *refreshStub) WriteSecret(path string, data map[string]string) error {
	if s.writeErr != nil {
		return s.writeErr
	}
	s.written = data
	return nil
}

func TestRefresh_DryRun_SkipsWrite(t *testing.T) {
	stub := &refreshStub{data: map[string]string{"k": "v"}, version: 3}
	var buf bytes.Buffer
	opts := cmd.DefaultRefreshOptions()
	opts.Path = "svc/api"
	opts.DryRun = true
	opts.Output = &buf

	_, err := cmd.Refresh(stub, stub, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.written != nil {
		t.Error("expected no write in dry-run mode")
	}
	if !strings.Contains(buf.String(), "dry-run") {
		t.Errorf("output missing dry-run notice: %q", buf.String())
	}
}

func TestRefresh_WritesSecret(t *testing.T) {
	stub := &refreshStub{data: map[string]string{"key": "secret"}, version: 2}
	var buf bytes.Buffer
	opts := cmd.DefaultRefreshOptions()
	opts.Path = "svc/db"
	opts.DryRun = false
	opts.Mask = false
	opts.Output = &buf

	result, err := cmd.Refresh(stub, stub, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Refreshed {
		t.Error("expected Refreshed=true")
	}
	if stub.written["key"] != "secret" {
		t.Errorf("unexpected written data: %v", stub.written)
	}
}

func TestRefreshResult_JSONShape(t *testing.T) {
	stub := &refreshStub{data: map[string]string{"a": "b"}, version: 1}
	var buf bytes.Buffer
	opts := cmd.DefaultRefreshOptions()
	opts.Path = "svc/x"
	opts.Format = "json"
	opts.Output = &buf

	_, err := cmd.Refresh(stub, stub, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, field := range []string{"path", "version", "refreshed", "timestamp"} {
		if _, ok := out[field]; !ok {
			t.Errorf("missing JSON field %q", field)
		}
	}
}

func TestRefresh_FetchError_ReturnsError(t *testing.T) {
	stub := &refreshStub{getErr: errors.New("vault unreachable")}
	opts := cmd.DefaultRefreshOptions()
	opts.Path = "svc/y"
	opts.Output = &bytes.Buffer{}

	_, err := cmd.Refresh(stub, stub, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "refresh") {
		t.Errorf("error missing context: %v", err)
	}
}

func TestDefaultRefreshOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultRefreshOptions()
	if !opts.Mask {
		t.Error("expected Mask=true by default")
	}
}
