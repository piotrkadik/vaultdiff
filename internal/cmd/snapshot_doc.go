// Package cmd provides the top-level command implementations for vaultdiff.
//
// # Snapshot
//
// The snapshot feature captures a point-in-time view of a single Vault secret
// version and serialises it as a JSON document. It is useful for:
//
//   - Archiving the state of a secret before a deployment.
//   - Providing a baseline for future diff operations.
//   - Auditing which keys existed at a specific version.
//
// Values can optionally be masked (replaced with "***") so that the snapshot
// file is safe to commit to a repository or share with teammates who should
// not see plaintext credentials.
//
// Example output:
//
//	{
//	  "path": "secret/data/myapp",
//	  "version": 3,
//	  "captured_at": "2024-01-15T10:30:00Z",
//	  "data": {
//	    "DB_PASS": "***",
//	    "API_KEY": "***"
//	  }
//	}
package cmd
