package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type decryptVaultStub struct {
	data map[string]string
	err  error
}

func (s *decryptVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &vault.SecretVersion{Path: path, Version: version, Data: s.data}, nil
}

func TestDecrypt_TextFormat_PrintsKeyValues(t *testing.T) {
	stub := &decryptVaultStub{data: map[string]string{"KEY": "val"}}
	var buf bytes.Buffer
	err := Decrypt(stub, DecryptOptions{Path: "secret/app", Version: 1, Mask: false, Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY=val") {
		t.Errorf("expected KEY=val in output, got: %s", buf.String())
	}
}

func TestDecrypt_MasksValues(t *testing.T) {
	stub := &decryptVaultStub{data: map[string]string{"SECRET": "plaintext"}}
	var buf bytes.Buffer
	err := Decrypt(stub, DecryptOptions{Path: "secret/app", Version: 2, Mask: true, Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(buf.String(), "plaintext") {
		t.Errorf("expected value to be masked")
	}
}

func TestDecrypt_JSONFormat_ValidOutput(t *testing.T) {
	stub := &decryptVaultStub{data: map[string]string{"K": "v"}}
	var buf bytes.Buffer
	err := Decrypt(stub, DecryptOptions{Path: "secret/app", Version: 1, Mask: false, Format: "json", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var res DecryptResult
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if res.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", res.Path)
	}
}

func TestDecrypt_InvalidAddress_ReturnsError(t *testing.T) {
	stub := &decryptVaultStub{err: errors.New("connection refused")}
	var buf bytes.Buffer
	err := Decrypt(stub, DecryptOptions{Path: "secret/app", Version: 1, Output: &buf})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDecrypt_NilOutput_DefaultsToStdout(t *testing.T) {
	stub := &decryptVaultStub{data: map[string]string{"X": "y"}}
	err := Decrypt(stub, DecryptOptions{Path: "secret/app", Version: 1, Mask: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
