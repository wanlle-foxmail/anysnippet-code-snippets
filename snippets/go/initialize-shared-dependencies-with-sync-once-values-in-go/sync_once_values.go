package main

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

type AppSettings struct {
	DSN string
}

// Flow:
//   first caller runs the initializer once
//      |
//      +-> success -> cache the value for every later caller
//      `-> failure or panic -> replay the same result for every later caller
func NewOnceLoader[T any](load func() (T, error)) func() (T, error) {
	if load == nil {
		return func() (T, error) {
			var zero T
			return zero, errors.New("load function is required")
		}
	}
	return sync.OnceValues(load)
}

func main() {
	var calls atomic.Int32
	loadSettings := NewOnceLoader(func() (AppSettings, error) {
		calls.Add(1)
		return AppSettings{DSN: "postgres://app@db/main"}, nil
	})

	settings, err := loadSettings()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("dsn=%s calls=%d", settings.DSN, calls.Load())
}