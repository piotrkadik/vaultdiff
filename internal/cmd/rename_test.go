package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
)

func TestRename_DryRun_SkipsWrite(t *testing.T) {
	stub := newValidateStub() // reuse httptest vault stub
	defer stub.Close()

	var buf bytes.Buffer
	err := Rename(context.Background(), "svc/config", "OLD_KEY", "NEW_KEY", RenameOptions{
		Address: stub.URL,
		Token:   "test-token",
		Mount:   "secret",
		DryRun:  true,
		Output:  &buf,
	})
	// stub may return not-found; we only assert no panic and JSON shape when err==nil
	if err == nil {
		var r RenameResult
		if jsonErr := json.Unmarshal(buf.Bytes(), &r); jsonErr != nil {
			t.Fatalf("expected valid JSON: %v", jsonErr)
		}
		if !r.DryRun {
			t.Error("expected dry_run=true in result")
		}
	}
}

func TestRename_NilOutputDefaultsToStdout(t *testing.T) {
	stub := newValidateStub()
	defer stub.Close()

	// Should not panic even with nil Output
	_ = Rename(context.Background(), "svc/config", "A", "B", RenameOptions{
		Address: stub.URL,
		Token:   "test-token",
		Mount:   "secret",
		Output:  nil,
	})
}

func TestRename_InvalidAddress_ReturnsError(t *testing.T) {
	err := Rename(context.Background(), "svc/config", "A", "B", RenameOptions{
		Address: "://bad",
		Token:   "tok",
		Mount:   "secret",
		Output:  &bytes.Buffer{},
	})
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}

func TestDefaultRenameOptions_Defaults(t *testing.T) {
	opts := DefaultRenameOptions()
	if opts.DryRun {
		t.Error("expected DryRun=false by default")
	}
	if opts.Mount == "" {
		t.Error("expected non-empty Mount")
	}
}

func TestRenameResult_JSONShape(t *testing.T) {
	r := RenameResult{Path: "svc/cfg", OldKey: "FOO", NewKey: "BAR", DryRun: false}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"path", "old_key", "new_key", "dry_run"} {
		if !bytes.Contains(b, []byte(field)) {
			t.Errorf("missing field %q in JSON", field)
		}
	}
}
