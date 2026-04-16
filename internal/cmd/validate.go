package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// ValidateOptions controls the behaviour of Validate.
type ValidateOptions struct {
	Address string
	Token   string
	Mount   string
	Mask    bool
	Output  io.Writer
}

// DefaultValidateOptions returns options populated from environment variables.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    true,
		Output:  os.Stdout,
	}
}

// Validate fetches the given secret path and version and reports whether all
// keys are present and non-empty. It writes a human-readable report to
// opts.Output and returns an error if any key is blank or the fetch fails.
func Validate(path string, version int, requiredKeys []string, opts ValidateOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	data, err := vault.GetSecretVersion(client, path, version)
	if err != nil {
		return fmt.Errorf("validate: fetch %s@%d: %w", path, version, err)
	}

	fmt.Fprintf(opts.Output, "Validating %s (version %d)\n", path, version)

	var missing []string
	for _, k := range requiredKeys {
		v, ok := data[k]
		if !ok || v == "" {
			missing = append(missing, k)
			fmt.Fprintf(opts.Output, "  MISSING  %s\n", k)
		} else {
			display := v
			if opts.Mask {
				display = "***"
			}
			fmt.Fprintf(opts.Output, "  OK       %s = %s\n", k, display)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("validate: %d required key(s) missing or empty: %v", len(missing), missing)
	}
	fmt.Fprintln(opts.Output, "All required keys present.")
	return nil
}

func mountFromEnv() string {
	if m := os.Getenv("VAULT_MOUNT"); m != "" {
		return m
	}
	return "secret"
}
