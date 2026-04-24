package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type countingJSONLWriter struct {
	*httptest.ResponseRecorder
	flushCount int
}

func (writer *countingJSONLWriter) Flush() {
	writer.flushCount++
	writer.ResponseRecorder.Flush()
}

type plainJSONLWriter struct {
	header http.Header
}

func (writer *plainJSONLWriter) Header() http.Header {
	return writer.header
}

func (writer *plainJSONLWriter) Write(buffer []byte) (int, error) {
	return len(buffer), nil
}

func (writer *plainJSONLWriter) WriteHeader(statusCode int) {
}

type failingJSONLWriter struct {
	header http.Header
}

func (writer *failingJSONLWriter) Header() http.Header {
	return writer.header
}

func (writer *failingJSONLWriter) Write(buffer []byte) (int, error) {
	return 0, errors.New("write failed")
}

func (writer *failingJSONLWriter) WriteHeader(statusCode int) {
}

func (writer *failingJSONLWriter) Flush() {
}

func TestStreamJSONLWritesHeadersAndRecords(t *testing.T) {
	records := make(chan interface{}, 2)
	records <- map[string]interface{}{"id": 1, "status": "ready"}
	records <- map[string]interface{}{"id": 2, "status": "ready"}
	close(records)

	recorder := &countingJSONLWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/records", nil)

	err := StreamJSONL(recorder, req, records)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if recorder.Header().Get("Content-Type") != jsonlContentType {
		t.Fatalf("expected content type %q, got %q", jsonlContentType, recorder.Header().Get("Content-Type"))
	}
	if recorder.Body.String() != "{\"id\":1,\"status\":\"ready\"}\n{\"id\":2,\"status\":\"ready\"}\n" {
		t.Fatalf("unexpected JSONL body %q", recorder.Body.String())
	}
}

func TestStreamJSONLFlushesAfterEachRecord(t *testing.T) {
	records := make(chan interface{}, 2)
	records <- map[string]interface{}{"id": 1}
	records <- map[string]interface{}{"id": 2}
	close(records)

	recorder := &countingJSONLWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/records", nil)

	err := StreamJSONL(recorder, req, records)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if recorder.flushCount != 2 {
		t.Fatalf("expected 2 flushes, got %d", recorder.flushCount)
	}
}

func TestStreamJSONLReturnsNilWhenClientDisconnects(t *testing.T) {
	records := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/records", nil).WithContext(ctx)
	recorder := &countingJSONLWriter{ResponseRecorder: httptest.NewRecorder()}

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- StreamJSONL(recorder, req, records)
	}()

	cancel()

	select {
	case err := <-resultCh:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected StreamJSONL to stop after client disconnect")
	}
}

func TestStreamJSONLReturnsNilWhenTheChannelCloses(t *testing.T) {
	records := make(chan interface{})
	close(records)

	recorder := &countingJSONLWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/records", nil)

	err := StreamJSONL(recorder, req, records)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestStreamJSONLRejectsWritersWithoutStreamingSupport(t *testing.T) {
	records := make(chan interface{})
	close(records)
	req := httptest.NewRequest(http.MethodGet, "/records", nil)

	err := StreamJSONL(&plainJSONLWriter{header: make(http.Header)}, req, records)
	if err == nil {
		t.Fatal("expected missing flusher error")
	}
}

func TestStreamJSONLPropagatesMarshalErrors(t *testing.T) {
	records := make(chan interface{}, 1)
	records <- func() {}
	close(records)

	recorder := &countingJSONLWriter{ResponseRecorder: httptest.NewRecorder()}
	req := httptest.NewRequest(http.MethodGet, "/records", nil)

	err := StreamJSONL(recorder, req, records)
	if err == nil {
		t.Fatal("expected marshal error")
	}
	if !strings.Contains(err.Error(), "marshal jsonl record") {
		t.Fatalf("expected wrapped marshal error, got %v", err)
	}
}

func TestStreamJSONLRejectsNilInputsAndPropagatesWriteErrors(t *testing.T) {
	records := make(chan interface{}, 1)
	records <- map[string]interface{}{"id": 1}
	close(records)
	req := httptest.NewRequest(http.MethodGet, "/records", nil)
	recorder := &countingJSONLWriter{ResponseRecorder: httptest.NewRecorder()}

	if err := StreamJSONL(nil, req, records); err == nil {
		t.Fatal("expected nil writer error")
	}
	if err := StreamJSONL(recorder, nil, records); err == nil {
		t.Fatal("expected nil request error")
	}
	if err := StreamJSONL(recorder, req, nil); err == nil {
		t.Fatal("expected nil records channel error")
	}

	failingWriter := &failingJSONLWriter{header: make(http.Header)}
	err := StreamJSONL(failingWriter, req, records)
	if err == nil {
		t.Fatal("expected write error")
	}
	if !strings.Contains(err.Error(), "write jsonl record") {
		t.Fatalf("expected wrapped write error, got %v", err)
	}
}