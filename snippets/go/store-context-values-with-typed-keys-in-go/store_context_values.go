package main

import (
	"context"
	"errors"
	"log"
	"strings"
)

type ContextKey[T any] struct {
	name string
}

func NewContextKey[T any](name string) (*ContextKey[T], error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("key name is required")
	}
	return &ContextKey[T]{name: name}, nil
}

// Flow: create one typed key -> attach a typed value -> read it back without colliding with other context keys.
func WithContextValue[T any](ctx context.Context, key *ContextKey[T], value T) (context.Context, error) {
	if key == nil {
		return nil, errors.New("key is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, key, value), nil
}

func ContextValue[T any](ctx context.Context, key *ContextKey[T]) (T, bool) {
	var zero T
	if ctx == nil || key == nil {
		return zero, false
	}

	rawValue := ctx.Value(key)
	if rawValue == nil {
		return zero, false
	}

	value, ok := rawValue.(T)
	if !ok {
		return zero, false
	}
	return value, true
}

func main() {
	requestIDKey, err := NewContextKey[string]("request_id")
	if err != nil {
		log.Fatal(err)
	}

	ctx, err := WithContextValue(context.Background(), requestIDKey, "req-42")
	if err != nil {
		log.Fatal(err)
	}

	requestID, ok := ContextValue(ctx, requestIDKey)
	if !ok {
		log.Fatal("request id missing")
	}

	log.Printf("request_id=%s", requestID)
}
