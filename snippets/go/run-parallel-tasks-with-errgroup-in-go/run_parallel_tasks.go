package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type Task func(context.Context) error

// Flow:
//
//	start every task in one errgroup
//	   |
//	   +-> all tasks return nil -> wait returns nil
//	   `-> one task returns an error -> errgroup cancels siblings and returns that error
func RunParallelTasks(ctx context.Context, tasks []Task) error {
	if ctx == nil {
		ctx = context.Background()
	}
	for index, task := range tasks {
		if task == nil {
			return fmt.Errorf("task %d is required", index)
		}
	}

	group, groupContext := errgroup.WithContext(ctx)
	for _, task := range tasks {
		task := task
		group.Go(func() error {
			return task(groupContext)
		})
	}

	return group.Wait()
}

func main() {
	tasks := []Task{
		func(ctx context.Context) error {
			select {
			case <-time.After(20 * time.Millisecond):
				log.Print("task one finished")
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
		func(ctx context.Context) error {
			select {
			case <-time.After(10 * time.Millisecond):
				log.Print("task two finished")
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	if err := RunParallelTasks(context.Background(), tasks); err != nil {
		log.Fatal(err)
	}
}
