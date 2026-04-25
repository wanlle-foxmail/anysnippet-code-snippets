package main

import (
	"context"
	"errors"
	"log"
	"time"
)

// Flow:
//
//	read input, timer, or ctx.Done()
//	   |
//	   +-> input value -> buffer it and start or reset the quiet-window timer
//	   +-> timer fires -> flush one batch downstream
//	   +-> input closes -> flush pending values and stop
//	   `-> context canceled -> stop without emitting more values
func CoalesceEvents[T any](ctx context.Context, input <-chan T, window time.Duration) (<-chan []T, error) {
	if input == nil {
		return nil, errors.New("input channel is required")
	}
	if window <= 0 {
		return nil, errors.New("window must be greater than 0")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	output := make(chan []T)
	go func() {
		defer close(output)

		var buffer []T
		var timer *time.Timer
		var timerChannel <-chan time.Time

		flush := func() bool {
			if len(buffer) == 0 {
				return true
			}

			batch := append([]T(nil), buffer...)
			buffer = nil

			select {
			case output <- batch:
				return true
			case <-ctx.Done():
				return false
			}
		}

		for {
			select {
			case <-ctx.Done():
				stopAndDrainCoalesceTimer(timer)
				return
			case value, ok := <-input:
				if !ok {
					stopAndDrainCoalesceTimer(timer)
					_ = flush()
					return
				}

				buffer = append(buffer, value)
				if timer == nil {
					timer = time.NewTimer(window)
					timerChannel = timer.C
					continue
				}

				stopAndDrainCoalesceTimer(timer)
				timer.Reset(window)
			case <-timerChannel:
				if !flush() {
					stopAndDrainCoalesceTimer(timer)
					return
				}
				timer = nil
				timerChannel = nil
			}
		}
	}()

	return output, nil
}

func stopAndDrainCoalesceTimer(timer *time.Timer) {
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan string, 3)
	output, err := CoalesceEvents(ctx, input, 75*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	input <- "save"
	input <- "save"
	input <- "save"
	close(input)

	for batch := range output {
		log.Printf("batch=%v", batch)
	}
}
