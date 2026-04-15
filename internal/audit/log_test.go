package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultdiff/internal/audit"
)

func TestWrite_ProducesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	e := audit.Entry{
		Environment: "production",
		Path:        "secret/data/app",
		VersionA:    1,
		VersionB:    2,
		Added:       1,
		Removed:     0,
		Modified:    2,
		HasDrift:    true,
	}

	if err := logger.Write(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if decoded.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", decoded.Environment)
	}
	if !decoded.HasDrift {
		t.Error("expected has_drift to be true")
	}
}

func TestWrite_StampsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	before := time.Now().UTC()
	_ = logger.Write(audit.Entry{Path: "secret/data/test"})
	after := time.Now().UTC()

	var decoded audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &decoded)

	if decoded.Timestamp.Before(before) || decoded.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", decoded.Timestamp, before, after)
	}
}

func TestWrite_PreservesExplicitTimestamp(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	_ = logger.Write(audit.Entry{Timestamp: fixed, Path: "x"})

	var decoded audit.Entry
	_ = json.Unmarshal(buf.Bytes(), &decoded)

	if !decoded.Timestamp.Equal(fixed) {
		t.Errorf("expected timestamp %v, got %v", fixed, decoded.Timestamp)
	}
}

func TestWrite_NewlineDelimited(t *testing.T) {
	var buf bytes.Buffer
	logger := audit.NewLogger(&buf)

	_ = logger.Write(audit.Entry{Path: "a"})
	_ = logger.Write(audit.Entry{Path: "b"})

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestNewLogger_NilWriterDefaultsToStderr(t *testing.T) {
	// Should not panic when nil is passed.
	logger := audit.NewLogger(nil)
	if logger == nil {
		t.Error("expected non-nil logger")
	}
}
