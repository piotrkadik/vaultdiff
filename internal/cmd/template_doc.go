// Package cmd — Template command
//
// Template fetches a specific version of a Vault KV-v2 secret and renders a
// caller-supplied Go-style template string by substituting {{.KEY}} placeholders
// with the corresponding secret values.
//
// # Usage
//
//	opts := cmd.DefaultTemplateOptions()
//	opts.Path    = "infra/database"
//	opts.Version = 0          // 0 = latest
//	opts.Template = "host={{.DB_HOST}} port={{.DB_PORT}}"
//
//	client, _ := vault.NewClient(opts.Address, opts.Token, opts.Mount)
//	result, err := cmd.Template(client, opts)
//
// # Output
//
// Template always writes a JSON object to opts.Output (defaults to os.Stdout):
//
//	{
//	  "path":     "infra/database",
//	  "version":  4,
//	  "rendered": "host=db.internal port=5432",
//	  "keys":     ["DB_HOST", "DB_PORT"],
//	  "data":     { "DB_HOST": "db.internal", "DB_PORT": "5432" }
//	}
//
// When opts.Mask is true (the default) the "data" field is omitted so that
// raw secret values are never written to the output stream.
package cmd
