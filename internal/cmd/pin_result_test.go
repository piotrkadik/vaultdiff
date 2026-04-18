package cmd

import (
	"encoding/json"
	"testing"
)

func TestPinResult_DryRunTrue_DataOmitted(t *testing.T) {
	r := PinResult{Path: "x/y", Version: 2, DryRun: true}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if _, ok := out["data"]; ok {
		t.Error("data should be omitted when nil")
	}
}

func TestPinResult_DataIncludedWhenSet(t *testing.T) {
	r := PinResult{
		Path:    "x/y",
		Version: 1,
		Data:    map[string]any{"foo": "bar"},
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if _, ok := out["data"]; !ok {
		t.Error("expected data field in JSON")
	}
}

func TestPinResult_PathAndVersion(t *testing.T) {
	r := PinResult{Path: "prod/db", Version: 5}
	if r.Path != "prod/db" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Version != 5 {
		t.Errorf("unexpected version: %d", r.Version)
	}
}
