package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

// historyStub satisfies historyLister.
type historyStub struct {
	metas []historyVersionMeta
	err   error
}

func (h *historyStub) ListVersions(_ string) ([]historyVersionMeta, error) {
	return h.metas, h.err
}

func makeHistoryMetas() []historyVersionMeta {
	return []historyVersionMeta{
		{Version: 1, CreatedAt: "2024-01-01T00:00:00Z"},
		{Version: 2, CreatedAt: "2024-01-02T00:00:00Z", DeletedAt: "2024-01-03T00:00:00Z"},
		{Version: 3, CreatedAt: "2024-01-04T00:00:00Z", Destroyed: true},
		{Version: 4, CreatedAt: "2024-01-05T00:00:00Z"},
	}
}

func TestPrintHistory_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	client := &historyStub{metas: makeHistoryMetas()}
	opts := DefaultHistoryOptions()
	opts.ShowDestroyed = true
	opts.Output = &buf

	if err := PrintHistory(client, "secret/app", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entries []VersionEntry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
}

func TestPrintHistory_LastEntryIsCurrent(t *testing.T) {
	var buf bytes.Buffer
	client := &historyStub{metas: makeHistoryMetas()}
	opts := DefaultHistoryOptions()
	opts.ShowDestroyed = true
	opts.ShowDeleted = true
	opts.MaxVersions = 0
	opts.Output = &buf

	if err := PrintHistory(client, "secret/app", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entries []VersionEntry
	json.Unmarshal(buf.Bytes(), &entries) //nolint:errcheck
	last := entries[len(entries)-1]
	if !last.Current {
		t.Errorf("expected last entry to be current, got version %d", last.Version)
	}
}

func TestPrintHistory_FiltersDestroyed(t *testing.T) {
	var buf bytes.Buffer
	client := &historyStub{metas: makeHistoryMetas()}
	opts := DefaultHistoryOptions()
	opts.ShowDestroyed = false
	opts.ShowDeleted = true
	opts.MaxVersions = 0
	opts.Output = &buf

	PrintHistory(client, "secret/app", opts) //nolint:errcheck

	var entries []VersionEntry
	json.Unmarshal(buf.Bytes(), &entries) //nolint:errcheck
	for _, e := range entries {
		if e.Destroyed {
			t.Errorf("destroyed version %d should have been filtered", e.Version)
		}
	}
}

func TestPrintHistory_RespectsMaxVersions(t *testing.T) {
	var buf bytes.Buffer
	client := &historyStub{metas: makeHistoryMetas()}
	opts := DefaultHistoryOptions()
	opts.ShowDestroyed = true
	opts.ShowDeleted = true
	opts.MaxVersions = 2
	opts.Output = &buf

	PrintHistory(client, "secret/app", opts) //nolint:errcheck

	var entries []VersionEntry
	json.Unmarshal(buf.Bytes(), &entries) //nolint:errcheck
	if len(entries) > 2 {
		t.Errorf("expected at most 2 entries, got %d", len(entries))
	}
}

func TestPrintHistory_ErrorPropagated(t *testing.T) {
	client := &historyStub{err: errors.New("vault unavailable")}
	opts := DefaultHistoryOptions()
	opts.Output = &bytes.Buffer{}

	if err := PrintHistory(client, "secret/app", opts); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDefaultHistoryOptions_Defaults(t *testing.T) {
	opts := DefaultHistoryOptions()
	if opts.MaxVersions != 10 {
		t.Errorf("expected MaxVersions=10, got %d", opts.MaxVersions)
	}
	if opts.ShowDestroyed {
		t.Error("expected ShowDestroyed=false")
	}
	if !opts.ShowDeleted {
		t.Error("expected ShowDeleted=true")
	}
}
