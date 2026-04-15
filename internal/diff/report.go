package diff

import (
	"fmt"
	"io"
	"time"
)

// Report holds a summary of a diff operation between two secret paths.
type Report struct {
	SourcePath  string
	TargetPath  string
	SourceVersion int
	TargetVersion int
	GeneratedAt time.Time
	Changes     []Change
}

// Summary returns counts of each change type.
func (r *Report) Summary() map[ChangeType]int {
	counts := map[ChangeType]int{
		Added:     0,
		Removed:   0,
		Modified:  0,
		Unchanged: 0,
	}
	for _, c := range r.Changes {
		counts[c.Type]++
	}
	return counts
}

// HasDrift returns true if any keys differ between source and target.
func (r *Report) HasDrift() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// WriteText writes a human-readable report to w.
func (r *Report) WriteText(w io.Writer, opts RenderOptions) error {
	fmt.Fprintf(w, "Vault Secret Diff Report\n")
	fmt.Fprintf(w, "Generated: %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Source: %s (v%d)\n", r.SourcePath, r.SourceVersion)
	fmt.Fprintf(w, "Target: %s (v%d)\n", r.TargetPath, r.TargetVersion)
	fmt.Fprintf(w, "\n")

	if err := Render(w, r.Changes, opts); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	summary := r.Summary()
	fmt.Fprintf(w, "\nSummary: +%d added  -%d removed  ~%d modified  =%d unchanged\n",
		summary[Added], summary[Removed], summary[Modified], summary[Unchanged])

	if r.HasDrift() {
		fmt.Fprintf(w, "Drift detected.\n")
	} else {
		fmt.Fprintf(w, "No drift detected.\n")
	}
	return nil
}
