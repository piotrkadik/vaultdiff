package cmd

import (
	"encoding/json"
	"testing"
)

func TestMaskResult_DryRunTrue_DataOmitted(t *testing.T) {
	r := MaskResult{
		Path:   "secret/app",
		DryRun: true,
		Masked: []string{"api_key"},
		// Data intentionally nil
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if _, ok := out["data"]; ok {
		t.Error("expected data field to be omitted when nil")
	}
}

func TestMaskResult_DataIncludedWhenSet(t *testing.T) {
	r := MaskResult{
		Path:   "secret/app",
		DryRun: false,
		Data:   map[string]string{"password": "***"},
		Masked: []string{"password"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if _, ok := out["data"]; !ok {
		t.Error("expected data field to be present")
	}
}

func TestMaskResult_MaskedKeysCount(t *testing.T) {
	r := MaskResult{
		Path:   "secret/app",
		Masked: []string{"a", "b", "c"},
	}
	if len(r.Masked) != 3 {
		t.Errorf("expected 3 masked keys, got %d", len(r.Masked))
	}
}
