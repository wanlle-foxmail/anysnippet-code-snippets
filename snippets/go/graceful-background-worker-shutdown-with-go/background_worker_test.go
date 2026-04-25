package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBackgroundWorkerShutdownDrainsAcceptedJobs(t *testing.T) {
	processed := make([]string, 0, 2)
	var mu sync.Mutex

	worker, err := NewBackgroundWorker(2, func(job string) {
		mu.Lock()
		processed = append(processed, job)
		mu.Unlock()
	}, nil)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	if err := worker.Submit("job-1"); err != nil {
		t.Fatalf("submit job-1: %v", err)
	}
	if err := worker.Submit("job-2"); err != nil {
		t.Fatalf("submit job-2: %v", err)
	}
	if err := worker.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(processed) != 2 || processed[0] != "job-1" || processed[1] != "job-2" {
		t.Fatalf("expected both accepted jobs to drain in order, got %v", processed)
	}
}

func TestBackgroundWorkerRejectsSubmitAfterShutdown(t *testing.T) {
	worker, err := NewBackgroundWorker(1, func(job string) {}, nil)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}
	if err := worker.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	if err := worker.Submit("job"); !errors.Is(err, ErrWorkerClosed) {
		t.Fatalf("expected ErrWorkerClosed, got %v", err)
	}
}

func TestBackgroundWorkerReturnsQueueFullWhenBufferIsFull(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	worker, err := NewBackgroundWorker(1, func(job string) {
		started <- struct{}{}
		<-release
	}, nil)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	if err := worker.Submit("job-1"); err != nil {
		t.Fatalf("submit job-1: %v", err)
	}
	<-started
	if err := worker.Submit("job-2"); err != nil {
		t.Fatalf("submit job-2: %v", err)
	}
	if err := worker.Submit("job-3"); !errors.Is(err, ErrWorkerQueueFull) {
		t.Fatalf("expected ErrWorkerQueueFull, got %v", err)
	}

	close(release)
	if err := worker.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestBackgroundWorkerRunsCleanupOnce(t *testing.T) {
	var cleanupCalls atomic.Int32
	worker, err := NewBackgroundWorker(1, func(job string) {}, func() {
		cleanupCalls.Add(1)
	})
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}

	if err := worker.Shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := worker.Shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
	if cleanupCalls.Load() != 1 {
		t.Fatalf("expected cleanup to run once, got %d", cleanupCalls.Load())
	}
}

func TestBackgroundWorkerShutdownRespectsContextDeadline(t *testing.T) {
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	worker, err := NewBackgroundWorker(1, func(job string) {
		started <- struct{}{}
		<-release
	}, nil)
	if err != nil {
		t.Fatalf("new worker: %v", err)
	}
	if err := worker.Submit("job"); err != nil {
		t.Fatalf("submit: %v", err)
	}
	<-started

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	if err := worker.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}

	close(release)
	if err := worker.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown after release: %v", err)
	}
}

func TestNewBackgroundWorkerRejectsInvalidInput(t *testing.T) {
	if _, err := NewBackgroundWorker(0, func(job string) {}, nil); err == nil {
		t.Fatal("expected error for zero buffer")
	}
	if _, err := NewBackgroundWorker(1, nil, nil); err == nil {
		t.Fatal("expected error for nil process function")
	}
}
