package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func waitForHTTPResponse(t *testing.T, url string) string {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		response, err := http.Get(url)
		if err == nil {
			defer response.Body.Close()
			body, readErr := io.ReadAll(response.Body)
			if readErr != nil {
				t.Fatalf("read response body: %v", readErr)
			}
			if response.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", response.StatusCode)
			}
			return string(body)
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("server did not become ready in time")
	return ""
}

func TestServeWithGracefulShutdownServesRequestsUntilContextCanceled(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	server := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("ok"))
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- ServeWithGracefulShutdown(ctx, server, listener, 500*time.Millisecond)
	}()

	body := waitForHTTPResponse(t, "http://"+listener.Addr().String())
	if body != "ok" {
		t.Fatalf("expected body ok, got %q", body)
	}

	cancel()

	if err := <-resultCh; err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestServeWithGracefulShutdownWaitsForInflightRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	requestStarted := make(chan struct{})
	allowFinish := make(chan struct{})

	server := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			close(requestStarted)
			<-allowFinish
			_, _ = writer.Write([]byte("done"))
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- ServeWithGracefulShutdown(ctx, server, listener, time.Second)
	}()

	responseCh := make(chan error, 1)
	go func() {
		response, getErr := http.Get("http://" + listener.Addr().String())
		if getErr != nil {
			responseCh <- getErr
			return
		}
		defer response.Body.Close()
		_, _ = io.ReadAll(response.Body)
		responseCh <- nil
	}()

	<-requestStarted
	cancel()
	close(allowFinish)

	if err := <-responseCh; err != nil {
		t.Fatalf("expected in-flight request to finish, got %v", err)
	}
	if err := <-resultCh; err != nil {
		t.Fatalf("expected graceful shutdown, got %v", err)
	}
}

func TestServeWithGracefulShutdownReturnsTimeoutErrorWhenRequestDoesNotFinish(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	requestStarted := make(chan struct{})
	allowFinish := make(chan struct{})

	server := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			close(requestStarted)
			<-allowFinish
			_, _ = writer.Write([]byte("late"))
		}),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- ServeWithGracefulShutdown(ctx, server, listener, 20*time.Millisecond)
	}()

	go func() {
		_, _ = http.Get("http://" + listener.Addr().String())
	}()

	<-requestStarted
	cancel()

	err = <-resultCh
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}

	close(allowFinish)
	_ = server.Close()
}

func TestServeWithGracefulShutdownReturnsServeErrorFromClosedListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	server := &http.Server{Handler: http.NewServeMux()}

	err = ServeWithGracefulShutdown(context.Background(), server, listener, time.Second)
	if err == nil {
		t.Fatal("expected serve error")
	}
}

func TestServeWithGracefulShutdownRejectsNilServer(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	err = ServeWithGracefulShutdown(context.Background(), nil, listener, time.Second)
	if err == nil {
		t.Fatal("expected error for nil server")
	}
}

func TestServeWithGracefulShutdownRejectsNilListenerOrNonPositiveTimeout(t *testing.T) {
	server := &http.Server{Handler: http.NewServeMux()}

	err := ServeWithGracefulShutdown(context.Background(), server, nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil listener")
	}

	listener, listenErr := net.Listen("tcp", "127.0.0.1:0")
	if listenErr != nil {
		t.Fatalf("listen: %v", listenErr)
	}
	defer listener.Close()

	err = ServeWithGracefulShutdown(context.Background(), server, listener, 0)
	if err == nil {
		t.Fatal("expected error for non-positive timeout")
	}
}
