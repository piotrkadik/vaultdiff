// Package audit provides structured, append-only audit logging for the
// vaultdiff tool.
//
// Each diff operation produces an [Entry] that captures the path compared,
// the two secret versions involved, a summary of changes (added, removed,
// modified counts), and whether drift was detected.
//
// Entries are serialised as newline-delimited JSON (NDJSON) so that logs can
// be streamed to a file, stdout, or any io.Writer and later processed by
// standard tooling such as jq or a SIEM.
//
// Basic usage:
//
//	f, _ := os.OpenFile("vaultdiff-audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
//	defer f.Close()
//
//	l := audit.NewLogger(f)
//	l.Write(audit.Entry{
//		Environment: "production",
//		Path:        "secret/data/myapp",
//		VersionA:    3,
//		VersionB:    4,
//		Added:       0,
//		Removed:     1,
//		Modified:    2,
//		HasDrift:    true,
//	})
package audit
