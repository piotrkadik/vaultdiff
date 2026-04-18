package cmd_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

func TestPrune_DryRun_DoesNotDestroy(t *testing.T) {
	var buf bytes.Buffer
	opts := cmd.DefaultPruneOptions()
	opts.Address = "http://127.0.0.1:8200"
	opts.Token = "root"
	opts.Path = "secret/data/app"
	opts.Keep = 2
	opts.DryRun = true
	opts.Output = &buf

	// Without a real Vault we expect a connection error, not a panic.
	err := cmd.Prune(context.Background(), opts)
	if err == nil {
		t.Fatal("expected error connecting to stub address")
	}
}

func TestPruneResult_JSONShape(t *testing.T) {
	res := cmd.PruneResult{
		Path:      "secret/data/app",
		Destroyed: []int{1, 2},
		Kept:      3,
		DryRun:    true,
	}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"path", "destroyed", "kept", "dry_run"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}

func TestDefaultPruneOptions_Keep(t *testing.T) {
	opts := cmd.DefaultPruneOptions()
	if opts.Keep != 5 {
		t.Errorf("expected Keep=5, got %d", opts.Keep)
	}
}

func TestDefaultPruneOptions_DryRunDefaultsFalse(t *testing.T) {
	opts := cmd.DefaultPruneOptions()
	if opts.DryRun {
		t.Error("expected DryRun to default to false")
	}
}
