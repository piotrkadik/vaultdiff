package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultHealthOptions returns HealthOptions populated from environment variables.
func DefaultHealthOptions() HealthOptions {
	return HealthOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// HealthOptions configures the Health command.
type HealthOptions struct {
	Address string
	Mount   string
	Format  string
	Output  io.Writer
}

// HealthResult holds the outcome of a vault connectivity check.
type HealthResult struct {
	Reachable bool      `json:"reachable"`
	Address   string    `json:"address"`
	Mount     string    `json:"mount"`
	CheckedAt time.Time `json:"checked_at"`
	Error     string    `json:"error,omitempty"`
}

// Health checks whether the configured Vault address is reachable and the
// KV mount is accessible, then writes a result to opts.Output.
func Health(path string, opts HealthOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	result := HealthResult{
		Address:   opts.Address,
		Mount:     opts.Mount,
		CheckedAt: time.Now().UTC(),
	}

	client, err := vault.NewClient(opts.Address, opts.Mount)
	if err != nil {
		result.Error = err.Error()
		return writeHealthResult(out, opts.Format, result)
	}

	_, err = client.GetSecretVersion(path, 0)
	if err != nil {
		result.Error = err.Error()
		return writeHealthResult(out, opts.Format, result)
	}

	result.Reachable = true
	return writeHealthResult(out, opts.Format, result)
}

func writeHealthResult(w io.Writer, format string, r HealthResult) error {
	if format == "json" {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	}
	status := "OK"
	if !r.Reachable {
		status = "UNREACHABLE"
	}
	fmt.Fprintf(w, "vault health: %s\n", status)
	fmt.Fprintf(w, "  address : %s\n", r.Address)
	fmt.Fprintf(w, "  mount   : %s\n", r.Mount)
	if r.Error != "" {
		fmt.Fprintf(w, "  error   : %s\n", r.Error)
	}
	return nil
}
