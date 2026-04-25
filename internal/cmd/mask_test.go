package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

// maskVaultStub satisfies the vault.Client interface used by Mask.
type maskVaultStub struct {
	data    map[string]string
	version int
	written map[string]string
	wantErr error
}

func (s *maskVaultStub) GetSecretVersion(path string, v int) (*vault.SecretVersion, error) {
	if s.wantErr != nil {
		return nil, s.wantErr
	}
	return &vault.SecretVersion{Data: s.data, Version: s.version}, nil
}

func (s *maskVaultStub) WriteSecret(path string, data map[string]string) error {
	s.written = data
	return nil
}

func TestMask_DryRun_SkipsWrite(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultMaskOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.DryRun = true
	opts.Output = &buf

	err := Mask("secret/app", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result MaskResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !result.DryRun {
		t.Error("expected dry_run=true")
	}
}

func TestMaskResult_JSONShape(t *testing.T) {
	r := MaskResult{
		Path:    "secret/app",
		Version: 2,
		DryRun:  false,
		Masked:  []string{"password"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{`"path"`, `"version"`, `"masked_keys"`} {
		if !strings.Contains(s, want) {
			t.Errorf("missing field %s in JSON", want)
		}
	}
}

func TestDefaultMaskOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultMaskOptions()
	if !opts.Mask {
		t.Error("expected Mask=true by default")
	}
}

func TestDefaultMaskOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultMaskOptions()
	if opts.Mask {
		t.Error("expected Mask=false when VAULTDIFF_MASK=false")
	}
}

func TestMask_InvalidAddress_ReturnsError(t *testing.T) {
	opts := DefaultMaskOptions()
	opts.Address = "://bad"
	opts.Output = &bytes.Buffer{}
	err := Mask("secret/app", opts)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}
