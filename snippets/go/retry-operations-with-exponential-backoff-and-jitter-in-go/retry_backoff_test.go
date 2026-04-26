package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryOperationSucceedsOnTheFirstAttempt(t *testing.T) {
	restoreRandom, restoreWait := stubRetryHelpers(t, 0.5, nil)
	defer restoreRandom()
	defer restoreWait()

	attempts := 0
	err := RetryOperation(context.Background(), func(context.Context) error {
		attempts++
		return nil
	}, func(error) bool {
		return true
	}, RetryConfig{
		MaxAttempts:    3,
		BaseDelay:      10 * time.Millisecond,
		MaxDelay:       40 * time.Millisecond,
		JitterFraction: 0.2,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected one attempt, got %d", attempts)
	}
}

func TestRetryOperationRetriesBeforeSuccess(t *testing.T) {
	var delays []time.Duration
	restoreRandom, restoreWait := stubRetryHelpers(t, 0.5, func(_ context.Context, delay time.Duration) error {
		delays = append(delays, delay)
		return nil
	})
	defer restoreRandom()
	defer restoreWait()

	attempts := 0
	err := RetryOperation(context.Background(), func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errTemporary
		}
		return nil
	}, func(err error) bool {
		return errors.Is(err, errTemporary)
	}, RetryConfig{
		MaxAttempts:    4,
		BaseDelay:      10 * time.Millisecond,
		MaxDelay:       40 * time.Millisecond,
		JitterFraction: 0.5,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected three attempts, got %d", attempts)
	}
	if len(delays) != 2 || delays[0] != 10*time.Millisecond || delays[1] != 20*time.Millisecond {
		t.Fatalf("unexpected delays %v", delays)
	}
}

func TestRetryOperationStopsOnNonRetryableErrors(t *testing.T) {
	restoreRandom, restoreWait := stubRetryHelpers(t, 0.5, nil)
	defer restoreRandom()
	defer restoreWait()

	wantErr := errors.New("bad request")
	err := RetryOperation(context.Background(), func(context.Context) error {
		return wantErr
	}, func(err error) bool {
		return errors.Is(err, errTemporary)
	}, RetryConfig{
		MaxAttempts:    4,
		BaseDelay:      10 * time.Millisecond,
		MaxDelay:       40 * time.Millisecond,
		JitterFraction: 0.2,
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}

func TestRetryOperationReturnsTheLastRetryableError(t *testing.T) {
	var delays []time.Duration
	restoreRandom, restoreWait := stubRetryHelpers(t, 0.5, func(_ context.Context, delay time.Duration) error {
		delays = append(delays, delay)
		return nil
	})
	defer restoreRandom()
	defer restoreWait()

	err := RetryOperation(context.Background(), func(context.Context) error {
		return errTemporary
	}, func(err error) bool {
		return errors.Is(err, errTemporary)
	}, RetryConfig{
		MaxAttempts:    3,
		BaseDelay:      10 * time.Millisecond,
		MaxDelay:       40 * time.Millisecond,
		JitterFraction: 0.2,
	})
	if !errors.Is(err, errTemporary) {
		t.Fatalf("expected error %v, got %v", errTemporary, err)
	}
	if len(delays) != 2 {
		t.Fatalf("expected two backoff delays, got %v", delays)
	}
}

func TestRetryOperationStopsWhenTheContextIsCanceledDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	restoreRandom, restoreWait := stubRetryHelpers(t, 0.5, func(ctx context.Context, delay time.Duration) error {
		cancel()
		<-ctx.Done()
		return ctx.Err()
	})
	defer restoreRandom()
	defer restoreWait()

	err := RetryOperation(ctx, func(context.Context) error {
		return errTemporary
	}, func(err error) bool {
		return errors.Is(err, errTemporary)
	}, RetryConfig{
		MaxAttempts:    3,
		BaseDelay:      10 * time.Millisecond,
		MaxDelay:       40 * time.Millisecond,
		JitterFraction: 0.2,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestRetryOperationRejectsInvalidInput(t *testing.T) {
	if err := RetryOperation(context.Background(), nil, func(error) bool { return true }, RetryConfig{MaxAttempts: 1, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}); err == nil {
		t.Fatal("expected error for a nil operation")
	}
	if err := RetryOperation(context.Background(), func(context.Context) error { return nil }, nil, RetryConfig{MaxAttempts: 1, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}); err == nil {
		t.Fatal("expected error for a nil retry predicate")
	}
	if err := RetryOperation(context.Background(), func(context.Context) error { return nil }, func(error) bool { return true }, RetryConfig{MaxAttempts: 0, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}); err == nil {
		t.Fatal("expected error for max attempts")
	}
	if err := RetryOperation(context.Background(), func(context.Context) error { return nil }, func(error) bool { return true }, RetryConfig{MaxAttempts: 1, BaseDelay: 0, MaxDelay: time.Millisecond}); err == nil {
		t.Fatal("expected error for base delay")
	}
	if err := RetryOperation(context.Background(), func(context.Context) error { return nil }, func(error) bool { return true }, RetryConfig{MaxAttempts: 1, BaseDelay: 2 * time.Millisecond, MaxDelay: time.Millisecond}); err == nil {
		t.Fatal("expected error for max delay")
	}
	if err := RetryOperation(context.Background(), func(context.Context) error { return nil }, func(error) bool { return true }, RetryConfig{MaxAttempts: 1, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, JitterFraction: 2}); err == nil {
		t.Fatal("expected error for jitter fraction")
	}
}

func TestBackoffDelayClampsVeryLargeAttemptsBeforeOverflow(t *testing.T) {
	restoreRandom, restoreWait := stubRetryHelpers(t, 0.5, nil)
	defer restoreRandom()
	defer restoreWait()

	delay := backoffDelay(100, time.Second, 10*time.Second, 0)
	if delay != 10*time.Second {
		t.Fatalf("expected delay 10s, got %v", delay)
	}
}

func stubRetryHelpers(t *testing.T, randomValue float64, wait func(context.Context, time.Duration) error) (func(), func()) {
	t.Helper()
	originalRandom := randomFloat64
	originalWait := waitForDelay
	randomFloat64 = func() float64 {
		return randomValue
	}
	if wait != nil {
		waitForDelay = wait
	} else {
		waitForDelay = func(context.Context, time.Duration) error {
			return nil
		}
	}
	return func() {
			randomFloat64 = originalRandom
		}, func() {
			waitForDelay = originalWait
		}
}
