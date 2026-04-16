package vault

import (
	"context"
	"testing"
)

func TestWriteSecret_InvalidAddress_ReturnsError(t *testing.T) {
	client, err := NewClient("://bad", "root", "secret")
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
	_ = client
}

func TestToAnyMap_ConvertsStrings(t *testing.T) {
	in := map[string]string{"key": "value", "foo": "bar"}
	out := toAnyMap(in)
	if len(out) != len(in) {
		t.Fatalf("expected %d keys, got %d", len(in), len(out))
	}
	for k, v := range in {
		got, ok := out[k]
		if !ok {
			t.Errorf("missing key %q", k)
		}
		if got != v {
			t.Errorf("key %q: expected %q, got %v", k, v, got)
		}
	}
}

func TestWriteSecret_ReturnsNoErrorOnSuccess(t *testing.T) {
	// Requires a live Vault; skip in unit context.
	t.Skip("integration only")
	client, _ := NewClient("http://127.0.0.1:8200", "root", "secret")
	err := client.WriteSecret(context.Background(), "test/write", map[string]string{"k": "v"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
