package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultdiff/internal/vault"
)

// AnnotateOptions controls the behaviour of Annotate.
type AnnotateOptions struct {
	Address string
	Token   string
	Mount   string
	Mask    bool
	DryRun  bool
	Output  io.Writer
}

// DefaultAnnotateOptions reads defaults from the environment.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   mountFromEnv(),
		Mask:    os.Getenv("VAULTDIFF_MASK") != "false",
		DryRun:  false,
		Output:  os.Stdout,
	}
}

// AnnotateResult is the JSON-serialisable result of an Annotate call.
type AnnotateResult struct {
	Path        string            `json:"path"`
	Version     int               `json:"version"`
	Annotations map[string]string `json:"annotations"`
	DryRun      bool              `json:"dry_run"`
	Timestamp   time.Time         `json:"timestamp"`
}

// Annotate attaches free-form key/value annotations to a secret version by
// writing them as a reserved "_annotations" key in a new version.
func Annotate(path string, version int, annotations map[string]string, opts AnnotateOptions) error {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	client, err := vault.NewClient(opts.Address, opts.Token, opts.Mount)
	if err != nil {
		return fmt.Errorf("annotate: %w", err)
	}

	existing, err := vault.GetSecretVersion(client, path, version)
	if err != nil {
		return fmt.Errorf("annotate: fetch version %d: %w", version, err)
	}

	merged := make(map[string]string, len(existing)+len(annotations))
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range annotations {
		merged["_annotation."+k] = v
	}

	result := AnnotateResult{
		Path:        path,
		Version:     version,
		Annotations: annotations,
		DryRun:      opts.DryRun,
		Timestamp:   time.Now().UTC(),
	}

	if !opts.DryRun {
		if err := vault.WriteSecret(client, path, merged); err != nil {
			return fmt.Errorf("annotate: write: %w", err)
		}
	}

	return json.NewEncoder(opts.Output).Encode(result)
}
