package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/vaultdiff/internal/cmd"
	"github.com/vaultdiff/internal/vault"
)

// compareVaultStub satisfies the interface used by Compare via vault.Client.
// We reuse the vaultStub helper pattern established in run_test.go.
func newCompareStub(dataA, dataB map[string]string, errA, errB error) *vaultStub {
	call := 0
	return &vaultStub{
		getFn: func(_ context.Context, path string, version int) (*vault.SecretVersion, error) {
			call++
			if call == 1 {
				if errA != nil {
					return nil, errA
				}
				return &vault.SecretVersion{Data: dataA, Version: 1}, nil
			}
			if errB != nil {
				return nil, errB
			}
			return &vault.SecretVersion{Data: dataB, Version: 2}, nil
		},
	}
}

func TestCompare_ShowsHeader(t *testing.T) {
	stub := newCompareStub(
		map[string]string{"KEY": "alpha"},
		map[string]string{"KEY": "beta"},
		nil, nil,
	)
	var buf bytes.Buffer
	opts := cmd.DefaultCompareOptions()
	opts.PathA = "secret/data/staging/db"
	opts.PathB = "secret/data/prod/db"
	opts.Output = &buf
	opts.Mask = false

	if err := cmd.Compare(context.Background(), stub.Client(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging/db") {
		t.Errorf("expected PathA in header, got: %s", out)
	}
	if !strings.Contains(out, "prod/db") {
		t.Errorf("expected PathB in header, got: %s", out)
	}
}

func TestCompare_ErrorOnFirstFetch(t *testing.T) {
	stub := newCompareStub(nil, nil, errors.New("vault unavailable"), nil)
	opts := cmd.DefaultCompareOptions()
	opts.PathA = "secret/data/a"
	opts.PathB = "secret/data/b"

	err := cmd.Compare(context.Background(), stub.Client(), opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCompare_NilOutputDefaultsToStdout(t *testing.T) {
	stub := newCompareStub(
		map[string]string{"X": "1"},
		map[string]string{"X": "1"},
		nil, nil,
	)
	opts := cmd.DefaultCompareOptions()
	opts.PathA = "secret/data/a"
	opts.PathB = "secret/data/b"
	opts.Output = nil // should default to os.Stdout without panic

	if err := cmd.Compare(context.Background(), stub.Client(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
