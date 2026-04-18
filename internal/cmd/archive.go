package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/exampleorg/vaultdiff/internal/vault"
)

// DefaultArchiveOptions returns ArchiveOptions populated from environment variables.
func DefaultArchiveOptions() ArchiveOptions {
	return ArchiveOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
	}
}

// ArchiveOptions configures the Archive command.
type ArchiveOptions struct {
	Address string
	Token   string
	Mount   string
	Path    string
	Version int
	Mask    bool
	Output  io.Writer
}

// ArchiveRecord holds a point-in-time snapshot of a secret version.
type ArchiveRecord struct {
	Path      string            `json:"path"`
	Version   int               `json:"version"`
	Data      map[string]string `json:"data,omitempty"`
	ArchivedAt time.Time        `json:"archived_at"`
}

// Archive fetches a specific secret version and writes it as a JSON archive record.
func Archive(opts ArchiveOptions) error {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("archive: %w", err)
	}

	secret, err := vault.GetSecretVersion(client, opts.Path, opts.Version)
	if err != nil {
		return fmt.Errorf("archive: %w", err)
	}

	data := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		if opts.Mask {
			data[k] = "***"
		} else {
			data[k] = v
		}
	}

	rec := ArchiveRecord{
		Path:       opts.Path,
		Version:    secret.Version,
		Data:       data,
		ArchivedAt: time.Now().UTC(),
	}

	return json.NewEncoder(out).Encode(rec)
}
