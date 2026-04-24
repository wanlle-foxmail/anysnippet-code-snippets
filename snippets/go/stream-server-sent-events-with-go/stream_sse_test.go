package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type countingFlushWriter struct {
	*httptest.ResponseRecorder
	flushCount int
}

func (writer *countingFlushWriter) Flush() {
	writer.flushCount++
	writer.ResponseRecorder.Flush()
}

type failingStreamWriter struct {
	header http.Header
}

type plainResponseWriter struct {
	header http.Header
}

func (writer *plainResponseWriter) Header() http.Header {
	return writer.header
}

func (writer *plainResponseWriter) Write(buffer []byte) (int, error) {
	return len(buffer), nil
}

func (writer *plainResponseWriter) WriteHeader(statusCode int) {
}

func (writer *failingStreamWriter) Header() http.Header {
	return writer.header
}

func (writer *failingStreamWriter) Write(buffer []byte) (int, error) {
	return 0, errors.New("write failed")
}

func (writer *failingStreamWriter) WriteHeader(statusCode int) {
}

func (writer *failingStreamWriter) Flush() {
}

func TestStreamSSEWritesHeadersAndFirstMessage(t *testing.T) {
	messages := make(chan string, 1)
	messages <- "hello"
	close(messages)

	recorder := &countingFlushWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	if err := StreamSSE(recorder, req, messages); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got := recorder.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("expected text/event-stream, got %q", got)
	}
	if !strings.Contains(recorder.Body.String(), "data: hello\n\n") {
		t.Fatalf("expected SSE body, got %q", recorder.Body.String())
	}
}

func TestStreamSSEFlushesAfterEachMessage(t *testing.T) {
	messages := make(chan string, 2)
	messages <- "first"
	messages <- "second"
	close(messages)

	recorder := &countingFlushWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	if err := StreamSSE(recorder, req, messages); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if recorder.flushCount != 2 {
		t.Fatalf("expected 2 flushes, got %d", recorder.flushCount)
	}
	if !strings.Contains(recorder.Body.String(), "data: first\n\ndata: second\n\n") {
		t.Fatalf("expected both SSE events, got %q", recorder.Body.String())
	}
}

func TestStreamSSEFormatsMultilineMessagesSafely(t *testing.T) {
	messages := make(chan string, 1)
	messages <- "first\nsecond"
	close(messages)

	recorder := &countingFlushWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	if err := StreamSSE(recorder, req, messages); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if recorder.Body.String() != "data: first\ndata: second\n\n" {
		t.Fatalf("expected multiline SSE block, got %q", recorder.Body.String())
	}
}

func TestStreamSSEReturnsNilWhenClientDisconnects(t *testing.T) {
	messages := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx)
	recorder := &countingFlushWriter{ResponseRecorder: httptest.NewRecorder()}

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- StreamSSE(recorder, req, messages)
	}()

	cancel()

	select {
	case err := <-resultCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected StreamSSE to stop after client disconnect")
	}
}

func TestStreamSSEReturnsNilWhenMessageChannelCloses(t *testing.T) {
	messages := make(chan string)
	close(messages)

	recorder := &countingFlushWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	if err := StreamSSE(recorder, req, messages); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestStreamSSERejectsWriterWithoutFlusher(t *testing.T) {
	messages := make(chan string)
	close(messages)
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	writer := &plainResponseWriter{header: make(http.Header)}

	if err := StreamSSE(writer, req, messages); err == nil {
		t.Fatal("expected missing flusher error")
	}
}

func TestStreamSSEPropagatesWriteErrors(t *testing.T) {
	messages := make(chan string, 1)
	messages <- "hello"
	close(messages)
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	writer := &failingStreamWriter{header: make(http.Header)}

	err := StreamSSE(writer, req, messages)
	if err == nil {
		t.Fatal("expected write error")
	}
	if !strings.Contains(err.Error(), "write sse event") {
		t.Fatalf("expected wrapped write error, got %v", err)
	}
}

func TestStreamSSERejectsNilInputs(t *testing.T) {
	messages := make(chan string)
	close(messages)
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	recorder := &countingFlushWriter{ResponseRecorder: httptest.NewRecorder()}

	if err := StreamSSE(nil, req, messages); err == nil {
		t.Fatal("expected nil writer error")
	}
	if err := StreamSSE(recorder, nil, messages); err == nil {
		t.Fatal("expected nil request error")
	}
	if err := StreamSSE(recorder, req, nil); err == nil {
		t.Fatal("expected nil channel error")
	}
	if bytes.Contains(recorder.Body.Bytes(), []byte("data:")) {
		t.Fatal("expected nil-input checks to run before streaming")
	}
}
