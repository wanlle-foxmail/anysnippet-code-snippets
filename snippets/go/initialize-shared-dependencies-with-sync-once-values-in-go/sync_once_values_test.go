package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewOnceLoaderCachesTheFirstSuccessfulValue(t *testing.T) {
	var calls atomic.Int32
	loader := NewOnceLoader(func() (string, error) {
		calls.Add(1)
		return "ready", nil
	})

	first, err := loader()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	second, err := loader()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if first != "ready" || second != "ready" {
		t.Fatalf("expected cached value ready, got %q and %q", first, second)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected one initializer call, got %d", calls.Load())
	}
}

func TestNewOnceLoaderCachesTheFirstError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	var calls atomic.Int32
	loader := NewOnceLoader(func() (string, error) {
		calls.Add(1)
		return "", wantErr
	})

	if _, err := loader(); !errors.Is(err, wantErr) {
		t.Fatalf("expected cached error %v, got %v", wantErr, err)
	}
	if _, err := loader(); !errors.Is(err, wantErr) {
		t.Fatalf("expected cached error %v, got %v", wantErr, err)
	}
	if calls.Load() != 1 {
		t.Fatalf("expected one initializer call, got %d", calls.Load())
	}
}

func TestNewOnceLoaderRunsTheInitializerOnceForConcurrentCallers(t *testing.T) {
	var calls atomic.Int32
	loader := NewOnceLoader(func() (string, error) {
		calls.Add(1)
		return "ready", nil
	})

	results := make(chan string, 12)
	var group sync.WaitGroup
	for range 12 {
		group.Add(1)
		go func() {
			defer group.Done()
			value, err := loader()
			if err != nil {
				t.Errorf("expected nil error, got %v", err)
				return
			}
			results <- value
		}()
	}
	group.Wait()
	close(results)

	for value := range results {
		if value != "ready" {
			t.Fatalf("expected ready, got %q", value)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("expected one initializer call, got %d", calls.Load())
	}
}

func TestNewOnceLoaderRejectsNilLoadFunctions(t *testing.T) {
	loader := NewOnceLoader[string](nil)
	if _, err := loader(); err == nil {
		t.Fatal("expected error for a nil load function")
	}
}

func TestNewOnceLoaderReplaysPanics(t *testing.T) {
	loader := NewOnceLoader(func() (string, error) {
		panic("boom")
	})

	assertPanicValue(t, func() {
		_, _ = loader()
	}, "boom")
	assertPanicValue(t, func() {
		_, _ = loader()
	}, "boom")
}

func assertPanicValue(t *testing.T, run func(), want interface{}) {
	t.Helper()
	defer func() {
		recovered := recover()
		if recovered != want {
			t.Fatalf("expected panic %v, got %v", want, recovered)
		}
	}()
	run()
}