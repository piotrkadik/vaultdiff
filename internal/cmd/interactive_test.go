package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/diff"
)

func makeReport(added, removed, modified int) *diff.Report {
	changes := make([]diff.Change, 0, added+removed+modified)
	for i := 0; i < added; i++ {
		changes = append(changes, diff.Change{Type: diff.Added, Path: "a"})
	}
	for i := 0; i < removed; i++ {
		changes = append(changes, diff.Change{Type: diff.Removed, Path: "r"})
	}
	for i := 0; i < modified; i++ {
		changes = append(changes, diff.Change{Type: diff.Modified, Path: "m"})
	}
	return diff.NewReport("secret/test", 1, 2, changes)
}

func TestConfirmDrift_NoDrift_SkipsPrompt(t *testing.T) {
	report := makeReport(0, 0, 0)
	var out bytes.Buffer
	ok, err := ConfirmDrift(report, PromptOptions{In: strings.NewReader(""), Out: &out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true when no drift")
	}
	if out.Len() != 0 {
		t.Errorf("expected no output, got: %q", out.String())
	}
}

func TestConfirmDrift_UserConfirms(t *testing.T) {
	report := makeReport(1, 0, 0)
	var out bytes.Buffer
	ok, err := ConfirmDrift(report, PromptOptions{In: strings.NewReader("y\n"), Out: &out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true when user answers y")
	}
}

func TestConfirmDrift_UserDeclines(t *testing.T) {
	report := makeReport(0, 1, 0)
	var out bytes.Buffer
	ok, err := ConfirmDrift(report, PromptOptions{In: strings.NewReader("n\n"), Out: &out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when user answers n")
	}
}

func TestConfirmDrift_OutputContainsCounts(t *testing.T) {
	report := makeReport(2, 1, 3)
	var out bytes.Buffer
	_, _ = ConfirmDrift(report, PromptOptions{In: strings.NewReader("n\n"), Out: &out})
	got := out.String()
	for _, want := range []string{"+2", "-1", "~3"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected output to contain %q, got: %q", want, got)
		}
	}
}

func TestConfirmDrift_CaseInsensitive(t *testing.T) {
	report := makeReport(1, 0, 0)
	var out bytes.Buffer
	ok, err := ConfirmDrift(report, PromptOptions{In: strings.NewReader("Y\n"), Out: &out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for uppercase Y")
	}
}

func TestConfirmDrift_EmptyInput_ReturnsFalse(t *testing.T) {
	// An empty reader simulates EOF without any user input (e.g. non-interactive pipe).
	// ConfirmDrift should treat this as a decline rather than an error.
	report := makeReport(1, 0, 0)
	var out bytes.Buffer
	ok, err := ConfirmDrift(report, PromptOptions{In: strings.NewReader(""), Out: &out})
	if err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if ok {
		t.Error("expected false when input is empty (EOF)")
	}
}
