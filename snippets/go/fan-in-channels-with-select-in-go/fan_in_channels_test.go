package main

import (
	"context"
	"sort"
	"testing"
	"time"
)

func TestFanInMergesValuesFromBothChannels(t *testing.T) {
	left := make(chan int, 2)
	right := make(chan int, 2)

	left <- 1
	left <- 3
	right <- 2
	right <- 4
	close(left)
	close(right)

	values := make([]int, 0, 4)
	for value := range FanIn(context.Background(), left, right) {
		values = append(values, value)
	}

	sort.Ints(values)
	want := []int{1, 2, 3, 4}
	if len(values) != len(want) {
		t.Fatalf("expected %d values, got %d", len(want), len(values))
	}
	for index, value := range want {
		if values[index] != value {
			t.Fatalf("expected values %v, got %v", want, values)
		}
	}
}

func TestFanInClosesOutputWhenInputsClose(t *testing.T) {
	left := make(chan int)
	right := make(chan int)
	close(left)
	close(right)

	_, ok := <-FanIn(context.Background(), left, right)
	if ok {
		t.Fatal("expected output channel to be closed")
	}
}

func TestFanInDrainsLeftWhenRightIsNil(t *testing.T) {
	left := make(chan string, 2)
	left <- "a"
	left <- "b"
	close(left)

	values := make([]string, 0, 2)
	for value := range FanIn[string](context.Background(), left, nil) {
		values = append(values, value)
	}

	sort.Strings(values)
	if len(values) != 2 || values[0] != "a" || values[1] != "b" {
		t.Fatalf("expected [a b], got %v", values)
	}
}

func TestFanInDrainsRightWhenLeftIsNil(t *testing.T) {
	right := make(chan string, 2)
	right <- "a"
	right <- "b"
	close(right)

	values := make([]string, 0, 2)
	for value := range FanIn[string](context.Background(), nil, right) {
		values = append(values, value)
	}

	sort.Strings(values)
	if len(values) != 2 || values[0] != "a" || values[1] != "b" {
		t.Fatalf("expected [a b], got %v", values)
	}
}

func TestFanInReturnsClosedChannelForNilInputs(t *testing.T) {
	_, ok := <-FanIn[int](context.Background(), nil, nil)
	if ok {
		t.Fatal("expected output channel to be closed")
	}
}

func TestFanInStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	left := make(chan int, 1)
	left <- 1

	output := FanIn(ctx, left, nil)
	if value, ok := <-output; !ok || value != 1 {
		t.Fatalf("expected first forwarded value to be 1, got value=%d ok=%t", value, ok)
	}

	cancel()

	select {
	case _, ok := <-output:
		if ok {
			t.Fatal("expected output channel to close after context cancellation")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for output channel to close after context cancellation")
	}
}
