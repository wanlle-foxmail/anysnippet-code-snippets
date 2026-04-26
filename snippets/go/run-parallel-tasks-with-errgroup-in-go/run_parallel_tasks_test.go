package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunParallelTasksSucceedsWhenAllTasksSucceed(t *testing.T) {
	var completed atomic.Int32
	tasks := []Task{
		func(context.Context) error {
			completed.Add(1)
			return nil
		},
		func(context.Context) error {
			completed.Add(1)
			return nil
		},
	}

	if err := RunParallelTasks(context.Background(), tasks); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if completed.Load() != 2 {
		t.Fatalf("expected two completed tasks, got %d", completed.Load())
	}
}

func TestRunParallelTasksReturnsTheFirstError(t *testing.T) {
	wantErr := errors.New("task failed")
	tasks := []Task{
		func(context.Context) error {
			return wantErr
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	}

	err := RunParallelTasks(context.Background(), tasks)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
}

func TestRunParallelTasksCancelsSiblingTasks(t *testing.T) {
	wantErr := errors.New("task failed")
	canceled := make(chan struct{})
	tasks := []Task{
		func(context.Context) error {
			return wantErr
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			close(canceled)
			return nil
		},
	}

	err := RunParallelTasks(context.Background(), tasks)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got %v", wantErr, err)
	}
	select {
	case <-canceled:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected the sibling task to observe cancellation")
	}
}

func TestRunParallelTasksUsesBackgroundWhenContextIsNil(t *testing.T) {
	if err := RunParallelTasks(nil, []Task{func(context.Context) error { return nil }}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestRunParallelTasksRejectsNilTasks(t *testing.T) {
	err := RunParallelTasks(context.Background(), []Task{nil})
	if err == nil {
		t.Fatal("expected error for a nil task")
	}
}

func TestRunParallelTasksAllowsEmptyTaskLists(t *testing.T) {
	if err := RunParallelTasks(context.Background(), nil); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
