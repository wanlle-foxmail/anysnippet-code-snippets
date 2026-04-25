package main

import (
	"context"
	"testing"
	"time"
)

func TestCoalesceEventsRejectsNilInputChannel(t *testing.T) {
	if _, err := CoalesceEvents[string](context.Background(), nil, 20*time.Millisecond); err == nil {
		t.Fatal("expected error for nil input channel")
	}
}

func TestCoalesceEventsRejectsInvalidWindow(t *testing.T) {
	input := make(chan string)
	if _, err := CoalesceEvents(context.Background(), input, 0); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestCoalesceEventsMergesQuickBurstsIntoOneBatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan string, 2)
	output, err := CoalesceEvents(ctx, input, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("coalesce events: %v", err)
	}

	input <- "save"
	input <- "save"

	select {
	case batch := <-output:
		if len(batch) != 2 || batch[0] != "save" || batch[1] != "save" {
			t.Fatalf("expected one merged batch, got %v", batch)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for merged batch")
	}
}

func TestCoalesceEventsSplitsSeparatedBursts(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan string, 2)
	output, err := CoalesceEvents(ctx, input, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("coalesce events: %v", err)
	}

	input <- "first"
	firstBatch := <-output
	if len(firstBatch) != 1 || firstBatch[0] != "first" {
		t.Fatalf("expected first batch [first], got %v", firstBatch)
	}

	input <- "second"
	secondBatch := <-output
	if len(secondBatch) != 1 || secondBatch[0] != "second" {
		t.Fatalf("expected second batch [second], got %v", secondBatch)
	}
}

func TestCoalesceEventsFlushesPendingBatchWhenInputCloses(t *testing.T) {
	input := make(chan string, 1)
	output, err := CoalesceEvents(context.Background(), input, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("coalesce events: %v", err)
	}

	input <- "save"
	close(input)

	batch, ok := <-output
	if !ok {
		t.Fatal("expected one flushed batch before output close")
	}
	if len(batch) != 1 || batch[0] != "save" {
		t.Fatalf("expected [save], got %v", batch)
	}
	if _, ok := <-output; ok {
		t.Fatal("expected output to close after flushing the pending batch")
	}
}

func TestCoalesceEventsClosesOutputForEmptyInput(t *testing.T) {
	input := make(chan string)
	close(input)

	output, err := CoalesceEvents(context.Background(), input, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("coalesce events: %v", err)
	}
	if _, ok := <-output; ok {
		t.Fatal("expected output to close when the input is already empty")
	}
}

func TestCoalesceEventsStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	input := make(chan string)
	output, err := CoalesceEvents(ctx, input, 30*time.Millisecond)
	if err != nil {
		t.Fatalf("coalesce events: %v", err)
	}

	cancel()

	select {
	case _, ok := <-output:
		if ok {
			t.Fatal("expected output to close after context cancellation")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for output to close after context cancellation")
	}
}
