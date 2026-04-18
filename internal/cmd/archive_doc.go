// Package cmd — Archive command
//
// Archive captures a point-in-time JSON record of a single Vault secret
// version. Unlike Snapshot (which records the current state), Archive is
// intended for long-term retention: you specify an explicit version number
// and the result is written as a self-contained JSON object that can be
// stored in an audit log, object store, or version-control system.
//
// Example usage:
//
//	err := cmd.Archive(cmd.ArchiveOptions{
//		Address: "https://vault.example.com",
//		Token:   token,
//		Path:    "secret/myapp/prod",
//		Version: 3,
//		Mask:    false,
//		Output:  os.Stdout,
//	})
//
// Output shape:
//
//	{
//	  "path": "secret/myapp/prod",
//	  "version": 3,
//	  "data": { "API_KEY": "abc123" },
//	  "archived_at": "2024-01-15T12:00:00Z"
//	}
package cmd
