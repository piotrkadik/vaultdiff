package vault

import (
	"context"
	"errors"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func makeVersionsSecret(versions map[string]interface{}) *vaultapi.Secret {
	return &vaultapi.Secret{
		Data: map[string]interface{}{
			"versions": versions,
		},
	}
}

func TestListVersions_ReturnsMetadata(t *testing.T) {
	versions := map[string]interface{}{
		"1": map[string]interface{}{
			"created_time":  "2024-01-01T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
		"2": map[string]interface{}{
			"created_time":  "2024-02-01T00:00:00Z",
			"deletion_time": "",
			"destroyed":     false,
		},
	}
	mock := &mockLogical{secret: makeVersionsSecret(versions)}
	c := &Client{mount: "secret", logical: mock}

	metas, err := c.ListVersions(context.Background(), "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(metas))
	}
}

func TestListVersions_NotFound(t *testing.T) {
	mock := &mockLogical{secret: nil}
	c := &Client{mount: "secret", logical: mock}

	_, err := c.ListVersions(context.Background(), "missing/path")
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestListVersions_ErrorPropagated(t *testing.T) {
	mock := &mockLogical{err: errors.New("vault unavailable")}
	c := &Client{mount: "secret", logical: mock}

	_, err := c.ListVersions(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestLatestVersion_SkipsDestroyed(t *testing.T) {
	versions := map[string]interface{}{
		"1": map[string]interface{}{
			"created_time": "2024-01-01T00:00:00Z",
			"destroyed":    false,
		},
		"2": map[string]interface{}{
			"created_time": "2024-02-01T00:00:00Z",
			"destroyed":    true,
		},
	}
	mock := &mockLogical{secret: makeVersionsSecret(versions)}
	c := &Client{mount: "secret", logical: mock}

	v, err := c.LatestVersion(context.Background(), "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 1 {
		t.Errorf("expected latest non-destroyed version 1, got %d", v)
	}
}
