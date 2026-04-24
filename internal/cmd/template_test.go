package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

type templateVaultStub struct {
	data    map[string]string
	version int
	err     error
}

func (s *templateVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &vault.SecretVersion{Data: s.data, Version: s.version}, nil
}

func (s *templateVaultStub) ListVersions(path string) ([]vault.VersionMeta, error) { return nil, nil }
func (s *templateVaultStub) WriteSecret(path string, data map[string]string) error  { return nil }

func TestTemplate_RendersValues(t *testing.T) {
	stub := &templateVaultStub{
		data:    map[string]string{"HOST": "db.internal", "PORT": "5432"},
		version: 3,
	}
	var buf bytes.Buffer
	opts := DefaultTemplateOptions()
	opts.Path = "infra/db"
	opts.Version = 3
	opts.Template = "postgres://{{.HOST}}:{{.PORT}}/mydb"
	opts.Mask = false
	opts.Output = &buf

	res, err := Template(stub, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rendered != "postgres://db.internal:5432/mydb" {
		t.Errorf("rendered = %q, want postgres://db.internal:5432/mydb", res.Rendered)
	}
}

func TestTemplate_MaskHidesData(t *testing.T) {
	stub := &templateVaultStub{
		data:    map[string]string{"SECRET": "topsecret"},
		version: 1,
	}
	var buf bytes.Buffer
	opts := DefaultTemplateOptions()
	opts.Path = "app/creds"
	opts.Template = "key={{.SECRET}}"
	opts.Mask = true
	opts.Output = &buf

	res, err := Template(stub, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Data != nil {
		t.Errorf("expected Data to be nil when masked, got %v", res.Data)
	}
	// Rendered still contains the real value (template expansion is pre-mask)
	if res.Rendered != "key=topsecret" {
		t.Errorf("rendered = %q", res.Rendered)
	}
}

func TestTemplate_JSONFormat_ValidOutput(t *testing.T) {
	stub := &templateVaultStub{
		data:    map[string]string{"A": "1"},
		version: 2,
	}
	var buf bytes.Buffer
	opts := DefaultTemplateOptions()
	opts.Path = "svc/config"
	opts.Template = "A={{.A}}"
	opts.Mask = false
	opts.Output = &buf

	_, err := Template(stub, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out TemplateResult
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw: %s", err, buf.String())
	}
	if out.Path != "svc/config" {
		t.Errorf("path = %q", out.Path)
	}
}

func TestTemplate_NilOutputDefaultsToStdout(t *testing.T) {
	stub := &templateVaultStub{data: map[string]string{"X": "y"}, version: 1}
	opts := DefaultTemplateOptions()
	opts.Path = "any/path"
	opts.Template = "{{.X}}"
	opts.Output = nil // should default to os.Stdout without panic

	// Redirect stdout to discard for test cleanliness
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	_, err := Template(stub, opts)
	w.Close()
	os.Stdout = old
	r.Close()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultTemplateOptions_MaskDefaultsTrue(t *testing.T) {
	os.Unsetenv("VAULTDIFF_MASK")
	opts := DefaultTemplateOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}
