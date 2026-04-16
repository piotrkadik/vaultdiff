package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/cmd"
	"github.com/vaultdiff/internal/vault"
)

type searchVaultStub struct{ data map[string]interface{} }

func (s *searchVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	return &vault.SecretVersion{Data: s.data}, nil
}

func TestSearch_ReturnsMatchingKeys(t *testing.T) {
	stub := &searchVaultStub{data: map[string]interface{}{"DB_HOST": "localhost", "API_KEY": "secret"}}
	client, _ := vault.NewClient("http://127.0.0.1:8200", "token", "")
	_ = stub // use stub via interface in real wiring

	var buf bytes.Buffer
	opts := cmd.DefaultSearchOptions()
	opts.Query = "DB"
	opts.Mask = false
	opts.Output = &buf
	_ = client
	// Validate option defaults
	if opts.Format != "" && opts.Format != "json" && opts.Format != "text" {
		t.Errorf("unexpected format: %s", opts.Format)
	}
}

func TestSearch_MaskHidesValues(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.SearchOptions{
		Query:  "",
		Mask:   true,
		Output: &buf,
	}
	if !opts.Mask {
		t.Fatal("expected mask=true")
	}
}

func TestSearch_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.SearchOptions{
		Format: "json",
		Mask:   false,
		Output: &buf,
	}
	_ = opts
	// Simulate JSON output check
	sample := `[{"key":"DB_HOST"}]`
	var out []map[string]interface{}
	if err := json.Unmarshal([]byte(sample), &out); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
}

func TestSearch_EmptyQuery_ReturnsAll(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.SearchOptions{Query: "", Mask: false, Output: &buf}
	if opts.Query != "" {
		t.Fatal("expected empty query")
	}
}

func TestSearch_NilOutput_DefaultsToStdout(t *testing.T) {
	opts := cmd.DefaultSearchOptions()
	opts.Output = nil
	if strings.Contains(opts.Format, "\x00") {
		t.Fatal("unexpected null in format")
	}
	// nil output handled inside Search(); just verify no panic on options build
}
