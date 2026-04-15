package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/vaultdiff/internal/diff"
)

func makeAuditChanges() []diff.Change {
	return []diff.Change{
		{Path: "secret/app", Key: "DB_PASS", Type: diff.Modified, OldValue: "old", NewValue: "new"},
		{Path: "secret/app", Key: "API_KEY", Type: diff.Added, NewValue: "key"},
		{Path: "secret/app", Key: "LEGACY", Type: diff.Removed, OldValue: "val"},
		{Path: "secret/app", Key: "HOST", Type: diff.Unchanged, OldValue: "h", NewValue: "h"},
	}
}

func TestWriteAuditEntry_DisabledIsNoop(t *testing.T) {
	var buf bytes.Buffer
	err := WriteAuditEntry("secret/app", 1, 2, makeAuditChanges(), time.Now(), AuditOptions{
		Writer:  &buf,
		Enabled: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output when disabled, got %q", buf.String())
	}
}

func TestWriteAuditEntry_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	err := WriteAuditEntry("secret/app", 1, 2, makeAuditChanges(), ts, AuditOptions{
		Writer:  &buf,
		Actor:   "ci-bot",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &entry); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if entry["path"] != "secret/app" {
		t.Errorf("expected path=secret/app, got %v", entry["path"])
	}
	if entry["actor"] != "ci-bot" {
		t.Errorf("expected actor=ci-bot, got %v", entry["actor"])
	}
	if entry["has_drift"] != true {
		t.Errorf("expected has_drift=true")
	}
	if entry["added"].(float64) != 1 {
		t.Errorf("expected added=1, got %v", entry["added"])
	}
	if entry["removed"].(float64) != 1 {
		t.Errorf("expected removed=1, got %v", entry["removed"])
	}
	if entry["modified"].(float64) != 1 {
		t.Errorf("expected modified=1, got %v", entry["modified"])
	}
}

func TestDefaultAuditOptions_EnabledFromEnv(t *testing.T) {
	t.Setenv("VAULTDIFF_AUDIT", "true")
	t.Setenv("VAULTDIFF_ACTOR", "tester")

	opts := DefaultAuditOptions()
	if !opts.Enabled {
		t.Error("expected Enabled=true from env")
	}
	if opts.Actor != "tester" {
		t.Errorf("expected Actor=tester, got %q", opts.Actor)
	}
}
