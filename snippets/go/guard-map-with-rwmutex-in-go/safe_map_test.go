package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestSafeMapZeroValueIsUsable(t *testing.T) {
	var safeMap SafeMap[string, int]
	safeMap.Set("count", 3)

	value, ok := safeMap.Get("count")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != 3 {
		t.Fatalf("expected 3, got %d", value)
	}
}

func TestSafeMapGetMissingKeyReturnsFalse(t *testing.T) {
	safeMap := NewSafeMap[string, int]()

	if _, ok := safeMap.Get("missing"); ok {
		t.Fatal("expected missing key lookup to return false")
	}
}

func TestSafeMapDeleteRemovesKey(t *testing.T) {
	safeMap := NewSafeMap[string, int]()
	safeMap.Set("count", 3)
	safeMap.Delete("count")

	if _, ok := safeMap.Get("count"); ok {
		t.Fatal("expected deleted key to be missing")
	}
}

func TestSafeMapLenTracksStoredKeys(t *testing.T) {
	safeMap := NewSafeMap[string, int]()
	safeMap.Set("a", 1)
	safeMap.Set("b", 2)
	safeMap.Set("c", 3)

	if safeMap.Len() != 3 {
		t.Fatalf("expected length 3, got %d", safeMap.Len())
	}
}

func TestSafeMapSnapshotReturnsIndependentCopy(t *testing.T) {
	safeMap := NewSafeMap[string, int]()
	safeMap.Set("count", 3)

	snapshot := safeMap.Snapshot()
	snapshot["count"] = 99

	value, ok := safeMap.Get("count")
	if !ok {
		t.Fatal("expected original key to exist")
	}
	if value != 3 {
		t.Fatalf("expected original value 3, got %d", value)
	}
}

func TestSafeMapSupportsConcurrentWrites(t *testing.T) {
	safeMap := NewSafeMap[string, int]()

	var waitGroup sync.WaitGroup
	for index := 0; index < 100; index++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			safeMap.Set(strconv.Itoa(index), index)
		}(index)
	}
	waitGroup.Wait()

	if safeMap.Len() != 100 {
		t.Fatalf("expected 100 keys, got %d", safeMap.Len())
	}
}

func TestSafeMapSupportsConcurrentReadsDuringWrites(t *testing.T) {
	safeMap := NewSafeMap[string, int]()

	var waitGroup sync.WaitGroup
	for index := 0; index < 100; index++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			key := strconv.Itoa(index)
			safeMap.Set(key, index)
			_, _ = safeMap.Get(key)
		}(index)
	}
	waitGroup.Wait()

	if safeMap.Len() != 100 {
		t.Fatalf("expected 100 keys after mixed reads and writes, got %d", safeMap.Len())
	}
}
