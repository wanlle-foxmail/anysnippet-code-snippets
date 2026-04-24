package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

type trackingReadCloser struct {
	reader io.Reader
	closed bool
}

func (closer *trackingReadCloser) Read(buffer []byte) (int, error) {
	return closer.reader.Read(buffer)
}

func (closer *trackingReadCloser) Close() error {
	closer.closed = true
	return nil
}

func TestRetryHTTPGetReturnsFirstSuccessfulResponse(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		_, _ = writer.Write([]byte("ok"))
	}))
	defer server.Close()

	response, err := RetryHTTPGet(server.Client(), server.URL, 3, 10*time.Millisecond)
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
		t.Fatalf("expected ok body, got %q", string(body))
	}
	if requestCount != 1 {
		t.Fatalf("expected one request, got %d", requestCount)
	}
}

func TestRetryHTTPGetRetriesOnServerErrorUntilSuccess(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		if requestCount == 1 {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("retry"))
			return
		}
		_, _ = writer.Write([]byte("ok"))
	}))
	defer server.Close()

	originalSleep := sleep
	sleepCalls := 0
	sleep = func(delay time.Duration) {
		sleepCalls++
	}
	defer func() {
		sleep = originalSleep
	}()

	response, err := RetryHTTPGet(server.Client(), server.URL, 3, 25*time.Millisecond)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
	if requestCount != 2 {
		t.Fatalf("expected two requests, got %d", requestCount)
	}
	if sleepCalls != 1 {
		t.Fatalf("expected one sleep call, got %d", sleepCalls)
	}
}

func TestRetryHTTPGetRetriesOnTooManyRequestsUntilSuccess(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		if requestCount == 1 {
			writer.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	response, err := RetryHTTPGet(server.Client(), server.URL, 2, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
	if requestCount != 2 {
		t.Fatalf("expected two requests, got %d", requestCount)
	}
}

func TestRetryHTTPGetDoesNotRetryOnBadRequest(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		writer.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	originalSleep := sleep
	sleepCalls := 0
	sleep = func(delay time.Duration) {
		sleepCalls++
	}
	defer func() {
		sleep = originalSleep
	}()

	response, err := RetryHTTPGet(server.Client(), server.URL, 3, time.Millisecond)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.StatusCode)
	}
	if requestCount != 1 {
		t.Fatalf("expected one request, got %d", requestCount)
	}
	if sleepCalls != 0 {
		t.Fatalf("expected no sleep calls, got %d", sleepCalls)
	}
}

func TestRetryHTTPGetRetriesTransportErrors(t *testing.T) {
	requestCount := 0
	client := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			requestCount++
			if requestCount == 1 {
				return nil, errors.New("temporary network failure")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
				Request:    request,
			}, nil
		}),
	}

	response, err := RetryHTTPGet(client, "http://example.test", 2, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	if requestCount != 2 {
		t.Fatalf("expected two attempts, got %d", requestCount)
	}
}

func TestRetryHTTPGetClosesRetryableResponseBodyBeforeRetrying(t *testing.T) {
	firstBody := &trackingReadCloser{reader: strings.NewReader("retry")}
	requestCount := 0
	client := &http.Client{
		Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			requestCount++
			if requestCount == 1 {
				return &http.Response{
					StatusCode: http.StatusServiceUnavailable,
					Body:       firstBody,
					Header:     make(http.Header),
					Request:    request,
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
				Header:     make(http.Header),
				Request:    request,
			}, nil
		}),
	}

	response, err := RetryHTTPGet(client, "http://example.test", 2, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	if !firstBody.closed {
		t.Fatal("expected retryable response body to be closed before retry")
	}
}

func TestRetryHTTPGetReturnsFinalRetryableResponseWhenAttemptsRunOut(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestCount++
		writer.WriteHeader(http.StatusServiceUnavailable)
		_, _ = writer.Write([]byte("still failing"))
	}))
	defer server.Close()

	response, err := RetryHTTPGet(server.Client(), server.URL, 2, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	defer response.Body.Close()

	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		t.Fatalf("read body: %v", readErr)
	}

	if response.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.StatusCode)
	}
	if string(body) != "still failing" {
		t.Fatalf("expected body still failing, got %q", string(body))
	}
	if requestCount != 2 {
		t.Fatalf("expected two attempts, got %d", requestCount)
	}
}

func TestRetryHTTPGetRejectsInvalidInput(t *testing.T) {
	if _, err := RetryHTTPGet(nil, "http://example.test", 2, 0); err == nil {
		t.Fatal("expected error for nil client")
	}
	if _, err := RetryHTTPGet(http.DefaultClient, "", 2, 0); err == nil {
		t.Fatal("expected error for empty url")
	}
	if _, err := RetryHTTPGet(http.DefaultClient, "http://example.test", 0, 0); err == nil {
		t.Fatal("expected error for invalid max attempts")
	}
	if _, err := RetryHTTPGet(http.DefaultClient, "http://example.test", 2, -time.Millisecond); err == nil {
		t.Fatal("expected error for negative delay")
	}
}
