package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Semaphore struct {
	tokens chan struct{}
}

func NewSemaphore(maxConcurrent int) (*Semaphore, error) {
	if maxConcurrent <= 0 {
		return nil, errors.New("max concurrent must be greater than 0")
	}

	return &Semaphore{tokens: make(chan struct{}, maxConcurrent)}, nil
}

func (semaphore *Semaphore) Acquire(ctx context.Context) error {
	if semaphore == nil {
		return errors.New("semaphore is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case semaphore.tokens <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (semaphore *Semaphore) TryAcquire() bool {
	if semaphore == nil {
		return false
	}

	select {
	case semaphore.tokens <- struct{}{}:
		return true
	default:
		return false
	}
}

func (semaphore *Semaphore) Release() error {
	if semaphore == nil {
		return errors.New("semaphore is required")
	}

	select {
	case <-semaphore.tokens:
		return nil
	default:
		return errors.New("semaphore has no acquired permits")
	}
}

func main() {
	semaphore, err := NewSemaphore(2)
	if err != nil {
		log.Fatal(err)
	}

	jobs := []string{"job-1", "job-2", "job-3"}
	var waitGroup sync.WaitGroup

	for _, job := range jobs {
		waitGroup.Add(1)
		if err := semaphore.Acquire(context.Background()); err != nil {
			log.Fatal(err)
		}

		go func(job string) {
			defer waitGroup.Done()
			defer func() {
				if releaseErr := semaphore.Release(); releaseErr != nil {
					log.Printf("release error: %v", releaseErr)
				}
			}()

			log.Printf("processing %s", job)
			time.Sleep(50 * time.Millisecond)
		}(job)
	}

	waitGroup.Wait()
}
