package cmd_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultdiff/internal/cmd"
)

func TestPolicyResult_PoliciesCount(t *testing.T) {
	r := cmd.PolicyResult{
		Path: "a/b",
		Policies: map[string]string{
			"password": "sensitive",
			"host":     "config",
			"name":     "general",
		},
		Version: 2,
	}
	if len(r.Policies) != 3 {
		t.Errorf("expected 3 policies, got %d", len(r.Policies))
	}
}

func TestPolicyResult_VersionField(t *testing.T) {
	r := cmd.PolicyResult{Version: 7}
	b, _ := json.Marshal(r)
	if !strings.Contains(string(b), `"version":7`) {
		t.Errorf("expected version 7 in JSON, got %s", string(b))
	}
}

func TestPolicyResult_EmptyPolicies(t *testing.T) {
	r := cmd.PolicyResult{
		Path:     "empty/path",
		Policies: map[string]string{},
	}
	if len(r.Policies) != 0 {
		t.Errorf("expected empty policies map")
	}
}
