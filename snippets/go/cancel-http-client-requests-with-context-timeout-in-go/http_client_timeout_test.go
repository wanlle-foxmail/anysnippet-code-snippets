package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type blockingTransport struct{}

func (blockingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	<-request.Context().Done()
	return nil, request.Context().Err()
}

func TestGetWithContextTimeoutReturnsSuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	}))
	defer server.Close()

	response, err := GetWithContextTimeout(context.Background(), server.Client(), server.URL, time.Second)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		t.Fatalf("read body: %v", readErr)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
	if string(body) != "ok" {
		t.Fatalf("expected body ok, got %q", string(body))
	}
}

func TestGetWithContextTimeoutCancelsSlowRequest(t *testing.T) {
	client := &http.Client{Transport: blockingTransport{}}

	_, err := GetWithContextTimeout(context.Background(), client, "http://example.test", 20*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
}

func TestGetWithContextTimeoutHonorsParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &http.Client{Transport: blockingTransport{}}

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	_, err := GetWithContextTimeout(ctx, client, "http://example.test", time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestGetWithContextTimeoutRejectsNilClient(t *testing.T) {
	if _, err := GetWithContextTimeout(context.Background(), nil, "http://example.test", time.Second); err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestGetWithContextTimeoutRejectsEmptyURL(t *testing.T) {
	if _, err := GetWithContextTimeout(context.Background(), http.DefaultClient, "", time.Second); err == nil {
		t.Fatal("expected error for empty url")
	}
}

func TestGetWithContextTimeoutRejectsInvalidTimeout(t *testing.T) {
	if _, err := GetWithContextTimeout(context.Background(), http.DefaultClient, "http://example.test", 0); err == nil {
		t.Fatal("expected error for zero timeout")
	}
}
