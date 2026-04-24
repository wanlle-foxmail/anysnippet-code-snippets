package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type failingResponseWriter struct {
	header http.Header
}

func (writer *failingResponseWriter) Header() http.Header {
	return writer.header
}

func (writer *failingResponseWriter) Write(buffer []byte) (int, error) {
	return 0, errors.New("write failed")
}

func (writer *failingResponseWriter) WriteHeader(statusCode int) {
}

func TestWriteJSONWithETagWritesTheBodyAndETag(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	rec := httptest.NewRecorder()

	err := WriteJSONWithETag(rec, req, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type application/json, got %q", rec.Header().Get("Content-Type"))
	}
	if rec.Header().Get("ETag") == "" {
		t.Fatal("expected ETag header to be set")
	}
	if !strings.Contains(rec.Body.String(), "product-42") {
		t.Fatalf("expected product body, got %q", rec.Body.String())
	}
}

func TestWriteJSONWithETagReturns304WhenTheETagMatches(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	firstRec := httptest.NewRecorder()

	err := WriteJSONWithETag(firstRec, req, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	secondReq := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	secondReq.Header.Set("If-None-Match", firstRec.Header().Get("ETag"))
	secondRec := httptest.NewRecorder()

	err = WriteJSONWithETag(secondRec, secondReq, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if secondRec.Code != http.StatusNotModified {
		t.Fatalf("expected status 304, got %d", secondRec.Code)
	}
	if secondRec.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", secondRec.Body.String())
	}
}

func TestWriteJSONWithETagIgnoresNonMatchingValidators(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	req.Header.Set("If-None-Match", `"other-tag"`)
	rec := httptest.NewRecorder()

	err := WriteJSONWithETag(rec, req, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestWriteJSONWithETagMatchesWildcardValidators(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	req.Header.Set("If-None-Match", "*")
	rec := httptest.NewRecorder()

	err := WriteJSONWithETag(rec, req, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if rec.Code != http.StatusNotModified {
		t.Fatalf("expected status 304, got %d", rec.Code)
	}
}

func TestWriteJSONWithETagMatchesOneValueInAValidatorList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	firstRec := httptest.NewRecorder()

	err := WriteJSONWithETag(firstRec, req, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	secondReq := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	secondReq.Header.Set("If-None-Match", `"other-tag", `+firstRec.Header().Get("ETag"))
	secondRec := httptest.NewRecorder()

	err = WriteJSONWithETag(secondRec, secondReq, map[string]string{"product_id": "product-42", "status": "ready"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if secondRec.Code != http.StatusNotModified {
		t.Fatalf("expected status 304, got %d", secondRec.Code)
	}
}

func TestWriteJSONWithETagRejectsNilInputs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)

	if err := WriteJSONWithETag(nil, req, map[string]string{"status": "ready"}); err == nil {
		t.Fatal("expected nil writer error")
	}
	if err := WriteJSONWithETag(httptest.NewRecorder(), nil, map[string]string{"status": "ready"}); err == nil {
		t.Fatal("expected nil request error")
	}
}

func TestWriteJSONWithETagPropagatesWriteErrors(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/42", nil)
	writer := &failingResponseWriter{header: make(http.Header)}

	err := WriteJSONWithETag(writer, req, map[string]string{"product_id": "product-42", "status": "ready"})
	if err == nil {
		t.Fatal("expected write error")
	}
	if !strings.Contains(err.Error(), "write response body") {
		t.Fatalf("expected wrapped write error, got %v", err)
	}
}