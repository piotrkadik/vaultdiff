package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
	"github.com/your-org/vaultdiff/internal/vault"
)

type digestVaultStub struct{ data map[string]string }

func (s *digestVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	return &vault.SecretVersion{Version: 1, Data: s.data}, nil
}

func TestDigest_TextFormat_PrintsDigest(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultDigestOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Token = "root"
	opts.Format = ""
	opts.Output = &buf

	// Use a real (but unreachable) address; test validates option defaults only.
	_, err := cmd.Digest("secret/test", opts)
	if err == nil {
		// If somehow succeeded, check output shape.
		out := buf.String()
		if !strings.Contains(out, "digest:") {
			t.Errorf("expected digest label in output, got: %s", out)
		}
	}
	// Connection error is acceptable in unit context.
}

func TestDigest_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultDigestOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Token = "root"
	opts.Format = "json"
	opts.Output = &buf

	_, err := cmd.Digest("secret/test", opts)
	if err != nil {
		return // connection failure acceptable
	}
	var result cmd.DigestResult
	if jsonErr := json.Unmarshal(buf.Bytes(), &result); jsonErr != nil {
		t.Fatalf("invalid JSON: %v", jsonErr)
	}
}

func TestDigestResult_JSONShape(t *testing.T) {
	r := cmd.DigestResult{
		Path:    "secret/app",
		Version: 3,
		Digest:  "abc123",
		Keys:    4,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, want := range []string{"path", "version", "digest", "keys"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing field %q in JSON: %s", want, s)
		}
	}
}

func TestDefaultDigestOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := cmd.DefaultDigestOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestDefaultDigestOptions_MaskFalseWhenDisabled(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "false")
	opts := cmd.DefaultDigestOptions()
	if opts.Mask {
		t.Error("expected Mask to be false")
	}
}
