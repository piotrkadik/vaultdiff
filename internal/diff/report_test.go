package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeReport(changes []Change) *Report {
	return &Report{
		SourcePath:    "secret/prod/app",
		TargetPath:    "secret/staging/app",
		SourceVersion: 3,
		TargetVersion: 2,
		GeneratedAt:   time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Changes:       changes,
	}
}

func TestReport_Summary(t *testing.T) {
	r := makeReport([]Change{
		{Key: "a", Type: Added},
		{Key: "b", Type: Removed},
		{Key: "c", Type: Modified},
		{Key: "d", Type: Unchanged},
		{Key: "e", Type: Added},
	})
	s := r.Summary()
	if s[Added] != 2 {
		t.Errorf("expected 2 added, got %d", s[Added])
	}
	if s[Removed] != 1 {
		t.Errorf("expected 1 removed, got %d", s[Removed])
	}
	if s[Modified] != 1 {
		t.Errorf("expected 1 modified, got %d", s[Modified])
	}
	if s[Unchanged] != 1 {
		t.Errorf("expected 1 unchanged, got %d", s[Unchanged])
	}
}

func TestReport_HasDrift_True(t *testing.T) {
	r := makeReport([]Change{{Key: "x", Type: Added}})
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestReport_HasDrift_False(t *testing.T) {
	r := makeReport([]Change{{Key: "x", Type: Unchanged}})
	if r.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestReport_WriteText_ContainsHeader(t *testing.T) {
	r := makeReport([]Change{{Key: "TOKEN", Type: Modified, OldValue: "old", NewValue: "new"}})
	var buf bytes.Buffer
	err := r.WriteText(&buf, RenderOptions{Color: false, MaskValues: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "secret/prod/app") {
		t.Error("expected source path in output")
	}
	if !strings.Contains(out, "Drift detected") {
		t.Error("expected drift message in output")
	}
	if !strings.Contains(out, "Summary:") {
		t.Error("expected summary in output")
	}
}
