package cmd

import (
	"flag"
	"fmt"
	"io"
)

// ParseFlags parses os.Args[1:] into RunOptions.
// It writes usage to out on error.
func ParseFlags(args []string, out io.Writer) (RunOptions, error) {
	fs := flag.NewFlagSet("vaultdiff", flag.ContinueOnError)
	fs.SetOutput(out)

	var opts RunOptions
	fs.StringVar(&opts.Path, "path", "", "secret path (required)")
	fs.IntVar(&opts.VersionA, "a", 0, "first version to compare (required)")
	fs.IntVar(&opts.VersionB, "b", 0, "second version to compare (required)")
	fs.BoolVar(&opts.ShowAll, "all", false, "show unchanged keys")
	fs.BoolVar(&opts.Mask, "mask", true, "mask secret values in output")
	fs.StringVar(&opts.Format, "format", "text", "output format: text or json")
	fs.StringVar(&opts.OutputFile, "out", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return RunOptions{}, err
	}

	if opts.Path == "" {
		return RunOptions{}, fmt.Errorf("-path is required")
	}
	if opts.VersionA <= 0 {
		return RunOptions{}, fmt.Errorf("-a must be a positive version number")
	}
	if opts.VersionB <= 0 {
		return RunOptions{}, fmt.Errorf("-b must be a positive version number")
	}
	if opts.VersionA == opts.VersionB {
		return RunOptions{}, fmt.Errorf("-a and -b must differ")
	}

	return opts, nil
}
