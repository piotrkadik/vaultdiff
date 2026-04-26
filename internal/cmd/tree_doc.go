// Package cmd provides CLI command implementations for vaultdiff.
//
// # Tree
//
// Tree fetches all keys stored at a given Vault KV path and renders them
// in a tree-like structure, similar to the Unix `tree` command.
//
// Text output example:
//
//	secret/myapp
//	├── database_url
//	├── jwt_secret
//	└── redis_url
//
//	3 key(s)
//
// JSON output includes the path, a sorted list of keys, and a count:
//
//	{
//	  "path": "secret/myapp",
//	  "keys": ["database_url", "jwt_secret", "redis_url"],
//	  "count": 3
//	}
//
// Set VAULTDIFF_FORMAT=json to switch output format, or pass
// TreeOptions.Format directly when calling Tree programmatically.
package cmd
