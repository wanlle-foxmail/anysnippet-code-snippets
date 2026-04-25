package main

import (
	"context"
	"errors"
	"log"
	"sync"
)

var ErrWorkerClosed = errors.New("worker is shutting down")
var ErrWorkerQueueFull = errors.New("worker queue is full")

type BackgroundWorker struct {
	mu      sync.RWMutex
	jobs    chan string
	done    chan struct{}
	process func(string)
	cleanup func()
	closing bool
}

func NewBackgroundWorker(buffer int, process func(string), cleanup func()) (*BackgroundWorker, error) {
	if buffer <= 0 {
		return nil, errors.New("buffer must be greater than 0")
	}
	if process == nil {
		return nil, errors.New("process function is required")
	}
	if cleanup == nil {
		cleanup = func() {}
	}

	worker := &BackgroundWorker{
		jobs:    make(chan string, buffer),
		done:    make(chan struct{}),
		process: process,
		cleanup: cleanup,
	}

	go worker.run()
	return worker, nil
}

func (worker *BackgroundWorker) run() {
	defer close(worker.done)
	defer worker.cleanup()

	for job := range worker.jobs {
		worker.process(job)
	}
}

func (worker *BackgroundWorker) Submit(job string) error {
	if worker == nil {
		return errors.New("worker is required")
	}

	worker.mu.RLock()
	defer worker.mu.RUnlock()

	if worker.closing {
		return ErrWorkerClosed
	}

	select {
	case worker.jobs <- job:
		return nil
	default:
		return ErrWorkerQueueFull
	}
}

// Flow:
//
//	mark worker as closing
//	   |
//	   +-> future submits return ErrWorkerClosed
//	   +-> accepted jobs keep draining until the queue empties
//	   `-> return nil on worker exit or ctx.Err() on shutdown timeout
func (worker *BackgroundWorker) Shutdown(ctx context.Context) error {
	if worker == nil {
		return errors.New("worker is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	worker.mu.Lock()
	if !worker.closing {
		worker.closing = true
		close(worker.jobs)
	}
	worker.mu.Unlock()

	select {
	case <-worker.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	worker, err := NewBackgroundWorker(2, func(job string) {
		log.Printf("processed %s", job)
	}, func() {
		log.Println("cleanup complete")
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := worker.Submit("job-1"); err != nil {
		log.Fatal(err)
	}
	if err := worker.Submit("job-2"); err != nil {
		log.Fatal(err)
	}
	if err := worker.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
