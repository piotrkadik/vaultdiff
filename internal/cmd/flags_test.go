package cmd_test

import (
	"bytes"
	"testing"

	"github.com/example/vaultdiff/internal/cmd"
)

func TestParseFlags_ValidArgs(t *testing.T) {
	opts, err := cmd.ParseFlags([]string{"-path", "myapp", "-a", "1", "-b", "2"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Path != "myapp" {
		t.Errorf("expected path myapp, got %s", opts.Path)
	}
	if opts.VersionA != 1 || opts.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", opts.VersionA, opts.VersionB)
	}
}

func TestParseFlags_MissingPath(t *testing.T) {
	_, err := cmd.ParseFlags([]string{"-a", "1", "-b", "2"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing -path")
	}
}

func TestParseFlags_SameVersions(t *testing.T) {
	_, err := cmd.ParseFlags([]string{"-path", "x", "-a", "3", "-b", "3"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error when -a == -b")
	}
}

func TestParseFlags_ZeroVersion(t *testing.T) {
	_, err := cmd.ParseFlags([]string{"-path", "x", "-a", "0", "-b", "2"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for zero version")
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	opts, err := cmd.ParseFlags([]string{"-path", "svc", "-a", "1", "-b", "2"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "text" {
		t.Errorf("expected default format text, got %s", opts.Format)
	}
	if !opts.Mask {
		t.Error("expected mask to default to true")
	}
}
