package diff

import (
	"testing"
)

func makeChanges() []Change {
	return []Change{
		{Key: "db/password", Type: ChangeTypeAdded, NewValue: "secret"},
		{Key: "db/user", Type: ChangeTypeUnchanged, OldValue: "admin", NewValue: "admin"},
		{Key: "app/token", Type: ChangeTypeModified, OldValue: "old", NewValue: "new"},
		{Key: "app/debug", Type: ChangeTypeRemoved, OldValue: "true"},
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	changes := makeChanges()
	got := Filter(changes, FilterOptions{})
	if len(got) != len(changes) {
		t.Fatalf("expected %d changes, got %d", len(changes), len(got))
	}
}

func TestFilter_ByType_Added(t *testing.T) {
	got := Filter(makeChanges(), FilterOptions{Types: []ChangeType{ChangeTypeAdded}})
	if len(got) != 1 {
		t.Fatalf("expected 1 change, got %d", len(got))
	}
	if got[0].Key != "db/password" {
		t.Errorf("unexpected key %q", got[0].Key)
	}
}

func TestFilter_ByType_MultipleTypes(t *testing.T) {
	got := Filter(makeChanges(), FilterOptions{
		Types: []ChangeType{ChangeTypeAdded, ChangeTypeRemoved},
	})
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
}

func TestFilter_ByPathPrefix(t *testing.T) {
	got := Filter(makeChanges(), FilterOptions{PathPrefix: "app/"})
	if len(got) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(got))
	}
	for _, c := range got {
		if !hasPrefix(c.Key, "app/") {
			t.Errorf("unexpected key %q does not match prefix", c.Key)
		}
	}
}

func TestFilter_ByTypeAndPrefix(t *testing.T) {
	got := Filter(makeChanges(), FilterOptions{
		Types:      []ChangeType{ChangeTypeModified},
		PathPrefix: "app/",
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 change, got %d", len(got))
	}
	if got[0].Key != "app/token" {
		t.Errorf("unexpected key %q", got[0].Key)
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	got := Filter(makeChanges(), FilterOptions{PathPrefix: "nonexistent/"})
	if len(got) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(got))
	}
}
