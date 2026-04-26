package cmd_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

type treeVaultStub struct {
	data map[string]string
	err  error
}

func (s *treeVaultStub) GetSecretVersion(path string, version int) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.data, nil
}

func TestTree_TextFormat_PrintsKeys(t *testing.T) {
	client := &treeVaultStub{
		data: map[string]string{"alpha": "1", "beta": "2", "gamma": "3"},
	}
	var buf bytes.Buffer
	opts := cmd.DefaultTreeOptions()
	opts.Path = "secret/app"
	opts.Output = &buf

	if err := cmd.Tree(client, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected 'alpha' in output, got: %s", out)
	}
	if !strings.Contains(out, "3 key(s)") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}

func TestTree_JSONFormat_ValidOutput(t *testing.T) {
	client := &treeVaultStub{
		data: map[string]string{"x": "1", "y": "2"},
	}
	var buf bytes.Buffer
	opts := cmd.DefaultTreeOptions()
	opts.Path = "secret/svc"
	opts.Format = "json"
	opts.Output = &buf

	if err := cmd.Tree(client, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result cmd.TreeResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected count=2, got %d", result.Count)
	}
	if result.Path != "secret/svc" {
		t.Errorf("expected path=secret/svc, got %s", result.Path)
	}
}

func TestTree_NilOutput_DefaultsToStdout(t *testing.T) {
	client := &treeVaultStub{data: map[string]string{"k": "v"}}
	opts := cmd.DefaultTreeOptions()
	opts.Path = "secret/test"
	opts.Output = nil
	// should not panic
	_ = cmd.Tree(client, opts)
}

func TestTree_FetchError_ReturnsError(t *testing.T) {
	client := &treeVaultStub{err: fmt.Errorf("connection refused")}
	var buf bytes.Buffer
	opts := cmd.DefaultTreeOptions()
	opts.Path = "secret/bad"
	opts.Output = &buf

	err := cmd.Tree(client, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "tree:") {
		t.Errorf("expected 'tree:' prefix in error, got: %v", err)
	}
}

func TestTree_KeysAreSorted(t *testing.T) {
	client := &treeVaultStub{
		data: map[string]string{"zebra": "z", "apple": "a", "mango": "m"},
	}
	var buf bytes.Buffer
	opts := cmd.DefaultTreeOptions()
	opts.Path = "secret/sorted"
	opts.Output = &buf

	if err := cmd.Tree(client, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	appleIdx := strings.Index(out, "apple")
	mangoIdx := strings.Index(out, "mango")
	zebraIdx := strings.Index(out, "zebra")
	if !(appleIdx < mangoIdx && mangoIdx < zebraIdx) {
		t.Errorf("keys not sorted: apple=%d mango=%d zebra=%d", appleIdx, mangoIdx, zebraIdx)
	}
}
