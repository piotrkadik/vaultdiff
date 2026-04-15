package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestRender_ShowsPath(t *testing.T) {
	result := &Result{Path: "secret/app", Changes: []SecretChange{}}
	var buf bytes.Buffer
	Render(&buf, result, FormatOptions{})
	if !strings.Contains(buf.String(), "secret/app") {
		t.Error("expected path in output")
	}
}

func TestRender_HidesUnchangedByDefault(t *testing.T) {
	result := &Result{
		Path: "secret/app",
		Changes: []SecretChange{
			{Key: "k", Type: Unchanged, OldValue: "v", NewValue: "v"},
		},
	}
	var buf bytes.Buffer
	Render(&buf, result, FormatOptions{ShowUnchanged: false})
	if strings.Contains(buf.String(), "\"v\"") {
		t.Error("unchanged value should be hidden")
	}
}

func TestRender_ShowsUnchangedWhenEnabled(t *testing.T) {
	result := &Result{
		Path: "secret/app",
		Changes: []SecretChange{
			{Key: "k", Type: Unchanged, OldValue: "v", NewValue: "v"},
		},
	}
	var buf bytes.Buffer
	Render(&buf, result, FormatOptions{ShowUnchanged: true})
	if !strings.Contains(buf.String(), "\"v\"") {
		t.Error("expected unchanged value in output")
	}
}

func TestRender_MasksValues(t *testing.T) {
	result := &Result{
		Path: "secret/app",
		Changes: []SecretChange{
			{Key: "password", Type: Added, NewValue: "supersecret"},
		},
	}
	var buf bytes.Buffer
	Render(&buf, result, FormatOptions{MaskValues: true})
	if strings.Contains(buf.String(), "supersecret") {
		t.Error("value should be masked")
	}
	if !strings.Contains(buf.String(), "***") {
		t.Error("expected masked placeholder")
	}
}

func TestRender_ColorDisabled(t *testing.T) {
	result := &Result{
		Path: "secret/app",
		Changes: []SecretChange{
			{Key: "k", Type: Modified, OldValue: "a", NewValue: "b"},
		},
	}
	var buf bytes.Buffer
	Render(&buf, result, FormatOptions{ColorEnabled: false})
	if strings.Contains(buf.String(), "\033[") {
		t.Error("ANSI codes should not appear when color is disabled")
	}
}
