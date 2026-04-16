package cmd

import (
	"context"
	"testing"

	"github.com/vaultdiff/internal/vault"
)

// resolveVaultStub satisfies the vault.Client contract for resolve tests.
func resolveVaultStub(t *testing.T, versions map[int]vault.VersionMeta) *vault.Client {
	t.Helper()
	client, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	_ = versions // stubs are exercised via the ListVersions mock in secret_test
	return client
}

func TestDefaultResolveOptions_Defaults(t *testing.T) {
	opts := DefaultResolveOptions()
	if opts.Mount != "secret" {
		t.Errorf("expected mount=secret, got %q", opts.Mount)
	}
	if !opts.SkipDestroyed {
		t.Error("expected SkipDestroyed=true")
	}
}

func TestResolveVersions_BothExplicit(t *testing.T) {
	client := resolveVaultStub(t, nil)
	opts := DefaultResolveOptions()

	got, err := ResolveVersions(context.Background(), client, "svc/config", 3, 5, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.From != 3 || got.To != 5 {
		t.Errorf("expected {3,5}, got %+v", got)
	}
}

func TestResolveVersions_FromEqualsTo_ReturnsError(t *testing.T) {
	client := resolveVaultStub(t, nil)
	opts := DefaultResolveOptions()

	_, err := ResolveVersions(context.Background(), client, "svc/config", 4, 4, opts)
	if err == nil {
		t.Fatal("expected error when from == to, got nil")
	}
}

func TestResolveVersions_InvalidAddress_ReturnsError(t *testing.T) {
	client, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:19999",
		Token:   "x",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	opts := DefaultResolveOptions()
	_, err = ResolveVersions(context.Background(), client, "svc/config", 0, 0, opts)
	if err == nil {
		t.Fatal("expected error for unreachable vault, got nil")
	}
}

func TestResolvedVersions_Fields(t *testing.T) {
	rv := ResolvedVersions{From: 1, To: 2}
	if rv.From != 1 {
		t.Errorf("From want 1, got %d", rv.From)
	}
	if rv.To != 2 {
		t.Errorf("To: want 2, got %d", rv.To)
	}
}
