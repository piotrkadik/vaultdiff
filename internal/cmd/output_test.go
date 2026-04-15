package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/diff"
)

func makeOutputChanges() []diff.Change {
	return []diff.Change{
		{Path: "secret/app", Key: "API_KEY", Type: diff.Added, NewValue: "abc123"},
		{Path: "secret/app", Key: "DB_PASS", Type: diff.Removed, OldValue: "hunter2"},
		{Path: "secret/app", Key: "HOST", Type: diff.Unchanged, OldValue: "localhost", NewValue: "localhost"},
	}
}

func TestWriteOutput_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultOutputOptions()
	opts.Writer = &buf
	opts.Color = false

	if err := WriteOutput(makeOutputChanges(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in text output, got: %s", out)
	}
}

func TestWriteOutput_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultOutputOptions()
	opts.Writer = &buf
	opts.Format = "json"
	opts.Mask = false

	if err := WriteOutput(makeOutputChanges(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in JSON output, got: %s", out)
	}
}

func TestWriteOutput_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultOutputOptions()
	opts.Writer = &buf
	opts.Format = "csv"
	opts.Mask = false

	if err := WriteOutput(makeOutputChanges(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "path") {
		t.Errorf("expected CSV header in output, got: %s", out)
	}
}

func TestWriteOutput_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := DefaultOutputOptions()
	opts.Writer = &buf
	opts.Format = "xml"

	if err := WriteOutput(makeOutputChanges(), opts); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestDefaultOutputOptions_Defaults(t *testing.T) {
	opts := DefaultOutputOptions()
	if opts.Format != "text" {
		t.Errorf("expected format=text, got %q", opts.Format)
	}
	if !opts.Mask {
		t.Error("expected Mask=true by default")
	}
	if !opts.Color {
		t.Error("expected Color=true by default")
	}
}
