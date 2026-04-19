package cmd

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEncryptResult_DryRunTrue_DataOmitted(t *testing.T) {
	res := EncryptResult{
		Path:   "app/cfg",
		DryRun: true,
		Keys:   []string{"password"},
	}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(b), `"data"`) {
		t.Error("data field should be omitted when nil")
	}
}

func TestEncryptResult_DataIncludedWhenSet(t *testing.T) {
	res := EncryptResult{
		Path:   "app/cfg",
		DryRun: false,
		Keys:   []string{"password"},
		Data:   map[string]string{"password": "enc:736563726574"},
	}
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(b), `"data"`) {
		t.Error("expected 'data' field in JSON when set")
	}
}

func TestEncryptResult_KeysCount(t *testing.T) {
	res := EncryptResult{
		Path: "app/cfg",
		Keys: []string{"a", "b", "c"},
	}
	if len(res.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Keys))
	}
}
