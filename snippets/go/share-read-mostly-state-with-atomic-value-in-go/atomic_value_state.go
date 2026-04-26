package main

import (
	"errors"
	"log"
	"sync/atomic"
)

type RoutingSnapshot struct {
	Version int
	Primary string
}

type AtomicSnapshot[T any] struct {
	value atomic.Value
}

func NewAtomicSnapshot[T any](initial T) *AtomicSnapshot[T] {
	holder := &AtomicSnapshot[T]{}
	holder.value.Store(initial)
	return holder
}

// Flow: store complete immutable snapshots -> load the latest snapshot without taking read locks.
func (holder *AtomicSnapshot[T]) Load() (T, error) {
	var zero T
	if holder == nil {
		return zero, errors.New("snapshot holder is required")
	}

	rawValue := holder.value.Load()
	if rawValue == nil {
		return zero, errors.New("snapshot holder is not initialized")
	}

	value, ok := rawValue.(T)
	if !ok {
		return zero, errors.New("snapshot holder has an unexpected type")
	}
	return value, nil
}

func (holder *AtomicSnapshot[T]) Store(next T) error {
	if holder == nil {
		return errors.New("snapshot holder is required")
	}
	holder.value.Store(next)
	return nil
}

func main() {
	routes := NewAtomicSnapshot(RoutingSnapshot{Version: 1, Primary: "api-a"})
	if err := routes.Store(RoutingSnapshot{Version: 2, Primary: "api-b"}); err != nil {
		log.Fatal(err)
	}

	current, err := routes.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("version=%d primary=%s", current.Version, current.Primary)
}