package main

import (
	"testing"
	"time"
)

type fakeClock struct {
	now time.Time
}

func (clock fakeClock) Now() time.Time {
	return clock.now
}

func TestHasExpiredReturnsFalseBeforeExpiration(t *testing.T) {
	clock := fakeClock{now: time.Date(2026, time.April, 25, 10, 0, 0, 0, time.UTC)}
	expiresAt := clock.now.Add(10 * time.Minute)

	expired, err := HasExpired(clock, expiresAt)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if expired {
		t.Fatal("expected token to be active")
	}
}

func TestHasExpiredReturnsTrueAfterExpiration(t *testing.T) {
	clock := fakeClock{now: time.Date(2026, time.April, 25, 10, 0, 0, 0, time.UTC)}
	expiresAt := clock.now.Add(-time.Minute)

	expired, err := HasExpired(clock, expiresAt)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !expired {
		t.Fatal("expected token to be expired")
	}
}

func TestHasExpiredTreatsExactDeadlineAsExpired(t *testing.T) {
	now := time.Date(2026, time.April, 25, 10, 0, 0, 0, time.UTC)

	expired, err := HasExpired(fakeClock{now: now}, now)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !expired {
		t.Fatal("expected exact deadline to count as expired")
	}
}

func TestHasExpiredReturnsErrorForNilClock(t *testing.T) {
	if _, err := HasExpired(nil, time.Now().Add(time.Minute)); err == nil {
		t.Fatal("expected error for nil clock")
	}
}

func TestHasExpiredReturnsErrorForZeroExpiration(t *testing.T) {
	if _, err := HasExpired(fakeClock{now: time.Now()}, time.Time{}); err == nil {
		t.Fatal("expected error for zero expiration time")
	}
}

func TestRealClockNowFallsWithinCallWindow(t *testing.T) {
	before := time.Now()
	now := RealClock{}.Now()
	after := time.Now()

	if now.Before(before) {
		t.Fatal("expected real clock time to be on or after the start of the call window")
	}
	if now.After(after) {
		t.Fatal("expected real clock time to be on or before the end of the call window")
	}
}
