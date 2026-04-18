package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultdiff/internal/vault"
)

// StatOptions configures the Stat command.
type StatOptions struct {
	Address string
	Mount   string
	Token   string
	Mask    bool
	Format  string
	Output  io.Writer
}

// DefaultStatOptions returns StatOptions populated from environment variables.
func DefaultStatOptions() StatOptions {
	mask := true
	if os.Getenv("VAULTDIFF_MASK") == "false" {
		mask = false
	}
	mount := os.Getenv("VAULT_MOUNT")
	if mount == "" {
		mount = "secret"
	}
	return StatOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mount,
		Mask:    mask,
		Format:  os.Getenv("VAULTDIFF_FORMAT"),
	}
}

// StatResult holds metadata about a secret path.
type StatResult struct {
	Path           string `json:"path"`
	LatestVersion  int    `json:"latest_version"`
	TotalVersions  int    `json:"total_versions"`
	DestroyedCount int    `json:"destroyed_count"`
	KeyCount       int    `json:"key_count"`
}

// Stat prints metadata about a secret path and its versions.
func Stat(path string, opts StatOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	c, err := vault.NewClient(opts.Address, opts.Mount, opts.Token)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}

	metas, err := vault.ListVersions(c, path)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}

	latest := vault.LatestVersion(metas)
	destroyed := 0
	for _, m := range metas {
		if m.Destroyed {
			destroyed++
		}
	ecret, err := vault.GetSecretVersion(c, path, latest)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}

	res := StatResult{
		Path:           path,
		LatestVersion:  latest,
		TotalVersions:  len(metas),
		DestroyedCount: destroyed,
		KeyCount:       len(secret.Data),
	}

	if opts.Format == "json" {
		return json.NewEncoder(out).Encode(res)
	}

	fmt.Fprintf(out, "Path:            %s\n", res.Path)
	fmt.Fprintf(out, "Latest Version:  %d\n", res.LatestVersion)
	fmt.Fprintf(out, "Total Versions:  %d\n", res.TotalVersions)
	fmt.Fprintf(out, "Destroyed:       %d\n", res.DestroyedCount)
	fmt.Fprintf(out, "Keys:            %d\n", res.KeyCount)
	return nil
}
