// Package diff provides utilities for comparing, filtering, summarizing,
// formatting, and exporting differences between Vault secret versions.
//
// # Export
//
// The Export function serializes a slice of [Change] values into a
// structured format suitable for downstream consumption — for example,
// feeding audit pipelines, storing historical drift records, or integrating
// with external reporting tools.
//
// Supported formats:
//
//   - [ExportJSON] — pretty-printed JSON array of [ExportRecord] objects.
//   - [ExportCSV]  — RFC 4180 CSV with a header row.
//
// Each exported record includes a UTC timestamp stamped at export time,
// the secret path, the change type, and optional old/new values.
// When mask is true, old and new values are replaced with "***" to
// prevent accidental secret exposure in output files.
//
// Example:
//
//	changes := diff.Compare(versionA, versionB)
//	f, _ := os.Create("drift.csv")
//	defer f.Close()
//	if err := diff.Export(f, changes, diff.ExportCSV, true); err != nil {
//		log.Fatal(err)
//	}
package diff
