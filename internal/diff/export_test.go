package diff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeExportChanges() []Change {
	return []Change{
		{Path: "secret/db/password", Type: ChangeAdded, OldValue: "", NewValue: "newpass"},
		{Path: "secret/db/user", Type: ChangeModified, OldValue: "root", NewValue: "admin"},
		{Path: "secret/api/key", Type: ChangeRemoved, OldValue: "oldkey", NewValue: ""},
	}
}

func TestExport_JSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, makeExportChanges(), ExportJSON, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
}

func TestExport_JSON_ContainsPath(t *testing.T) {
	var buf bytes.Buffer
	_ = Export(&buf, makeExportChanges(), ExportJSON, false)
	if !strings.Contains(buf.String(), "secret/db/password") {
		t.Error("expected path in JSON output")
	}
}

func TestExport_JSON_MasksValues(t *testing.T) {
	var buf bytes.Buffer
	_ = Export(&buf, makeExportChanges(), ExportJSON, true)
	if strings.Contains(buf.String(), "newpass") || strings.Contains(buf.String(), "oldkey") {
		t.Error("expected values to be masked")
	}
}

func TestExport_CSV_HasHeader(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, makeExportChanges(), ExportCSV, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "timestamp,path") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExport_CSV_RowCount(t *testing.T) {
	var buf bytes.Buffer
	_ = Export(&buf, makeExportChanges(), ExportCSV, false)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestExport_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, makeExportChanges(), ExportFormat("xml"), false)
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestExport_EmptyChanges_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Export(&buf, []Change{}, ExportJSON, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []ExportRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}
