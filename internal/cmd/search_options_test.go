package cmd_test

import (
	"os"
	"testing"

	"github.com/vaultdiff/internal/cmd"
)

func TestDefaultSearchOptions_FormatFromEnv(t *testing.T) {
	os.Setenv("VAULTDIFF_FORMAT", "json")
	t.Cleanup(func() { os.Unsetenv("VAULTDIFF_FORMAT") })

	opts := cmd.DefaultSearchOptions()
	if opts.Format != "json" {
		t.Errorf("expected json, got %s", opts.Format)
	}
}

func TestDefaultSearchOptions_MaskDefaultsTrue(t *testing.T) {
	os.Unsetenv("VAULTDIFF_MASK")
	opts := cmd.DefaultSearchOptions()
	if !opts.Mask {
		t.Error("expected mask=true by default")
	}
}

func TestDefaultSearchOptions_MaskFalseWhenDisabled(t *testing.T) {
	os.Setenv("VAULTDIFF_MASK", "false")
	t.Cleanup(func() { os.Unsetenv("VAULTDIFF_MASK") })

	opts := cmd.DefaultSearchOptions()
	if opts.Mask {
		t.Error("expected mask=false when env is 'false'")
	}
}

func TestSearchOptions_ZeroValue(t *testing.T) {
	var opts cmd.SearchOptions
	if opts.Query != "" {
		t.Error("expected empty query")
	}
	if opts.Output != nil {
		t.Error("expected nil output")
	}
}
