package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewSemaphoreRejectsInvalidCapacity(t *testing.T) {
	if _, err := NewSemaphore(0); err == nil {
		t.Fatal("expected error for zero capacity")
	}
	if _, err := NewSemaphore(-1); err == nil {
		t.Fatal("expected error for negative capacity")
	}
}

func TestSemaphoreHonorsConcurrencyLimit(t *testing.T) {
	semaphore, err := NewSemaphore(2)
	if err != nil {
		t.Fatalf("new semaphore: %v", err)
	}

	var current atomic.Int32
	var maxSeen atomic.Int32
	entered := make(chan struct{}, 5)
	releaseGate := make(chan struct{})
	releaseErrors := make(chan error, 5)

	var waitGroup sync.WaitGroup
	for range 5 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if acquireErr := semaphore.Acquire(context.Background()); acquireErr != nil {
				releaseErrors <- acquireErr
				return
			}

			active := current.Add(1)
			for {
				previous := maxSeen.Load()
				if active <= previous || maxSeen.CompareAndSwap(previous, active) {
					break
				}
			}

			entered <- struct{}{}
			<-releaseGate
			current.Add(-1)
			releaseErrors <- semaphore.Release()
		}()
	}

	<-entered
	<-entered

	if current.Load() != 2 {
		t.Fatalf("expected 2 active workers, got %d", current.Load())
	}
	select {
	case <-entered:
		t.Fatal("expected remaining workers to block until a permit is released")
	default:
	}

	close(releaseGate)
	waitGroup.Wait()
	close(releaseErrors)

	if maxSeen.Load() != 2 {
		t.Fatalf("expected max concurrent workers to stay at 2, got %d", maxSeen.Load())
	}
	for releaseErr := range releaseErrors {
		if releaseErr != nil {
			t.Fatalf("unexpected acquire or release error: %v", releaseErr)
		}
	}
}

func TestSemaphoreReleaseMakesPermitAvailableAgain(t *testing.T) {
	semaphore, err := NewSemaphore(1)
	if err != nil {
		t.Fatalf("new semaphore: %v", err)
	}

	if err := semaphore.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if err := semaphore.Release(); err != nil {
		t.Fatalf("release: %v", err)
	}
	if !semaphore.TryAcquire() {
		t.Fatal("expected permit to become available again after release")
	}
}

func TestSemaphoreTryAcquireReturnsFalseWhenFull(t *testing.T) {
	semaphore, err := NewSemaphore(1)
	if err != nil {
		t.Fatalf("new semaphore: %v", err)
	}

	if err := semaphore.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if semaphore.TryAcquire() {
		t.Fatal("expected TryAcquire to return false when the semaphore is full")
	}
}

func TestSemaphoreAcquireRespectsContextCancellation(t *testing.T) {
	semaphore, err := NewSemaphore(1)
	if err != nil {
		t.Fatalf("new semaphore: %v", err)
	}
	if err := semaphore.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := semaphore.Acquire(ctx); err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestSemaphoreReleaseWithoutPermitReturnsError(t *testing.T) {
	semaphore, err := NewSemaphore(1)
	if err != nil {
		t.Fatalf("new semaphore: %v", err)
	}

	if err := semaphore.Release(); err == nil {
		t.Fatal("expected error when releasing without an acquired permit")
	}
}
