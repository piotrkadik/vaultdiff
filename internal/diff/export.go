package diff

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ExportFormat defines the output format for exported diff results.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)

// ExportRecord represents a single exportable diff entry.
type ExportRecord struct {
	Timestamp time.Time  `json:"timestamp"`
	Path      string     `json:"path"`
	ChangeType ChangeType `json:"change_type"`
	OldValue  string     `json:"old_value,omitempty"`
	NewValue  string     `json:"new_value,omitempty"`
}

// Export writes the given changes to w in the specified format.
func Export(w io.Writer, changes []Change, format ExportFormat, mask bool) error {
	records := toRecords(changes, mask)
	switch format {
	case ExportJSON:
		return exportJSON(w, records)
	case ExportCSV:
		return exportCSV(w, records)
	default:
		return fmt.Errorf("unknown export format: %q", format)
	}
}

func toRecords(changes []Change, mask bool) []ExportRecord {
	now := time.Now().UTC()
	records := make([]ExportRecord, 0, len(changes))
	for _, c := range changes {
		r := ExportRecord{
			Timestamp:  now,
			Path:       c.Path,
			ChangeType: c.Type,
			OldValue:   maskIf(c.OldValue, mask),
			NewValue:   maskIf(c.NewValue, mask),
		}
		records = append(records, r)
	}
	return records
}

func exportJSON(w io.Writer, records []ExportRecord) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func exportCSV(w io.Writer, records []ExportRecord) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "path", "change_type", "old_value", "new_value"}); err != nil {
		return err
	}
	for _, r := range records {
		row := []string{
			r.Timestamp.Format(time.RFC3339),
			r.Path,
			string(r.ChangeType),
			r.OldValue,
			r.NewValue,
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
