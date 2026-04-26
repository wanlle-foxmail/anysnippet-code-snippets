package main

import (
	"sync"
	"testing"
)

func TestNewAtomicSnapshotLoadsTheInitialValue(t *testing.T) {
	routes := NewAtomicSnapshot(RoutingSnapshot{Version: 1, Primary: "api-a"})
	current, err := routes.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if current.Version != 1 || current.Primary != "api-a" {
		t.Fatalf("unexpected snapshot %+v", current)
	}
}

func TestAtomicSnapshotStoresUpdatedSnapshots(t *testing.T) {
	routes := NewAtomicSnapshot(RoutingSnapshot{Version: 1, Primary: "api-a"})
	if err := routes.Store(RoutingSnapshot{Version: 2, Primary: "api-b"}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	current, err := routes.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if current.Version != 2 || current.Primary != "api-b" {
		t.Fatalf("unexpected snapshot %+v", current)
	}
}

func TestAtomicSnapshotSupportsConcurrentLoadsAndStores(t *testing.T) {
	routes := NewAtomicSnapshot(RoutingSnapshot{Version: 1, Primary: "api-a"})

	var writers sync.WaitGroup
	for version := 2; version <= 8; version++ {
		writers.Add(1)
		go func(version int) {
			defer writers.Done()
			_ = routes.Store(RoutingSnapshot{Version: version, Primary: "api-b"})
		}(version)
	}

	var readers sync.WaitGroup
	for range 20 {
		readers.Add(1)
		go func() {
			defer readers.Done()
			if _, err := routes.Load(); err != nil {
				t.Errorf("expected nil error, got %v", err)
			}
		}()
	}

	writers.Wait()
	readers.Wait()

	current, err := routes.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if current.Version < 1 || current.Version > 8 {
		t.Fatalf("expected a written version between 1 and 8, got %+v", current)
	}
}

func TestAtomicSnapshotRejectsNilReceivers(t *testing.T) {
	var routes *AtomicSnapshot[RoutingSnapshot]
	if _, err := routes.Load(); err == nil {
		t.Fatal("expected error for a nil receiver on Load")
	}
	if err := routes.Store(RoutingSnapshot{Version: 1, Primary: "api-a"}); err == nil {
		t.Fatal("expected error for a nil receiver on Store")
	}
}

func TestNewAtomicSnapshotSupportsZeroValueSnapshots(t *testing.T) {
	routes := NewAtomicSnapshot(RoutingSnapshot{})
	current, err := routes.Load()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if current != (RoutingSnapshot{}) {
		t.Fatalf("expected a zero-value snapshot, got %+v", current)
	}
}