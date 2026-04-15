package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/vaultdiff/internal/diff"
)

// OutputOptions controls how diff results are written to the terminal or file.
type OutputOptions struct {
	Writer    io.Writer
	Format    string // "text", "json", "csv"
	Mask      bool
	ShowAll   bool
	Color     bool
	PatchFile string
}

// DefaultOutputOptions returns sensible defaults for terminal output.
func DefaultOutputOptions() OutputOptions {
	return OutputOptions{
		Writer:  os.Stdout,
		Format:  "text",
		Mask:    true,
		ShowAll: false,
		Color:   true,
	}
}

// WriteOutput renders the diff result according to the chosen format and options.
func WriteOutput(changes []diff.Change, opts OutputOptions) error {
	switch opts.Format {
	case "json", "csv":
		return diff.Export(changes, opts.Format, opts.Mask, opts.Writer)
	case "text", "":
		return diff.Render(changes, opts.Writer, diff.RenderOptions{
			Mask:    opts.Mask,
			ShowAll: opts.ShowAll,
			Color:   opts.Color,
		})
	default:
		return fmt.Errorf("unknown output format: %q", opts.Format)
	}
}
