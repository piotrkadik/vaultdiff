// Package cmd provides high-level command functions for the vaultdiff CLI.
//
// # Validate
//
// Validate fetches a specific version of a Vault secret and checks that every
// key in a caller-supplied list is present and non-empty.
//
// It is useful in CI pipelines to assert that a deployment target has all
// required secrets before an application is started.
//
// Basic usage:
//
//	opts := cmd.DefaultValidateOptions()
//	opts.Mask = true // hide values in output (default)
//
//	required := []string{"DB_PASSWORD", "API_TOKEN", "JWT_SECRET"}
//	if err := cmd.Validate("services/api", 3, required, opts); err != nil {
//		log.Fatal(err)
//	}
//
// Exit behaviour:
//
// Validate returns a non-nil error when one or more required keys are absent or
// empty, or when the secret cannot be fetched. The caller is responsible for
// translating this into a process exit code.
package cmd
