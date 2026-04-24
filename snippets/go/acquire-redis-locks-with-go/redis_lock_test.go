package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRedisLockStore struct {
	setNXFunc            func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	compareAndDeleteFunc func(ctx context.Context, key string, expectedValue string) (bool, error)
}

func (store fakeRedisLockStore) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return store.setNXFunc(ctx, key, value, ttl)
}

func (store fakeRedisLockStore) CompareAndDelete(ctx context.Context, key string, expectedValue string) (bool, error) {
	return store.compareAndDeleteFunc(ctx, key, expectedValue)
}

func TestAcquireRedisLockReturnsTrueWhenSetNXSucceeds(t *testing.T) {
	store := fakeRedisLockStore{
		setNXFunc: func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
			if key != "jobs:123" || value != "worker-1" || ttl != 15*time.Second {
				t.Fatalf("unexpected SetNX input: %q %q %v", key, value, ttl)
			}
			return true, nil
		},
		compareAndDeleteFunc: func(ctx context.Context, key string, expectedValue string) (bool, error) {
			return false, nil
		},
	}

	acquired, err := AcquireRedisLock(context.Background(), store, "jobs:123", "worker-1", 15*time.Second)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !acquired {
		t.Fatal("expected lock acquisition to succeed")
	}
}

func TestAcquireRedisLockReturnsFalseWhenAnotherWorkerAlreadyHoldsLock(t *testing.T) {
	store := fakeRedisLockStore{
		setNXFunc: func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
			return false, nil
		},
		compareAndDeleteFunc: func(ctx context.Context, key string, expectedValue string) (bool, error) {
			return false, nil
		},
	}

	acquired, err := AcquireRedisLock(context.Background(), store, "jobs:123", "worker-1", 15*time.Second)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if acquired {
		t.Fatal("expected lock acquisition to fail")
	}
}

func TestAcquireRedisLockWrapsStoreErrors(t *testing.T) {
	store := fakeRedisLockStore{
		setNXFunc: func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
			return false, errors.New("redis unavailable")
		},
		compareAndDeleteFunc: func(ctx context.Context, key string, expectedValue string) (bool, error) {
			return false, nil
		},
	}

	_, err := AcquireRedisLock(context.Background(), store, "jobs:123", "worker-1", 15*time.Second)
	if err == nil {
		t.Fatal("expected store error")
	}
}

func TestReleaseRedisLockReturnsTrueWhenOwnerMatches(t *testing.T) {
	store := fakeRedisLockStore{
		setNXFunc: func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
			return false, nil
		},
		compareAndDeleteFunc: func(ctx context.Context, key string, expectedValue string) (bool, error) {
			if key != "jobs:123" || expectedValue != "worker-1" {
				t.Fatalf("unexpected CompareAndDelete input: %q %q", key, expectedValue)
			}
			return true, nil
		},
	}

	released, err := ReleaseRedisLock(context.Background(), store, "jobs:123", "worker-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !released {
		t.Fatal("expected release to succeed")
	}
}

func TestReleaseRedisLockReturnsFalseWhenOwnerDoesNotMatch(t *testing.T) {
	store := fakeRedisLockStore{
		setNXFunc: func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
			return false, nil
		},
		compareAndDeleteFunc: func(ctx context.Context, key string, expectedValue string) (bool, error) {
			return false, nil
		},
	}

	released, err := ReleaseRedisLock(context.Background(), store, "jobs:123", "worker-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if released {
		t.Fatal("expected release to fail when owner token does not match")
	}
}

func TestRedisLockFunctionsRejectInvalidInput(t *testing.T) {
	store := fakeRedisLockStore{
		setNXFunc: func(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
			return true, nil
		},
		compareAndDeleteFunc: func(ctx context.Context, key string, expectedValue string) (bool, error) {
			return true, nil
		},
	}

	if _, err := AcquireRedisLock(context.Background(), nil, "jobs:123", "worker-1", 15*time.Second); err == nil {
		t.Fatal("expected nil store error")
	}
	if _, err := AcquireRedisLock(context.Background(), store, "", "worker-1", 15*time.Second); err == nil {
		t.Fatal("expected empty key error")
	}
	if _, err := AcquireRedisLock(context.Background(), store, "jobs:123", "", 15*time.Second); err == nil {
		t.Fatal("expected empty owner token error")
	}
	if _, err := AcquireRedisLock(context.Background(), store, "jobs:123", "worker-1", 0); err == nil {
		t.Fatal("expected non-positive ttl error")
	}
	if _, err := ReleaseRedisLock(context.Background(), nil, "jobs:123", "worker-1"); err == nil {
		t.Fatal("expected nil store error on release")
	}
	if _, err := ReleaseRedisLock(context.Background(), store, "jobs:123", "   "); err == nil {
		t.Fatal("expected empty owner token error on release")
	}
}
