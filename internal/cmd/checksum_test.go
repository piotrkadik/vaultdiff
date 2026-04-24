package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
	"github.com/your-org/vaultdiff/internal/vault"
)

type checksumVaultStub struct {
	data    map[string]string
	version int
}

func (s *checksumVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	return &vault.SecretVersion{Version: s.version, Data: s.data}, nil
}

func newChecksumStub() *checksumVaultStub {
	return &checksumVaultStub{
		version: 3,
		data:    map[string]string{"alpha": "a", "beta": "b"},
	}
}

func TestChecksum_TextFormat_PrintsChecksum(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultChecksumOptions()
	opts.Output = &buf
	opts.Path = "secret/app"
	opts.Version = 3

	if err := cmd.Checksum(newChecksumStub(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "checksum:") {
		t.Errorf("expected 'checksum:' label, got: %s", got)
	}
	if !strings.Contains(got, "secret/app") {
		t.Errorf("expected path in output, got: %s", got)
	}
}

func TestChecksum_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultChecksumOptions()
	opts.Output = &buf
	opts.Path = "secret/app"
	opts.Version = 3
	opts.Format = "json"

	if err := cmd.Checksum(newChecksumStub(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var res cmd.ChecksumResult
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if res.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if res.Keys != 2 {
		t.Errorf("expected 2 keys, got %d", res.Keys)
	}
}

func TestChecksum_IsDeterministic(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	stub := newChecksumStub()

	opts := cmd.DefaultChecksumOptions()
	opts.Format = "json"

	opts.Output = &buf1
	_ = cmd.Checksum(stub, opts)

	opts.Output = &buf2
	_ = cmd.Checksum(stub, opts)

	var r1, r2 cmd.ChecksumResult
	_ = json.Unmarshal(buf1.Bytes(), &r1)
	_ = json.Unmarshal(buf2.Bytes(), &r2)

	if r1.Checksum != r2.Checksum {
		t.Errorf("checksum not deterministic: %s != %s", r1.Checksum, r2.Checksum)
	}
}

func TestChecksum_NilOutput_DefaultsToStdout(t *testing.T) {
	opts := cmd.DefaultChecksumOptions()
	opts.Output = nil
	opts.Path = "secret/app"

	// should not panic
	_ = cmd.Checksum(newChecksumStub(), opts)
}

func TestChecksumResult_JSONShape(t *testing.T) {
	res := cmd.ChecksumResult{
		Path:     "secret/app",
		Version:  2,
		Checksum: "abc123",
		Keys:     4,
	}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	got := string(b)
	for _, want := range []string{"path", "version", "checksum", "keys"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected field %q in JSON, got: %s", want, got)
		}
	}
}
