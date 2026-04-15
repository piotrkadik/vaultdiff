package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// TestWatch_CancelledContextExits verifies that Watch returns promptly when
// the context is cancelled before the first tick.
func TestWatch_CancelledContextExits(t *testing.T) {
	client := vaultStub(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	opts := WatchOptions{
		Interval:  10 * time.Millisecond,
		MaxRounds: 5,
		Output:    &bytes.Buffer{},
	}

	err := Watch(ctx, client, "secret/data/app", opts)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}

// TestWatch_MaxRoundsStops verifies that Watch stops after MaxRounds ticks.
func TestWatch_MaxRoundsStops(t *testing.T) {
	client := vaultStub(t)

	var buf bytes.Buffer
	opts := WatchOptions{
		Interval:  5 * time.Millisecond,
		MaxRounds: 2,
		Output:    &buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := Watch(ctx, client, "secret/data/app", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestWatch_NilOutputDefaultsToStdout ensures a nil Output does not panic.
func TestWatch_NilOutputDefaultsToStdout(t *testing.T) {
	client := vaultStub(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	opts := WatchOptions{
		Interval:  5 * time.Millisecond,
		MaxRounds: 1,
		Output:    nil, // should default internally
	}

	// Should not panic
	_ = Watch(ctx, client, "secret/data/app", opts)
}

// TestDefaultWatchOptions_Interval checks that the default interval is set.
func TestDefaultWatchOptions_Interval(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", opts.Interval)
	}
	if opts.MaxRounds != 0 {
		t.Errorf("expected MaxRounds=0 (unlimited), got %d", opts.MaxRounds)
	}
	if opts.Output == nil {
		t.Error("expected non-nil default output")
	}
}

// TestWatch_PrintsTrackingMessage verifies the initial tracking banner.
func TestWatch_PrintsTrackingMessage(t *testing.T) {
	client := vaultStub(t)

	var buf bytes.Buffer
	opts := WatchOptions{
		Interval:  5 * time.Millisecond,
		MaxRounds: 1,
		Output:    &buf,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = Watch(ctx, client, "secret/data/app", opts)

	if !strings.Contains(buf.String(), "tracking") {
		t.Errorf("expected tracking message in output, got: %s", buf.String())
	}
}
