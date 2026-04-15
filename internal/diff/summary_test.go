package diff

import (
	"strings"
	"testing"
)

func makeSummaryChanges() []Change {
	return []Change{
		{Key: "DB_HOST", Type: Added, NewValue: "localhost"},
		{Key: "DB_PORT", Type: Removed, OldValue: "5432"},
		{Key: "DB_PASS", Type: Modified, OldValue: "old", NewValue: "new"},
		{Key: "APP_ENV", Type: Unchanged, OldValue: "prod", NewValue: "prod"},
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := Summarize(makeSummaryChanges())

	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected Modified=1, got %d", s.Modified)
	}
	if s.Unchanged != 1 {
		t.Errorf("expected Unchanged=1, got %d", s.Unchanged)
	}
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 || s.HasDrift() {
		t.Errorf("expected empty summary with no drift, got %+v", s)
	}
}

func TestSummary_HasDrift_True(t *testing.T) {
	s := Summarize(makeSummaryChanges())
	if !s.HasDrift() {
		t.Error("expected HasDrift=true")
	}
}

func TestSummary_HasDrift_False(t *testing.T) {
	changes := []Change{
		{Key: "APP_ENV", Type: Unchanged, OldValue: "prod", NewValue: "prod"},
	}
	s := Summarize(changes)
	if s.HasDrift() {
		t.Error("expected HasDrift=false")
	}
}

func TestSummary_String_ContainsFields(t *testing.T) {
	s := Summarize(makeSummaryChanges())
	out := s.String()

	for _, want := range []string{"total=", "added=", "removed=", "modified=", "unchanged=", "drift="} {
		if !strings.Contains(out, want) {
			t.Errorf("String() missing field %q in %q", want, out)
		}
	}
}
