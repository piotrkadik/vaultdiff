package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makePatch() Patch {
	changes := []Change{
		{Key: "DB_HOST", Type: ChangeAdded, ToValue: "db.prod"},
		{Key: "DB_PASS", Type: ChangeModified, FromValue: "old", ToValue: "new"},
		{Key: "LEGACY", Type: ChangeRemoved, FromValue: "yes"},
	}
	return NewPatch("secret/myapp", 3, 4, changes)
}

func TestNewPatch_Fields(t *testing.T) {
	p := makePatch()
	if p.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", p.Path)
	}
	if p.From != 3 || p.To != 4 {
		t.Errorf("unexpected versions: from=%d to=%d", p.From, p.To)
	}
	if len(p.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(p.Changes))
	}
}

func TestWritePatch_JSON(t *testing.T) {
	var buf bytes.Buffer
	p := makePatch()
	if err := WritePatch(&buf, p, PatchFormatJSON); err != nil {
		t.Fatalf("WritePatch JSON error: %v", err)
	}
	var decoded Patch
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if decoded.Path != p.Path {
		t.Errorf("decoded path mismatch: %s", decoded.Path)
	}
	if len(decoded.Changes) != len(p.Changes) {
		t.Errorf("decoded changes count mismatch")
	}
}

func TestWritePatch_Text(t *testing.T) {
	var buf bytes.Buffer
	p := makePatch()
	if err := WritePatch(&buf, p, PatchFormatText); err != nil {
		t.Fatalf("WritePatch text error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "secret/myapp") {
		t.Error("expected path in text output")
	}
	if !strings.Contains(out, "v3 -> v4") {
		t.Error("expected version range in text output")
	}
	if !strings.Contains(out, "+ DB_HOST") {
		t.Error("expected added symbol for DB_HOST")
	}
	if !strings.Contains(out, "- LEGACY") {
		t.Error("expected removed symbol for LEGACY")
	}
	if !strings.Contains(out, "~ DB_PASS") {
		t.Error("expected modified symbol for DB_PASS")
	}
}

func TestWritePatch_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := WritePatch(&buf, makePatch(), PatchFormat("xml")); err == nil {
		t.Error("expected error for unknown format")
	}
}
