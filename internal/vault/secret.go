package vault

import (
	"context"
	"fmt"
	"strconv"
)

// SecretVersion holds the data and metadata for a specific version of a secret.
type SecretVersion struct {
	Path    string
	Version int
	Data    map[string]string
}

// GetSecretVersion fetches a specific version of a KV v2 secret from Vault.
// If version is 0, the latest version is returned.
func (c *Client) GetSecretVersion(ctx context.Context, path string, version int) (*SecretVersion, error) {
	kvPath := fmt.Sprintf("%s/data/%s", c.mount, path)

	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{strconv.Itoa(version)}
	}

	secret, err := c.logical.ReadWithDataWithContext(ctx, kvPath, params)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q (version %d): %w", path, version, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret %q not found", path)
	}

	rawData, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("secret %q has no data field", path)
	}

	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("secret %q data is not a map", path)
	}

	result := &SecretVersion{
		Path:    path,
		Version: version,
		Data:    make(map[string]string, len(dataMap)),
	}

	for k, v := range dataMap {
		result.Data[k] = fmt.Sprintf("%v", v)
	}

	// Resolve the actual version from metadata if we requested latest.
	if version == 0 {
		if meta, ok := secret.Data["metadata"].(map[string]interface{}); ok {
			if v, ok := meta["version"].(json.Number); ok {
				n, _ := v.Int64()
				result.Version = int(n)
			}
		}
	}

	return result, nil
}
