package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/example/vaultdiff/internal/vault"
)

// CopyOptions configures the Copy command.
type CopyOptions struct {
	Address string
	Token   string
	Mount   string
	DryRun  bool
	Mask    bool
	Output  io.Writer
}

// DefaultCopyOptions returns options populated from environment variables.
func DefaultCopyOptions() CopyOptions {
	return CopyOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   envOr("VAULT_MOUNT", "secret"),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		Output:  os.Stdout,
	}
}

// CopyResult holds the outcome of a copy operation.
type CopyResult struct {
	SourcePath string            `json:"source_path"`
	DestPath   string            `json:"dest_path"`
	Version    int               `json:"version"`
	Keys       int               `json:"keys"`
	DryRun     bool              `json:"dry_run"`
	CopiedAt   time.Time         `json:"copied_at"`
	Data       map[string]string `json:"data,omitempty"`
}

// Copy reads a secret version from srcPath and writes it to dstPath.
func Copy(ctx context.Context, srcPath, dstPath string, version int, opts CopyOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("copy: create client: %w", err)
	}

	data, err := client.GetSecretVersion(ctx, srcPath, version)
	if err != nil {
		return fmt.Errorf("copy: fetch source %s@%d: %w", srcPath, version, err)
	}

	result := CopyResult{
		SourcePath: srcPath,
		DestPath:   dstPath,
		Version:    version,
		Keys:       len(data),
		DryRun:     opts.DryRun,
		CopiedAt:   time.Now().UTC(),
	}

	if !opts.DryRun {
		if err := client.WriteSecret(ctx, dstPath, data); err != nil {
			return fmt.Errorf("copy: write dest %s: %w", dstPath, err)
		}
		result.Data = data
	}

	return json.NewEncoder(opts.Output).Encode(result)
}
