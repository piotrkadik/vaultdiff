package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

// TouchOptions configures the Touch command.
type TouchOptions struct {
	Address string
	Token   string
	Mount   string
	Mask    bool
	DryRun  bool
	Output  io.Writer
}

// DefaultTouchOptions returns options populated from environment variables.
func DefaultTouchOptions() TouchOptions {
	return TouchOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// TouchResult holds the outcome of a Touch operation.
type TouchResult struct {
	Path      string         `json:"path"`
	Version   int            `json:"version"`
	TouchedAt time.Time      `json:"touched_at"`
	DryRun    bool           `json:"dry_run"`
	Data      map[string]any `json:"data,omitempty"`
}

// Touch re-writes the current secret version at path to create a new version
// with identical data, useful for bumping the version counter or audit trail.
func Touch(path string, opts TouchOptions) (*TouchResult, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return nil, fmt.Errorf("touch: create client: %w", err)
	}

	versions, err := vault.ListVersions(client, path)
	if err != nil {
		return nil, fmt.Errorf("touch: list versions: %w", err)
	}
	latest := vault.LatestVersion(versions)

	secret, err := vault.GetSecretVersion(client, path, latest)
	if err != nil {
		return nil, fmt.Errorf("touch: fetch secret: %w", err)
	}

	result := &TouchResult{
		Path:      path,
		Version:   latest + 1,
		TouchedAt: time.Now().UTC(),
		DryRun:    opts.DryRun,
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, path, secret.Data); err != nil {
			return nil, fmt.Errorf("touch: write secret: %w", err)
		}
		if !opts.Mask {
			result.Data = secret.Data
		}
	}

	return result, json.NewEncoder(opts.Output).Encode(result)
}
