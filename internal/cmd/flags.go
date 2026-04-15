package cmd

import (
	"errors"
	"flag"
	"fmt"
)

// Flags holds the parsed CLI arguments for a diff invocation.
type Flags struct {
	Path       string
	VersionA   int
	VersionB   int
	Mount      string
	Format     string
	Mask       bool
	ShowAll    bool
	Snapshot   bool
}

// ParseFlags parses os.Args[1:] using the provided FlagSet and returns Flags.
func ParseFlags(fs *flag.FlagSet, args []string) (Flags, error) {
	var f Flags

	fs.StringVar(&f.Path, "path", "", "Vault secret path (required)")
	fs.IntVar(&f.VersionA, "version-a", 0, "First version to compare (required, >0)")
	fs.IntVar(&f.VersionB, "version-b", 0, "Second version to compare (required, >0)")
	fs.StringVar(&f.Mount, "mount", "secret", "KV v2 mount path")
	fs.StringVar(&f.Format, "format", "text", "Output format: text, json, csv")
	fs.BoolVar(&f.Mask, "mask", true, "Mask secret values in output")
	fs.BoolVar(&f.ShowAll, "show-all", false, "Include unchanged keys in output")
	fs.BoolVar(&f.Snapshot, "snapshot", false, "Write a JSON snapshot instead of a diff")

	if err := fs.Parse(args); err != nil {
		return Flags{}, fmt.Errorf("flags: %w", err)
	}

	if f.Path == "" {
		return Flags{}, errors.New("flags: -path is required")
	}
	if f.VersionA <= 0 || f.VersionB <= 0 {
		return Flags{}, errors.New("flags: -version-a and -version-b must be greater than zero")
	}
	if f.VersionA == f.VersionB {
		return Flags{}, errors.New("flags: -version-a and -version-b must differ")
	}

	return f, nil
}
