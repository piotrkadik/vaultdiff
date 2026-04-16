package vault

import (
	"context"
	"fmt"
	"path"
)

// WriteSecret writes data to the KV v2 secret at secretPath under the client's mount.
func (c *Client) WriteSecret(ctx context.Context, secretPath string, data map[string]string) error {
	kv := c.logical
	fullPath := path.Join(c.mount, "data", secretPath)

	payload := map[string]interface{}{
		"data": toAnyMap(data),
	}

	_, err := kv.WriteWithContext(ctx, fullPath, payload)
	if err != nil {
		return fmt.Errorf("write secret %s: %w", secretPath, err)
	}
	return nil
}

func toAnyMap(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
