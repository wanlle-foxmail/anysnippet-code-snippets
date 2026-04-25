package main

import (
	"context"
	"testing"
	"time"
)

func TestStartHeartbeatsRejectsInvalidInterval(t *testing.T) {
	if _, err := StartHeartbeats(context.Background(), 0); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestStartHeartbeatsEmitsHeartbeat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeats, err := StartHeartbeats(ctx, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("start heartbeats: %v", err)
	}

	select {
	case <-heartbeats:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for a heartbeat")
	}
}

func TestStartHeartbeatsDoesNotEmitImmediately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeats, err := StartHeartbeats(ctx, 40*time.Millisecond)
	if err != nil {
		t.Fatalf("start heartbeats: %v", err)
	}

	select {
	case <-heartbeats:
		t.Fatal("expected no immediate heartbeat before the first tick")
	case <-time.After(15 * time.Millisecond):
	}
}

func TestStartHeartbeatsClosesOutputOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	heartbeats, err := StartHeartbeats(ctx, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("start heartbeats: %v", err)
	}

	cancel()

	select {
	case _, ok := <-heartbeats:
		if ok {
			t.Fatal("expected output to close after cancellation")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for output close after cancellation")
	}
}

func TestStartHeartbeatsDropsWhenReceiverFallsBehind(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeats, err := StartHeartbeats(ctx, 15*time.Millisecond)
	if err != nil {
		t.Fatalf("start heartbeats: %v", err)
	}

	time.Sleep(70 * time.Millisecond)
	if len(heartbeats) != 1 {
		t.Fatalf("expected exactly one buffered heartbeat, got %d", len(heartbeats))
	}
}

func TestStartHeartbeatsEmitsMultipleTicksWhenDrained(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeats, err := StartHeartbeats(ctx, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("start heartbeats: %v", err)
	}

	select {
	case <-heartbeats:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for first heartbeat")
	}

	select {
	case <-heartbeats:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for second heartbeat")
	}
}
