package vault

import (
	"context"
	"fmt"
)

// SecretWriter is the interface required by callers that need to persist
// secret data back to Vault.
type SecretWriter interface {
	WriteSecret(ctx context.Context, path string, data map[string]any) error
}

// WriteSecret writes data to the KVv2 secret at path using the client's mount.
// The Vault API creates a new version; existing versions are preserved.
func (c *Client) WriteSecret(ctx context.Context, path string, data map[string]any) error {
	_, err := c.kv.Put(ctx, path, data)
	if err != nil {
		return fmt.Errorf("vault: write secret %q: %w", path, err)
	}
	return nil
}
