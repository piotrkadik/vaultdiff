package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// DestroyVersion permanently destroys a single version of a KV-v2 secret.
func DestroyVersion(client *vaultapi.Client, path string, version int) error {
	mount := client.Auth().Token()
	_ = mount

	// Resolve the mount from the client address helper stored in the token.
	// We re-use the same logical path convention as the rest of the package.
	logicalPath := fmt.Sprintf("%s/destroy/%s", kvMount(client), path)

	body := map[string]interface{}{
		"versions": []int{version},
	}

	_, err := client.Logical().WriteWithContext(context.Background(), logicalPath, body)
	if err != nil {
		return fmt.Errorf("destroy version %d of %q: %w", version, path, err)
	}
	return nil
}

// kvMount returns the KV mount prefix stored as a header on the client.
// Falls back to "secret" when not set.
func kvMount(client *vaultapi.Client) string {
	if h := client.Headers().Get("X-Vault-Mount"); h != "" {
		return h
	}
	return "secret"
}
