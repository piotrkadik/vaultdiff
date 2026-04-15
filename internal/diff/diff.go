package diff

import (
	"fmt"
	"sort"
)

// ChangeType represents the type of change between two secret versions.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// SecretChange represents a single key-level change between two secret versions.
type SecretChange struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff result between two secret maps.
type Result struct {
	Changes []SecretChange
	Path    string
}

// HasChanges returns true if any non-unchanged entries exist.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Compare diffs two secret data maps and returns a Result.
func Compare(path string, oldData, newData map[string]interface{}) *Result {
	result := &Result{Path: path}

	allKeys := mergeKeys(oldData, newData)
	for _, key := range allKeys {
		oldVal, oldOk := oldData[key]
		newVal, newOk := newData[key]

		switch {
		case !oldOk && newOk:
			result.Changes = append(result.Changes, SecretChange{
				Key: key, Type: Added,
				NewValue: fmt.Sprintf("%v", newVal),
			})
		case oldOk && !newOk:
			result.Changes = append(result.Changes, SecretChange{
				Key: key, Type: Removed,
				OldValue: fmt.Sprintf("%v", oldVal),
			})
		case fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal):
			result.Changes = append(result.Changes, SecretChange{
				Key: key, Type: Modified,
				OldValue: fmt.Sprintf("%v", oldVal),
				NewValue: fmt.Sprintf("%v", newVal),
			})
		default:
			result.Changes = append(result.Changes, SecretChange{
				Key: key, Type: Unchanged,
				OldValue: fmt.Sprintf("%v", oldVal),
				NewValue: fmt.Sprintf("%v", newVal),
			})
		}
	}
	return result
}

func mergeKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
