package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultPolicyOptions returns PolicyOptions populated from environment variables.
func DefaultPolicyOptions() PolicyOptions {
	return PolicyOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Mount:   mountFromEnv(),
		Token:   os.Getenv("VAULT_TOKEN"),
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// PolicyOptions configures the Policy command.
type PolicyOptions struct {
	Address string
	Mount   string
	Token   string
	Path    string
	Format  string
	Mask    bool
	Output  io.Writer
}

// PolicyResult holds the result of a policy check.
type PolicyResult struct {
	Path     string            `json:"path"`
	Policies map[string]string `json:"policies"`
	Version  int               `json:"version"`
}

// Policy fetches secret metadata and reports which keys match policy patterns.
func Policy(client vault.SecretGetter, opts PolicyOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	secret, err := client.GetSecretVersion(opts.Path, 0)
	if err != nil {
		return fmt.Errorf("policy: fetch %q: %w", opts.Path, err)
	}

	policies := make(map[string]string, len(secret.Data))
	for k := range secret.Data {
		policies[k] = classifyKey(k)
	}

	result := PolicyResult{
		Path:     opts.Path,
		Policies: policies,
		Version:  secret.Version,
	}

	if opts.Format == "json" {
		return json.NewEncoder(out).Encode(result)
	}

	fmt.Fprintf(out, "Path: %s (version %d)\n", result.Path, result.Version)
	for k, p := range result.Policies {
		fmt.Fprintf(out, "  %-30s %s\n", k, p)
	}
	return nil
}

func classifyKey(key string) string {
	lower := strings.ToLower(key)
	switch {
	case strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "token"):
		return "sensitive"
	case strings.Contains(lower, "url") || strings.Contains(lower, "host") || strings.Contains(lower, "addr"):
		return "config"
	default:
		return "general"
	}
}
