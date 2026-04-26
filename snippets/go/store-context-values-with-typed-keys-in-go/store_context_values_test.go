package main

import (
	"context"
	"testing"
)

func TestWithContextValueStoresAndReadsOneTypedValue(t *testing.T) {
	requestIDKey, err := NewContextKey[string]("request_id")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	ctx, err := WithContextValue(context.Background(), requestIDKey, "req-42")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	requestID, ok := ContextValue(ctx, requestIDKey)
	if !ok || requestID != "req-42" {
		t.Fatalf("expected request_id req-42, got %q ok=%v", requestID, ok)
	}
}

func TestContextValueReturnsFalseWhenMissing(t *testing.T) {
	requestIDKey, err := NewContextKey[string]("request_id")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	requestID, ok := ContextValue(context.Background(), requestIDKey)
	if ok || requestID != "" {
		t.Fatalf("expected missing value, got %q ok=%v", requestID, ok)
	}
}

func TestWithContextValueRejectsNilKeys(t *testing.T) {
	_, err := WithContextValue(context.Background(), nil, "req-42")
	if err == nil {
		t.Fatal("expected error for a nil key")
	}
}

func TestWithContextValueUsesBackgroundWhenContextIsNil(t *testing.T) {
	requestIDKey, err := NewContextKey[string]("request_id")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	ctx, err := WithContextValue[string](nil, requestIDKey, "req-42")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	requestID, ok := ContextValue(ctx, requestIDKey)
	if !ok || requestID != "req-42" {
		t.Fatalf("expected request_id req-42, got %q ok=%v", requestID, ok)
	}
}

func TestContextValueSeparatesDistinctTypedKeys(t *testing.T) {
	requestIDKey, err := NewContextKey[string]("request_id")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	otherRequestIDKey, err := NewContextKey[string]("request_id")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	ctx, err := WithContextValue(context.Background(), requestIDKey, "req-42")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	requestID, ok := ContextValue(ctx, otherRequestIDKey)
	if ok || requestID != "" {
		t.Fatalf("expected the second key to miss, got %q ok=%v", requestID, ok)
	}
}

func TestNewContextKeyRejectsBlankNames(t *testing.T) {
	if _, err := NewContextKey[string]("  "); err == nil {
		t.Fatal("expected error for a blank key name")
	}
}