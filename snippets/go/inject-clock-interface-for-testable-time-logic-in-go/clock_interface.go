package main

import (
	"errors"
	"log"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func HasExpired(clock Clock, expiresAt time.Time) (bool, error) {
	if clock == nil {
		return false, errors.New("clock is required")
	}
	if expiresAt.IsZero() {
		return false, errors.New("expiration time is required")
	}

	return !clock.Now().Before(expiresAt), nil
}

func main() {
	expiresAt := time.Now().Add(30 * time.Minute)
	expired, err := HasExpired(RealClock{}, expiresAt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("expired=%t", expired)
}
