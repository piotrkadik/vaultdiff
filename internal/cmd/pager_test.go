package cmd

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewPager_EmptyBinary_FallsBackToWriter(t *testing.T) {
	var buf bytes.Buffer
	p, err := NewPager("", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Close()
	if p.Writer != &buf {
		t.Error("expected writer to be the provided fallback")
	}
}

func TestNewPager_NotFoundBinary_FallsBackToWriter(t *testing.T) {
	var buf bytes.Buffer
	p, err := NewPager("__nonexistent_pager_binary__", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Close()
	if p.Writer != &buf {
		t.Error("expected writer to be the provided fallback when binary not found")
	}
}

func TestNewPager_NilFallback_UsesStdout(t *testing.T) {
	p, err := NewPager("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Close()
	if p.Writer == nil {
		t.Error("expected non-nil writer when fallback is nil")
	}
}

func TestPager_Close_NilSafe(t *testing.T) {
	var buf bytes.Buffer
	p := &Pager{Writer: &buf}
	if err := p.Close(); err != nil {
		t.Errorf("Close on no-op pager should not error: %v", err)
	}
}

func TestNewPager_WritesToFallback(t *testing.T) {
	var buf bytes.Buffer
	p, err := NewPager("", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer p.Close()
	fmt.Fprint(p.Writer, "hello pager")
	if buf.String() != "hello pager" {
		t.Errorf("expected 'hello pager', got %q", buf.String())
	}
}
