package main

import (
	"log"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{}
}

func (safeMap *SafeMap[K, V]) Set(key K, value V) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	if safeMap.data == nil {
		safeMap.data = make(map[K]V)
	}
	safeMap.data[key] = value
}

func (safeMap *SafeMap[K, V]) Get(key K) (V, bool) {
	safeMap.mu.RLock()
	defer safeMap.mu.RUnlock()

	var zero V
	if safeMap.data == nil {
		return zero, false
	}

	value, ok := safeMap.data[key]
	return value, ok
}

func (safeMap *SafeMap[K, V]) Delete(key K) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	if safeMap.data == nil {
		return
	}
	delete(safeMap.data, key)
}

func (safeMap *SafeMap[K, V]) Len() int {
	safeMap.mu.RLock()
	defer safeMap.mu.RUnlock()

	return len(safeMap.data)
}

func (safeMap *SafeMap[K, V]) Snapshot() map[K]V {
	safeMap.mu.RLock()
	defer safeMap.mu.RUnlock()

	snapshot := make(map[K]V, len(safeMap.data))
	for key, value := range safeMap.data {
		snapshot[key] = value
	}

	return snapshot
}

func main() {
	cache := NewSafeMap[string, int]()
	cache.Set("active-users", 42)

	value, ok := cache.Get("active-users")
	if !ok {
		log.Fatal("expected key to exist")
	}

	log.Printf("active-users=%d", value)
}
