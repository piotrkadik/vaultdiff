package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// DefaultTemplateOptions returns TemplateOptions populated from environment variables.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// TemplateOptions controls the behaviour of Template.
type TemplateOptions struct {
	Address  string
	Token    string
	Mount    string
	Path     string
	Version  int
	Template string // Go text/template string; keys referenced as {{.KEY}}
	Mask     bool
	Output   io.Writer
}

// TemplateResult is the JSON-serialisable result of a Template call.
type TemplateResult struct {
	Path     string            `json:"path"`
	Version  int               `json:"version"`
	Rendered string            `json:"rendered"`
	Keys     []string          `json:"keys"`
	Data     map[string]string `json:"data,omitempty"`
}

// Template fetches a secret version and renders a Go template string against its data.
func Template(vc vault.Client, opts TemplateOptions) (*TemplateResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	secret, err := vc.GetSecretVersion(opts.Path, opts.Version)
	if err != nil {
		return nil, fmt.Errorf("template: fetch %s@%d: %w", opts.Path, opts.Version, err)
	}

	keys := make([]string, 0, len(secret.Data))
	for k := range secret.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rendered := renderTemplate(opts.Template, secret.Data)

	result := &TemplateResult{
		Path:     opts.Path,
		Version:  secret.Version,
		Rendered: rendered,
		Keys:     keys,
	}
	if !opts.Mask {
		result.Data = secret.Data
	}

	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintln(opts.Output, string(out))
	return result, nil
}

// renderTemplate performs simple {{.KEY}} substitution without importing text/template
// to keep the dependency surface small. Full template support can be layered later.
func renderTemplate(tmpl string, data map[string]string) string {
	result := tmpl
	for k, v := range data {
		result = strings.ReplaceAll(result, "{{."+k+"}}", v)
	}
	return result
}
