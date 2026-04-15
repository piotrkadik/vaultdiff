package cmd

import "testing"

func TestExitCodeForDrift_NoDrift(t *testing.T) {
	got := ExitCodeForDrift(false)
	if got != ExitOK {
		t.Errorf("expected %d, got %d", ExitOK, got)
	}
}

func TestExitCodeForDrift_WithDrift(t *testing.T) {
	got := ExitCodeForDrift(true)
	if got != ExitDrift {
		t.Errorf("expected %d, got %d", ExitDrift, got)
	}
}

func TestExitConstants_AreDistinct(t *testing.T) {
	codes := []int{ExitOK, ExitDrift, ExitError, ExitCancelled}
	seen := make(map[int]bool)
	for _, c := range codes {
		if seen[c] {
			t.Errorf("duplicate exit code: %d", c)
		}
		seen[c] = true
	}
}

func TestExitOK_IsZero(t *testing.T) {
	if ExitOK != 0 {
		t.Errorf("ExitOK must be 0, got %d", ExitOK)
	}
}
