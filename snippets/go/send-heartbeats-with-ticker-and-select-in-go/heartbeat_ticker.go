package main

import (
	"context"
	"errors"
	"log"
	"time"
)

// Flow:
//
//	wait for ticker or ctx.Done()
//	   |
//	   +-> tick -> try to queue one heartbeat without blocking
//	   `-> context canceled -> stop the ticker, close output, and return
func StartHeartbeats(ctx context.Context, interval time.Duration) (<-chan time.Time, error) {
	if interval <= 0 {
		return nil, errors.New("interval must be greater than 0")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	output := make(chan time.Time, 1)
	go func() {
		defer close(output)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case tickTime := <-ticker.C:
				select {
				case output <- tickTime:
				default:
				}
			}
		}
	}()

	return output, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeats, err := StartHeartbeats(ctx, 100*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	for beatNumber := 0; beatNumber < 3; beatNumber++ {
		tickTime := <-heartbeats
		log.Printf("heartbeat at %s", tickTime.Format(time.RFC3339Nano))
	}
}
