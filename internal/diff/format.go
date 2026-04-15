package diff

import (
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorGray   = "\033[90m"
)

// FormatOptions controls output rendering.
type FormatOptions struct {
	ShowUnchanged bool
	ColorEnabled  bool
	MaskValues    bool
}

// Render writes a human-readable diff to the provided writer.
func Render(w io.Writer, result *Result, opts FormatOptions) {
	fmt.Fprintf(w, "Path: %s\n", result.Path)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, c := range result.Changes {
		if c.Type == Unchanged && !opts.ShowUnchanged {
			continue
		}
		line := formatChange(c, opts)
		fmt.Fprintln(w, line)
	}
}

func formatChange(c SecretChange, opts FormatOptions) string {
	oldVal := maskIf(c.OldValue, opts.MaskValues)
	newVal := maskIf(c.NewValue, opts.MaskValues)

	var symbol, color, detail string
	switch c.Type {
	case Added:
		symbol, color = "+", colorGreen
		detail = fmt.Sprintf("%s = %q", c.Key, newVal)
	case Removed:
		symbol, color = "-", colorRed
		detail = fmt.Sprintf("%s = %q", c.Key, oldVal)
	case Modified:
		symbol, color = "~", colorYellow
		detail = fmt.Sprintf("%s: %q -> %q", c.Key, oldVal, newVal)
	case Unchanged:
		symbol, color = " ", colorGray
		detail = fmt.Sprintf("%s = %q", c.Key, oldVal)
	}

	if opts.ColorEnabled {
		return fmt.Sprintf("%s%s %s%s", color, symbol, detail, colorReset)
	}
	return fmt.Sprintf("%s %s", symbol, detail)
}

func maskIf(val string, mask bool) string {
	if mask && val != "" {
		return "***"
	}
	return val
}
