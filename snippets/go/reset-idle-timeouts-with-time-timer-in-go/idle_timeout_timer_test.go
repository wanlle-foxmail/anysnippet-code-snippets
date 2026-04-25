package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitForIdleTimeoutReturnsTimeoutWithoutActivity(t *testing.T) {
	start := time.Now()
	err := WaitForIdleTimeout(context.Background(), nil, 40*time.Millisecond)
	if !errors.Is(err, ErrIdleTimeout) {
		t.Fatalf("expected ErrIdleTimeout, got %v", err)
	}
	if elapsed := time.Since(start); elapsed < 35*time.Millisecond {
		t.Fatalf("expected timeout to wait for the idle window, got %v", elapsed)
	}
}

func TestWaitForIdleTimeoutResetsAfterActivity(t *testing.T) {
	activity := make(chan struct{}, 1)
	go func() {
		time.Sleep(25 * time.Millisecond)
		activity <- struct{}{}
	}()

	start := time.Now()
	err := WaitForIdleTimeout(context.Background(), activity, 50*time.Millisecond)
	if !errors.Is(err, ErrIdleTimeout) {
		t.Fatalf("expected ErrIdleTimeout, got %v", err)
	}
	if elapsed := time.Since(start); elapsed < 70*time.Millisecond {
		t.Fatalf("expected timeout to reset after activity, got %v", elapsed)
	}
}

func TestWaitForIdleTimeoutSupportsMultipleResets(t *testing.T) {
	activity := make(chan struct{}, 2)
	go func() {
		time.Sleep(25 * time.Millisecond)
		activity <- struct{}{}
		time.Sleep(25 * time.Millisecond)
		activity <- struct{}{}
	}()

	start := time.Now()
	err := WaitForIdleTimeout(context.Background(), activity, 50*time.Millisecond)
	if !errors.Is(err, ErrIdleTimeout) {
		t.Fatalf("expected ErrIdleTimeout, got %v", err)
	}
	if elapsed := time.Since(start); elapsed < 95*time.Millisecond {
		t.Fatalf("expected timeout to honor both resets, got %v", elapsed)
	}
}

func TestWaitForIdleTimeoutReturnsNilWhenActivityChannelCloses(t *testing.T) {
	activity := make(chan struct{})
	close(activity)

	if err := WaitForIdleTimeout(context.Background(), activity, 50*time.Millisecond); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWaitForIdleTimeoutHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	if err := WaitForIdleTimeout(ctx, nil, 100*time.Millisecond); err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestWaitForIdleTimeoutRejectsInvalidTimeout(t *testing.T) {
	if err := WaitForIdleTimeout(context.Background(), nil, 0); err == nil {
		t.Fatal("expected error for zero idle timeout")
	}
	if err := WaitForIdleTimeout(context.Background(), nil, -time.Second); err == nil {
		t.Fatal("expected error for negative idle timeout")
	}
}
