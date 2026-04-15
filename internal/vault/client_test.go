package vault_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourorg/vaultdiff/internal/vault"
)

func TestNewClient_DefaultMount(t *testing.T) {
	client, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
	})
	require.NoError(t, err)
	assert.Equal(t, "secret", client.Mount)
}

func TestNewClient_CustomMount(t *testing.T) {
	client, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
		Token:   "root",
		Mount:   "kv",
	})
	require.NoError(t, err)
	assert.Equal(t, "kv", client.Mount)
}

func TestNewClient_InvalidAddress(t *testing.T) {
	// The Vault SDK accepts any string as address during construction;
	// errors surface on actual requests, so construction should succeed.
	client, err := vault.NewClient(vault.Config{
		Address: "not-a-real-address",
		Token:   "token",
	})
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestSecretVersion_Fields(t *testing.T) {
	sv := vault.SecretVersion{
		Version: 3,
		Data:    map[string]string{"key": "value"},
		Deleted: false,
	}

	assert.Equal(t, 3, sv.Version)
	assert.Equal(t, "value", sv.Data["key"])
	assert.False(t, sv.Deleted)
}
