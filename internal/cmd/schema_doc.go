/*
Package cmd — Schema command

Schema fetches a secret version from Vault and prints the key schema: which
keys are present and whether each has a non-empty value.

This is useful for auditing secret shape across environments without exposing
the actual values.

Example usage:

	opts := DefaultSchemaOptions()
	opts.Path    = "myapp/config"
	opts.Version = 3
	opts.Format  = "json"

	client, _ := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err := Schema(client, opts); err != nil {
		log.Fatal(err)
	}
*/
package cmd
