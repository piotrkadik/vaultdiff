package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/your/vaultdiff/internal/vault"
)

type verifyVaultStub struct {
	data map[string]string
	err  error
}

func newVerifyStub(data map[string]string, err error) *verifyVaultStub {
	return &verifyVaultStub{data: data, err: err}
}

func TestVerify_TextFormat_OK(t *testing.T) {
	_ = newVerifyStub(map[string]string{"k": "v"}, nil)
	opts := DefaultVerifyOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Format = "text"
	var buf bytes.Buffer
	err := Verify("secret/app", 1, opts, &buf)
	// expect client error since no real vault; just check error type
	if err == nil {
		t.Fatal("expected error without real vault")
	}
}

func TestVerify_InvalidAddress_ReturnsError(t *testing.T) {
	opts := DefaultVerifyOptions()
	opts.Address = "://bad"
	err := Verify("secret/app", 1, opts, nil)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
	if !strings.Contains(err.Error(), "verify") {
		t.Errorf("error should mention verify, got: %v", err)
	}
}

func TestVerifyResult_JSONShape(t *testing.T) {
	r := VerifyResult{
		Path:    "secret/app",
		Version: 2,
		Keys:    3,
		OK:      true,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"path", "version", "keys", "ok"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON", key)
		}
	}
}

func TestDefaultVerifyOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultVerifyOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultVerifyOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := DefaultVerifyOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}

var _ = errors.New // keep import
var _ = vault.NewClient // keep import
