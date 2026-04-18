package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// DefaultSyncOptions returns SyncOptions populated from environment variables.
func DefaultSyncOptions() SyncOptions {
	return SyncOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		DryRun:  false,
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// SyncOptions configures a Sync operation.
type SyncOptions struct {
	Address  string
	Token    string
	Mount    string
	SrcPath  string
	DstPath  string
	DryRun   bool
	Mask     bool
	Output   io.Writer
}

// SyncResult holds the outcome of a Sync call.
type SyncResult struct {
	SrcPath string            `json:"src_path"`
	DstPath string            `json:"dst_path"`
	DryRun  bool              `json:"dry_run"`
	Keys    int               `json:"keys"`
	Data    map[string]string `json:"data,omitempty"`
}

// Sync copies all keys from SrcPath to DstPath, optionally skipping the write.
func Sync(opts SyncOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("sync: new client: %w", err)
	}

	data, err := vault.GetSecretVersion(client, opts.SrcPath, 0)
	if err != nil {
		return fmt.Errorf("sync: fetch src: %w", err)
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, opts.DstPath, data); err != nil {
			return fmt.Errorf("sync: write dst: %w", err)
		}
	}

	result := SyncResult{
		SrcPath: opts.SrcPath,
		DstPath: opts.DstPath,
		DryRun:  opts.DryRun,
		Keys:    len(data),
	}
	if !opts.Mask {
		result.Data = data
	}

	return json.NewEncoder(out).Encode(result)
}
