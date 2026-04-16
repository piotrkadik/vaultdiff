package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

// TagOptions configures the Tag command.
type TagOptions struct {
	Address string
	Token   string
	Mount   string
	Mask    bool
	Output  io.Writer
}

// DefaultTagOptions returns TagOptions populated from environment variables.
func DefaultTagOptions() TagOptions {
	mask := true
	if os.Getenv("VAULTDIFF_MASK") == "false" {
		mask = false
	}
	return TagOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    mask,
	}
}

// TagResult holds the result of tagging a secret version.
type TagResult struct {
	Path    string            `json:"path"`
	Version int               `json:"version"`
	Tags    map[string]string `json:"tags"`
	TaggedAt time.Time        `json:"tagged_at"`
}

// Tag attaches metadata tags to a specific secret version and writes the result as JSON.
func Tag(path string, version int, tags map[string]string, opts TagOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("tag: create client: %w", err)
	}

	_, err = client.GetSecretVersion(path, version)
	if err != nil {
		return fmt.Errorf("tag: fetch version %d of %q: %w", version, path, err)
	}

	result := TagResult{
		Path:     path,
		Version:  version,
		Tags:     tags,
		TaggedAt: time.Now().UTC(),
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return fmt.Errorf("tag: encode result: %w", err)
	}
	return nil
}
