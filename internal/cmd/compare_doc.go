/*
Package cmd — Compare command

The Compare command diffs two Vault secret paths, optionally at specific
versions. Unlike the Diff command (which compares two versions of the same
path), Compare is designed for cross-path auditing — for example, comparing
a secret in a staging environment against its production counterpart.

Usage example:

	opts := cmd.DefaultCompareOptions()
	opts.PathA = "secret/data/staging/db"
	opts.PathB = "secret/data/prod/db"
	opts.VersionA = 0  // latest
	opts.VersionB = 3
	opts.Mask = true

	if err := cmd.Compare(ctx, client, opts); err != nil {
		log.Fatal(err)
	}

Output format

The output mirrors the unified-diff style used by the Diff command:

	--- secret/data/staging/db (v5)
	+++ secret/data/prod/db (v3)
	+ DB_HOST  : db-prod.internal
	- DB_HOST  : db-staging.internal

Values are masked by default (Mask: true).
*/
package cmd
