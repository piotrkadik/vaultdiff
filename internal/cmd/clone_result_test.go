package cmd_test

import (
	"encoding/json"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
)

func TestCloneResult_DryRunTrue_DataOmitted(t *testing.T) {
	r := cmd.CloneResult{
		SrcPath: "app/prod",
		DstPath: "app/staging",
		Keys:    3,
		DryRun:  true,
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := m["data"]; ok {
		t.Error("data field should be omitted when nil")
	}
}

func TestCloneResult_DataIncludedWhenSet(t *testing.T) {
	r := cmd.CloneResult{
		SrcPath: "a",
		DstPath: "b",
		Keys:    1,
		DryRun:  false,
		Data:    map[string]any{"KEY": "val"},
	}
	b, _ := json.Marshal(r)
	var m map[string]any
	json.Unmarshal(b, &m)
	if _, ok := m["data"]; !ok {
		t.Error("expected data field in JSON output")
	}
}

func TestCloneResult_KeysCount(t *testing.T) {
	r := cmd.CloneResult{Keys: 7}
	if r.Keys != 7 {
		t.Errorf("expected 7, got %d", r.Keys)
	}
}
