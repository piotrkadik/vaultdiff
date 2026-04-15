package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]interface{}{}
	new := map[string]interface{}{"key": "value"}

	result := Compare("secret/test", old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]interface{}{"key": "value"}
	new := map[string]interface{}{}

	result := Compare("secret/test", old, new)
	if result.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].OldValue != "value" {
		t.Errorf("unexpected OldValue: %s", result.Changes[0].OldValue)
	}
}

func TestCompare_Modified(t *testing.T) {
	old := map[string]interface{}{"key": "old"}
	new := map[string]interface{}{"key": "new"}

	result := Compare("secret/test", old, new)
	if result.Changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Unchanged(t *testing.T) {
	old := map[string]interface{}{"key": "same"}
	new := map[string]interface{}{"key": "same"}

	result := Compare("secret/test", old, new)
	if result.Changes[0].Type != Unchanged {
		t.Errorf("expected Unchanged, got %s", result.Changes[0].Type)
	}
	if result.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestCompare_KeyOrder(t *testing.T) {
	old := map[string]interface{}{"b": "1", "a": "2"}
	new := map[string]interface{}{"b": "1", "a": "2"}

	result := Compare("secret/test", old, new)
	if result.Changes[0].Key != "a" {
		t.Errorf("expected first key to be 'a', got '%s'", result.Changes[0].Key)
	}
}

func TestResult_Path(t *testing.T) {
	result := Compare("secret/myapp", nil, nil)
	if result.Path != "secret/myapp" {
		t.Errorf("unexpected path: %s", result.Path)
	}
}
