package main

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestProcessBatchWithContextProcessesAllItems(t *testing.T) {
	items := []string{"a", "b", "c"}
	processed := make([]string, 0, len(items))

	err := ProcessBatchWithContext(context.Background(), items, func(ctx context.Context, item string) error {
		processed = append(processed, item)
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(processed) != len(items) {
		t.Fatalf("expected %d processed items, got %d", len(items), len(processed))
	}
}

func TestProcessBatchWithContextStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processed := make([]string, 0, 3)
	err := ProcessBatchWithContext(ctx, []string{"a", "b", "c"}, func(ctx context.Context, item string) error {
		processed = append(processed, item)
		cancel()
		return nil
	})
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(processed) != 1 {
		t.Fatalf("expected one processed item before cancellation, got %d", len(processed))
	}
}

func TestProcessBatchWithContextReturnsWorkerError(t *testing.T) {
	workerErr := errors.New("worker failed")
	err := ProcessBatchWithContext(context.Background(), []string{"a", "b"}, func(ctx context.Context, item string) error {
		if item == "b" {
			return workerErr
		}
		return nil
	})
	if !errors.Is(err, workerErr) {
		t.Fatalf("expected wrapped worker error, got %v", err)
	}
	if !strings.Contains(err.Error(), `process "b"`) {
		t.Fatalf("expected error to mention item b, got %q", err.Error())
	}
}

func TestProcessBatchWithContextRejectsNilProcessFunction(t *testing.T) {
	if err := ProcessBatchWithContext(context.Background(), []string{"a"}, nil); err == nil {
		t.Fatal("expected error for nil process function")
	}
}

func TestProcessBatchWithContextUsesBackgroundForNilContext(t *testing.T) {
	processed := 0
	err := ProcessBatchWithContext(nil, []string{"a", "b"}, func(ctx context.Context, item string) error {
		processed++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if processed != 2 {
		t.Fatalf("expected 2 processed items, got %d", processed)
	}
}

func TestProcessBatchWithContextHandlesEmptyItemList(t *testing.T) {
	called := false
	err := ProcessBatchWithContext(context.Background(), nil, func(ctx context.Context, item string) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if called {
		t.Fatal("expected process function to stay unused for an empty batch")
	}
}
