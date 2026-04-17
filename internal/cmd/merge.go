package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultdiff/internal/vault"
)

// MergeOptions configures the Merge command.
type MergeOptions struct {
	Address  string
	Token    string
	Mount    string
	SrcPath  string
	DstPath  string
	DryRun   bool
	Mask     bool
	Output   io.Writer
}

// DefaultMergeOptions returns options populated from environment variables.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	SrcPath  string            `json:"src_path"`
	DstPath  string            `json:"dst_path"`
	Merged   map[string]string `json:"merged,omitempty"`
	DryRun   bool              `json:"dry_run"`
	KeyCount int               `json:"key_count"`
}

// Merge copies keys from src into dst without overwriting existing dst keys.
func Merge(opts MergeOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("merge: new client: %w", err)
	}

	srcData, err := vault.GetSecretVersion(client, opts.SrcPath, 0)
	if err != nil {
		return fmt.Errorf("merge: fetch src: %w", err)
	}

	dstData, err := vault.GetSecretVersion(client, opts.DstPath, 0)
	if err != nil {
		return fmt.Errorf("merge: fetch dst: %w", err)
	}

	merged := make(map[string]string)
	for k, v := range dstData {
		merged[k] = v
	}
	for k, v := range srcData {
		if _, exists := merged[k]; !exists {
			merged[k] = v
		}
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, opts.DstPath, merged); err != nil {
			return fmt.Errorf("merge: write dst: %w", err)
		}
	}

	result := MergeResult{
		SrcPath:  opts.SrcPath,
		DstPath:  opts.DstPath,
		DryRun:   opts.DryRun,
		KeyCount: len(merged),
	}
	if !opts.Mask {
		result.Merged = merged
	}

	return json.NewEncoder(out).Encode(result)
}
