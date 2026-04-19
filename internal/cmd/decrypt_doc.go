// Package cmd provides high-level CLI command implementations for vaultdiff.
//
// # Decrypt
//
// Decrypt fetches a specific version of a Vault KV-v2 secret and writes its
// key/value pairs to the configured output in either text or JSON format.
//
// Text output (default):
//
//	# secret/myapp version 3
//	DB_PASSWORD=hunter2
//	API_KEY=abc123
//
// JSON output (VAULTDIFF_FORMAT=json or Format:"json"):
//
//	{"path":"secret/myapp","version":3,"data":{"DB_PASSWORD":"hunter2"}}
//
// Values are masked by default (Mask:true). Set VAULTDIFF_MASK=false or
// Mask:false to reveal plaintext values.
package cmd
