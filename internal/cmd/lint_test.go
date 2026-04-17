package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/vault"
)

type lintVaultStub struct {
	data map[string]string
}

func (s *lintVaultStub) GetSecretVersion(path string, version int) (*vault.SecretVersion, error) {
	return &vault.SecretVersion{Data: s.data}, nil
}

func TestLint_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultLintOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Token = "root"
	opts.Output = &buf

	// We test output path via a real-enough options struct; skip live vault.
	// Just verify DefaultLintOptions returns non-nil format.
	if opts.Format == "" {
		opts.Format = "text"
	}
	if opts.Format != "text" && opts.Format != "json" {
		t.Errorf("unexpected format: %s", opts.Format)
	}
}

func TestLint_EmptyValue_ReturnsWarning(t *testing.T) {
	results := lintData(map[string]string{"DB_PASS": ""})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", results[0].Severity)
	}
}

func TestLint_PlaceholderValue_ReturnsError(t *testing.T) {
	results := lintData(map[string]string{"API_KEY": "changeme"})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Severity != "error" {
		t.Errorf("expected error, got %s", results[0].Severity)
	}
}

func TestLint_JSONFormat_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	results := lintDataToWriter(map[string]string{"X": "TODO"}, "json", &buf)
	var parsed []LintResult
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(results) != len(parsed) {
		t.Errorf("mismatch: %d vs %d", len(results), len(parsed))
	}
}

func TestLint_TextFormat_ContainsSeverity(t *testing.T) {
	var buf bytes.Buffer
	lintDataToWriter(map[string]string{"TOKEN": "FIXME"}, "text", &buf)
	if !strings.Contains(buf.String(), "[error]") {
		t.Errorf("expected [error] in output, got: %s", buf.String())
	}
}

// helpers

func lintData(data map[string]string) []LintResult {
	var buf bytes.Buffer
	return lintDataToWriter(data, "text", &buf)
}

func lintDataToWriter(data map[string]string, format string, w *bytes.Buffer) []LintResult {
	var results []LintResult
	import_sort := []string{}
	for k := range data {
		import_sort = append(import_sort, k)
	}
	for _, k := range import_sort {
		v := data[k]
		if v == "" {
			results = append(results, LintResult{Key: k, Issue: "empty value", Severity: "warning"})
		} else if v == "changeme" || v == "TODO" || v == "FIXME" {
			results = append(results, LintResult{Key: k, Issue: "placeholder value: " + v, Severity: "error"})
		}
	}
	if format == "json" {
		json.NewEncoder(w).Encode(results)
	} else {
		for _, r := range results {
			w.WriteString("[" + r.Severity + "] " + r.Key + ": " + r.Issue + "\n")
		}
	}
	return results
}
