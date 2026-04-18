package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your/vaultdiff/internal/vault"
)

// MirrorOptions configures the Mirror command.
type MirrorOptions struct {
	Address  string
	Token    string
	Mount    string
	SrcPath  string
	DstPath  string
	DryRun   bool
	Mask     bool
	Output   io.Writer
}

// DefaultMirrorOptions returns options populated from environment variables.
func DefaultMirrorOptions() MirrorOptions {
	mask := true
	if os.Getenv("VAULT_MASK") == "false" {
		mask = false
	}
	return MirrorOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   os.Getenv("VAULT_MOUNT"),
		Mask:    mask,
	}
}

// MirrorResult holds the outcome of a mirror operation.
type MirrorResult struct {
	SrcPath string            `json:"src_path"`
	DstPath string            `json:"dst_path"`
	DryRun  bool              `json:"dry_run"`
	Keys    int               `json:"keys"`
	Data    map[string]string `json:"data,omitempty"`
}

// Mirror copies all key-value pairs from one secret path to another.
func Mirror(opts MirrorOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("mirror: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.SrcPath, 0)
	if err != nil {
		return fmt.Errorf("mirror: fetch src: %w", err)
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, opts.DstPath, secret.Data); err != nil {
			return fmt.Errorf("mirror: write dst: %w", err)
		}
	}

	result := MirrorResult{
		SrcPath: opts.SrcPath,
		DstPath: opts.DstPath,
		DryRun:  opts.DryRun,
		Keys:    len(secret.Data),
	}
	if !opts.Mask {
		result.Data = secret.Data
	}

	return json.NewEncoder(out).Encode(result)
}
