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

type policyStub struct {
	data map[string]string
	err  error
}

func (s *policyStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &vault.SecretVersion{Data: s.data, Version: 3}, nil
}

func TestPolicy_TextFormat_PrintsKeys(t *testing.T) {
	var buf bytes.Buffer
	client := &policyStub{data: map[string]string{"db_password": "s3cr3t", "host": "localhost"}}
	opts := cmd.DefaultPolicyOptions()
	opts.Path = "app/config"
	opts.Output = &buf

	if err := cmd.Policy(client, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "app/config") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "sensitive") {
		t.Errorf("expected 'sensitive' classification, got: %s", out)
	}
}

func TestPolicy_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	client := &policyStub{data: map[string]string{"api_token": "abc", "region": "us-east-1"}}
	opts := cmd.DefaultPolicyOptions()
	opts.Path = "app/config"
	opts.Format = "json"
	opts.Output = &buf

	if err := cmd.Policy(client, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result cmd.PolicyResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Path != "app/config" {
		t.Errorf("expected path app/config, got %s", result.Path)
	}
	if result.Policies["api_token"] != "sensitive" {
		t.Errorf("expected api_token to be sensitive")
	}
}

func TestPolicy_InvalidAddress_ReturnsError(t *testing.T) {
	client := &policyStub{err: errors.New("connection refused")}
	opts := cmd.DefaultPolicyOptions()
	opts.Path = "app/config"
	opts.Output = &bytes.Buffer{}

	if err := cmd.Policy(client, opts); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultPolicyOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultPolicyOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestPolicyResult_JSONShape(t *testing.T) {
	r := cmd.PolicyResult{
		Path:     "x/y",
		Policies: map[string]string{"k": "general"},
		Version:  1,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if !strings.Contains(string(b), `"path"`) {
		t.Errorf("expected path field in JSON")
	}
}
