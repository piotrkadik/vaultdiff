package cmd

// Mirror copies all key-value pairs from a source secret path to a destination
// secret path within the same Vault cluster.
//
// It reads the latest version of the source secret and writes its data to the
// destination path as a new version. The original source secret is not modified.
//
// When DryRun is true the write is skipped and the result is still returned,
// allowing callers to preview what would be written.
//
// Example:
//
//	err := cmd.Mirror(cmd.MirrorOptions{
//		Address: "http://127.0.0.1:8200",
//		Token:   "root",
//		Mount:   "secret",
//		SrcPath: "myapp/production",
//		DstPath: "myapp/staging",
//		DryRun:  false,
//		Mask:    true,
//	})
