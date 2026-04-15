package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

// fakeLogical is a minimal stub that satisfies the logical interface used by Client.
type fakeLogical struct {
	response *vaultapi.Secret
	err      error
}

func (f *fakeLogical) ReadWithDataWithContext(_ context.Context, _ string, _ map[string][]string) (*vaultapi.Secret, error) {
	return f.response, f.err
}

func makeSecret(data map[string]interface{}, version int) *vaultapi.Secret {
	return &vaultapi.Secret{
		Data: map[string]interface{}{
			"data": data,
			"metadata": map[string]interface{}{
				"version": json.Number(strconv.Itoa(version)),
			},
		},
	}
}

func TestGetSecretVersion_ReturnsData(t *testing.T) {
	c := &Client{
		mount:   "secret",
		logical: &fakeLogical{response: makeSecret(map[string]interface{}{"API_KEY": "abc123"}, 3)},
	}

	sv, err := c.GetSecretVersion(context.Background(), "myapp/config", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sv.Data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", sv.Data["API_KEY"])
	}
	if sv.Version != 3 {
		t.Errorf("expected version 3, got %d", sv.Version)
	}
}

func TestGetSecretVersion_NotFound(t *testing.T) {
	c := &Client{
		mount:   "secret",
		logical: &fakeLogical{response: nil},
	}

	_, err := c.GetSecretVersion(context.Background(), "missing/path", 1)
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestGetSecretVersion_ErrorPropagated(t *testing.T) {
	c := &Client{
		mount:   "secret",
		logical: &fakeLogical{err: fmt.Errorf("permission denied")},
	}

	_, err := c.GetSecretVersion(context.Background(), "myapp/config", 1)
	if err == nil {
		t.Fatal("expected error to be propagated")
	}
}
