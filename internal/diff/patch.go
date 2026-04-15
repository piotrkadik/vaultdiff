package diff

import (
	"encoding/json"
	"fmt"
	"io"
)

// PatchFormat controls the serialization format for a patch.
type PatchFormat string

const (
	PatchFormatJSON PatchFormat = "json"
	PatchFormatText PatchFormat = "text"
)

// Patch represents a serializable set of changes that can be replayed or
// stored for later auditing.
type Patch struct {
	Path    string   `json:"path"`
	From    int      `json:"from_version"`
	To      int      `json:"to_version"`
	Changes []Change `json:"changes"`
}

// NewPatch constructs a Patch from a slice of Change values.
func NewPatch(path string, from, to int, changes []Change) Patch {
	return Patch{
		Path:    path,
		From:    from,
		To:      to,
		Changes: changes,
	}
}

// WritePatch serializes the patch to w in the requested format.
func WritePatch(w io.Writer, p Patch, format PatchFormat) error {
	switch format {
	case PatchFormatJSON:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(p)
	case PatchFormatText:
		return writePatchText(w, p)
	default:
		return fmt.Errorf("diff: unknown patch format %q", format)
	}
}

func writePatchText(w io.Writer, p Patch) error {
	_, err := fmt.Fprintf(w, "patch: %s  v%d -> v%d  (%d change(s))\n",
		p.Path, p.From, p.To, len(p.Changes))
	if err != nil {
		return err
	}
	for _, c := range p.Changes {
		symbol := changeSymbol(c.Type)
		_, err = fmt.Fprintf(w, "  %s %s\n", symbol, c.Key)
		if err != nil {
			return err
		}
	}
	return nil
}

func changeSymbol(t ChangeType) string {
	switch t {
	case ChangeAdded:
		return "+"
	case ChangeRemoved:
		return "-"
	case ChangeModified:
		return "~"
	default:
		return " "
	}
}
