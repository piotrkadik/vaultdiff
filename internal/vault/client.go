package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client with helper methods.
type Client struct {
	api     *vaultapi.Client
	Mount   string
}

// Config holds configuration for connecting to a Vault instance.
type Config struct {
	Address string
	Token   string
	Mount   string
}

// NewClient creates a new Vault client from the given config.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()
	vcfg.Address = cfg.Address

	api, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	api.SetToken(cfg.Token)

	mount := cfg.Mount
	if mount == "" {
		mount = "secret"
	}

	return &Client{api: api, Mount: mount}, nil
}

// SecretVersion represents a single version of a KV v2 secret.
type SecretVersion struct {
	Version  int
	Data     map[string]string
	Deleted  bool
}

// GetSecretVersion retrieves a specific version of a KV v2 secret.
// Pass version 0 to get the latest version.
func (c *Client) GetSecretVersion(ctx context.Context, path string, version int) (*SecretVersion, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	secret, err := c.api.KVv2(c.Mount).GetVersion(ctx, path, version)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q version %d: %w", path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q not found", path)
	}

	data := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		if s, ok := v.(string); ok {
			data[k] = s
		} else {
			data[k] = fmt.Sprintf("%v", v)
		}
	}

	return &SecretVersion{
		Version: secret.VersionMetadata.Version,
		Data:    data,
		Deleted: secret.VersionMetadata.DeletionTime.IsZero() == false,
	}, nil
}
