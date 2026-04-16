package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// DefaultHistoryOptions returns a HistoryOptions populated with sensible defaults.
func DefaultHistoryOptions() HistoryOptions {
	return HistoryOptions{
		MaxVersions:   10,
		ShowDestroyed: false,
		ShowDeleted:   true,
	}
}

// HistoryOptions controls how version history is retrieved and displayed.
type HistoryOptions struct {
	MaxVersions   int
	ShowDestroyed bool
	ShowDeleted   bool
	Output        io.Writer
}

// VersionEntry is a single record in the history output.
type VersionEntry struct {
	Version   int    `json:"version"`
	CreatedAt string `json:"created_at"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Destroyed bool   `json:"destroyed"`
	Current   bool   `json:"current"`
}

// historyLister is the subset of vault.Client used by PrintHistory.
type historyLister interface {
	ListVersions(path string) ([]historyVersionMeta, error)
}

// historyVersionMeta carries per-version metadata returned by the vault client.
type historyVersionMeta struct {
	Version   int
	CreatedAt string
	DeletedAt string
	Destroyed bool
}

// PrintHistory fetches version metadata for path and writes a JSON array of
// VersionEntry records to opts.Output (defaults to os.Stdout).
func PrintHistory(client historyLister, path string, opts HistoryOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	metas, err := client.ListVersions(path)
	if err != nil {
		return fmt.Errorf("history: list versions: %w", err)
	}

	var entries []VersionEntry
	for i, m := range metas {
		if m.Destroyed && !opts.ShowDestroyed {
			continue
		}
		if m.DeletedAt != "" && !opts.ShowDeleted {
			continue
		}
		isCurrent := i == len(metas)-1
		entries = append(entries, VersionEntry{
			Version:   m.Version,
			CreatedAt: m.CreatedAt,
			DeletedAt: m.DeletedAt,
			Destroyed: m.Destroyed,
			Current:   isCurrent,
		})
		if opts.MaxVersions > 0 && len(entries) >= opts.MaxVersions {
			break
		}
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
