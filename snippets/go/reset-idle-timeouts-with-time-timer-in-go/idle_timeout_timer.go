package main

import (
	"context"
	"errors"
	"log"
	"time"
)

var ErrIdleTimeout = errors.New("idle timeout exceeded")

// Flow:
//
//	wait for activity, timer, or ctx.Done()
//	   |
//	   +-> activity -> stop and reset the idle timer
//	   +-> timer fires -> return ErrIdleTimeout
//	   `-> input closes or context cancels -> return without waiting again
func WaitForIdleTimeout(ctx context.Context, activity <-chan struct{}, idleTimeout time.Duration) error {
	if idleTimeout <= 0 {
		return errors.New("idle timeout must be greater than 0")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	timer := time.NewTimer(idleTimeout)
	defer stopAndDrainTimer(timer)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return ErrIdleTimeout
		case _, ok := <-activity:
			if !ok {
				return nil
			}

			stopAndDrainTimer(timer)
			timer.Reset(idleTimeout)
		}
	}
}

func stopAndDrainTimer(timer *time.Timer) {
	if timer == nil {
		return
	}
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func main() {
	activity := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		activity <- struct{}{}
		time.Sleep(100 * time.Millisecond)
		close(activity)
	}()

	err := WaitForIdleTimeout(context.Background(), activity, 150*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("finished before idle timeout")
}
