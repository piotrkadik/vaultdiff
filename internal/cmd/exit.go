package cmd

// Exit codes used by the vaultdiff CLI.
const (
	// ExitOK indicates a successful run with no drift detected.
	ExitOK = 0

	// ExitDrift indicates the run completed but drift was detected.
	ExitDrift = 1

	// ExitError indicates a fatal error prevented the diff from completing.
	ExitError = 2

	// ExitCancelled indicates the user cancelled an interactive prompt.
	ExitCancelled = 3
)

// ExitCodeForDrift returns ExitDrift when drift is present, ExitOK otherwise.
func ExitCodeForDrift(hasDrift bool) int {
	if hasDrift {
		return ExitDrift
	}
	return ExitOK
}
