package vault

import (
	"context"
	"fmt"
	"sort"
)

// VersionMeta holds metadata for a single secret version.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns metadata for all versions of a KV v2 secret at path.
// Versions are returned sorted ascending by version number.
func (c *Client) ListVersions(ctx context.Context, path string) ([]VersionMeta, error) {
	metaPath := fmt.Sprintf("%s/metadata/%s", c.mount, path)

	secret, err := c.logical.ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("listing versions for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found for path %q", path)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("metadata response missing 'versions' field for path %q", path)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for versions field")
	}

	var metas []VersionMeta
	for _, v := range versionsMap {
		entry, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		meta := VersionMeta{}
		if ct, ok := entry["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := entry["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if d, ok := entry["destroyed"].(bool); ok {
			meta.Destroyed = d
		}
		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})

	return metas, nil
}

// LatestVersion returns the highest available (non-destroyed) version number
// for the given secret path, or 0 if none found.
func (c *Client) LatestVersion(ctx context.Context, path string) (int, error) {
	versions, err := c.ListVersions(ctx, path)
	if err != nil {
		return 0, err
	}
	for i := len(versions) - 1; i >= 0; i-- {
		if !versions[i].Destroyed {
			return versions[i].Version, nil
		}
	}
	return 0, fmt.Errorf("no available versions for path %q", path)
}
