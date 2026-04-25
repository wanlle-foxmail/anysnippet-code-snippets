package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestInFlightStringLoaderSharesConcurrentCallsForSameKey(t *testing.T) {
	var loadCalls atomic.Int32
	release := make(chan struct{})
	loader, err := NewInFlightStringLoader(func(key string) (string, error) {
		loadCalls.Add(1)
		<-release
		return "value-for-" + key, nil
	})
	if err != nil {
		t.Fatalf("new loader: %v", err)
	}

	results := make(chan struct {
		value  string
		shared bool
		err    error
	}, 2)

	var waitGroup sync.WaitGroup
	for range 2 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			value, shared, loadErr := loader.Load("user:1")
			results <- struct {
				value  string
				shared bool
				err    error
			}{value: value, shared: shared, err: loadErr}
		}()
	}

	time.Sleep(20 * time.Millisecond)
	close(release)
	waitGroup.Wait()
	close(results)

	if loadCalls.Load() != 1 {
		t.Fatalf("expected one underlying load call, got %d", loadCalls.Load())
	}
	for result := range results {
		if result.err != nil {
			t.Fatalf("expected nil error, got %v", result.err)
		}
		if result.value != "value-for-user:1" {
			t.Fatalf("expected shared value, got %q", result.value)
		}
		if !result.shared {
			t.Fatal("expected shared=true for concurrent duplicate calls")
		}
	}
}

func TestInFlightStringLoaderDoesNotShareDifferentKeys(t *testing.T) {
	var loadCalls atomic.Int32
	loader, err := NewInFlightStringLoader(func(key string) (string, error) {
		loadCalls.Add(1)
		return "value-for-" + key, nil
	})
	if err != nil {
		t.Fatalf("new loader: %v", err)
	}

	firstValue, firstShared, firstErr := loader.Load("user:1")
	secondValue, secondShared, secondErr := loader.Load("user:2")
	if firstErr != nil || secondErr != nil {
		t.Fatalf("expected nil errors, got %v and %v", firstErr, secondErr)
	}
	if firstValue != "value-for-user:1" || secondValue != "value-for-user:2" {
		t.Fatalf("unexpected values %q and %q", firstValue, secondValue)
	}
	if firstShared || secondShared {
		t.Fatal("expected different keys not to be marked as shared")
	}
	if loadCalls.Load() != 2 {
		t.Fatalf("expected two underlying load calls, got %d", loadCalls.Load())
	}
}

func TestInFlightStringLoaderRunsAgainAfterFirstCallCompletes(t *testing.T) {
	var loadCalls atomic.Int32
	loader, err := NewInFlightStringLoader(func(key string) (string, error) {
		loadCalls.Add(1)
		return "value-for-" + key, nil
	})
	if err != nil {
		t.Fatalf("new loader: %v", err)
	}

	_, firstShared, firstErr := loader.Load("user:1")
	_, secondShared, secondErr := loader.Load("user:1")
	if firstErr != nil || secondErr != nil {
		t.Fatalf("expected nil errors, got %v and %v", firstErr, secondErr)
	}
	if firstShared || secondShared {
		t.Fatal("expected sequential calls not to be marked as shared")
	}
	if loadCalls.Load() != 2 {
		t.Fatalf("expected two underlying load calls after sequential requests, got %d", loadCalls.Load())
	}
}

func TestInFlightStringLoaderSharesConcurrentErrors(t *testing.T) {
	var loadCalls atomic.Int32
	loadErr := errors.New("backend unavailable")
	release := make(chan struct{})
	loader, err := NewInFlightStringLoader(func(key string) (string, error) {
		loadCalls.Add(1)
		<-release
		return "", loadErr
	})
	if err != nil {
		t.Fatalf("new loader: %v", err)
	}

	results := make(chan struct {
		shared bool
		err    error
	}, 2)

	var waitGroup sync.WaitGroup
	for range 2 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			_, shared, gotErr := loader.Load("user:1")
			results <- struct {
				shared bool
				err    error
			}{shared: shared, err: gotErr}
		}()
	}

	time.Sleep(20 * time.Millisecond)
	close(release)
	waitGroup.Wait()
	close(results)

	if loadCalls.Load() != 1 {
		t.Fatalf("expected one underlying load call, got %d", loadCalls.Load())
	}
	for result := range results {
		if !errors.Is(result.err, loadErr) {
			t.Fatalf("expected shared backend error, got %v", result.err)
		}
		if !result.shared {
			t.Fatal("expected shared=true for concurrent duplicate errors")
		}
	}
}

func TestInFlightStringLoaderRejectsEmptyKey(t *testing.T) {
	loader, err := NewInFlightStringLoader(func(key string) (string, error) {
		return key, nil
	})
	if err != nil {
		t.Fatalf("new loader: %v", err)
	}

	if _, _, err := loader.Load(""); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewInFlightStringLoaderRejectsNilLoadFunction(t *testing.T) {
	if _, err := NewInFlightStringLoader(nil); err == nil {
		t.Fatal("expected error for nil load function")
	}
}
