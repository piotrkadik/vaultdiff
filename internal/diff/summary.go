package diff

import "fmt"

// ChangeType represents the kind of change detected between two secret versions.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Summary holds aggregated counts of changes from a diff operation.
type Summary struct {
	Added     int
	Removed   int
	Modified  int
	Unchanged int
	Total     int
}

// Summarize computes a Summary from a slice of Change values.
func Summarize(changes []Change) Summary {
	var s Summary
	for _, c := range changes {
		s.Total++
		switch c.Type {
		case Added:
			s.Added++
		case Removed:
			s.Removed++
		case Modified:
			s.Modified++
		case Unchanged:
			s.Unchanged++
		}
	}
	return s
}

// HasDrift returns true when any keys were added, removed, or modified.
func (s Summary) HasDrift() bool {
	return s.Added > 0 || s.Removed > 0 || s.Modified > 0
}

// String returns a human-readable one-line representation of the summary.
func (s Summary) String() string {
	return fmt.Sprintf(
		"total=%d added=%d removed=%d modified=%d unchanged=%d drift=%v",
		s.Total, s.Added, s.Removed, s.Modified, s.Unchanged, s.HasDrift(),
	)
}
