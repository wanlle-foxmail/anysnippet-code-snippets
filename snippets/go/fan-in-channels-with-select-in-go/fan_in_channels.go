package main

import (
	"context"
	"log"
)

// Flow:
//
//	read from left, right, or ctx.Done()
//	   |
//	   +-> input closed -> disable that branch and keep draining the other side
//	   +-> value received -> forward it downstream
//	   `-> context canceled or both inputs drained -> close output and stop
func FanIn[T any](ctx context.Context, left, right <-chan T) <-chan T {
	if ctx == nil {
		ctx = context.Background()
	}

	output := make(chan T)
	go func() {
		defer close(output)

		leftChannel := left
		rightChannel := right

		for leftChannel != nil || rightChannel != nil {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-leftChannel:
				if !ok {
					leftChannel = nil
					continue
				}

				select {
				case output <- value:
				case <-ctx.Done():
					return
				}
			case value, ok := <-rightChannel:
				if !ok {
					rightChannel = nil
					continue
				}

				select {
				case output <- value:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return output
}

func main() {
	left := make(chan string, 2)
	right := make(chan string, 2)

	left <- "left: one"
	right <- "right: one"
	left <- "left: two"
	close(left)
	close(right)

	for value := range FanIn(context.Background(), left, right) {
		log.Printf("value=%s", value)
	}
}
