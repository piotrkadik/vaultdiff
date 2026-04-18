package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/vault"
)

type statVaultStub struct{}

func newStatStub(t *testing.T) string {
	t.Helper()
	srv := newVaultTestServer(t, map[string]vault.SecretVersion{
		"secret/data/myapp": {Data: map[string]string{"key1": "val1", "key2": "val2"}},
	})
	return srv
}

func TestStat_TextFormat_PrintsMetadata(t *testing.T) {
	addr := newStatStub(t)
	var buf bytes.Buffer
	err := Stat("myapp", StatOptions{
		Address: addr,
		Mount:   "secret",
		Token:   "root",
		Format:  "text",
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "Keys:") {
		t.Errorf("expected Keys line in output, got: %s", out)
	}
}

func TestStat_JSONFormat_ValidOutput(t *testing.T) {
	addr := newStatStub(t)
	var buf bytes.Buffer
	err := Stat("myapp", StatOptions{
		Address: addr,
		Mount:   "secret",
		Token:   "root",
		Format:  "json",
		Output:  &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var res StatResult
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if res.Path != "myapp" {
		t.Errorf("expected path myapp, got %s", res.Path)
	}
}

func TestStat_InvalidAddress_ReturnsError(t *testing.T) {
	err := Stat("myapp", StatOptions{
		Address: "http://127.0.0.1:0",
		Mount:   "secret",
		Token:   "root",
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultStatOptions_MaskDefaultsTrue(t *testing.T) {
	t.Setenv("VAULTDIFF_MASK", "")
	opts := DefaultStatOptions()
	if !opts.Mask {
		t.Error("expected Mask to default to true")
	}
}

func TestStatResult_JSONShape(t *testing.T) {
	res := StatResult{Path: "p", LatestVersion: 3, TotalVersions: 5, DestroyedCount: 1, KeyCount: 4}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "latest_version") {
		t.Errorf("expected latest_version key in JSON")
	}
}
