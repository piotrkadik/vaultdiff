package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/vaultdiff/internal/diff"
)

// PromptOptions controls interactive prompt behaviour.
type PromptOptions struct {
	In  io.Reader
	Out io.Writer
}

// ConfirmDrift asks the user whether to proceed when drift is detected.
// Returns true if the user confirms, false otherwise.
func ConfirmDrift(report *diff.Report, opts PromptOptions) (bool, error) {
	if !report.HasDrift() {
		return true, nil
	}

	summary := report.Summary
	fmt.Fprintf(opts.Out, "\nDrift detected: +%d added, -%d removed, ~%d modified\n",
		summary.Added, summary.Removed, summary.Modified)
	fmt.Fprint(opts.Out, "Proceed? [y/N]: ")

	var answer string
	_, err := fmt.Fscan(opts.In, &answer)
	if err != nil {
		return false, fmt.Errorf("reading prompt response: %w", err)
	}

	return strings.EqualFold(strings.TrimSpace(answer), "y"), nil
}
