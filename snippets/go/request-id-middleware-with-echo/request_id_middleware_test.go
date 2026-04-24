package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func doRequest(t *testing.T, requestID string) (*httptest.ResponseRecorder, map[string]string) {
	t.Helper()
	e := newServer()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	if requestID != "" {
		req.Header.Set(echo.HeaderXRequestID, requestID)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return rec, body
}

func TestUsesExistingRequestIDHeader(t *testing.T) {
	rec, body := doRequest(t, "incoming-id")
	if got := rec.Header().Get(echo.HeaderXRequestID); got != "incoming-id" {
		t.Fatalf("expected response request id incoming-id, got %q", got)
	}
	if got := body["request_id"]; got != "incoming-id" {
		t.Fatalf("expected body request_id incoming-id, got %q", got)
	}
}

func TestGeneratesRequestIDWhenHeaderIsMissing(t *testing.T) {
	rec, body := doRequest(t, "")
	requestID := rec.Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		t.Fatal("expected generated request id")
	}
	if len(requestID) != 32 {
		t.Fatalf("expected 32-character request id, got %q", requestID)
	}
	if _, err := hex.DecodeString(requestID); err != nil {
		t.Fatalf("expected hex request id, got %q", requestID)
	}
	if body["request_id"] != requestID {
		t.Fatalf("expected body request_id %q, got %q", requestID, body["request_id"])
	}
}

func TestTreatsEmptyHeaderAsMissing(t *testing.T) {
	e := newServer()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header[echo.HeaderXRequestID] = []string{""}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get(echo.HeaderXRequestID) == "" {
		t.Fatal("expected middleware to replace empty request id")
	}
}

func TestTreatsWhitespaceOnlyHeaderAsMissing(t *testing.T) {
	e := newServer()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set(echo.HeaderXRequestID, "   ")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	requestID := rec.Header().Get(echo.HeaderXRequestID)
	if requestID == "" {
		t.Fatal("expected middleware to replace whitespace-only request id")
	}
	if requestID == "   " {
		t.Fatal("expected middleware to trim and replace whitespace-only request id")
	}
}

func TestResponseHeaderAlwaysContainsRequestID(t *testing.T) {
	rec, _ := doRequest(t, "")
	if rec.Header().Get(echo.HeaderXRequestID) == "" {
		t.Fatal("expected response header to include request id")
	}
}

func TestHandlerReceivesSameRequestIDAsResponseHeader(t *testing.T) {
	rec, body := doRequest(t, "")
	if rec.Header().Get(echo.HeaderXRequestID) != body["request_id"] {
		t.Fatalf("expected body request_id %q, got %q", rec.Header().Get(echo.HeaderXRequestID), body["request_id"])
	}
}

func TestSeparateRequestsReceiveDifferentGeneratedIDs(t *testing.T) {
	firstRec, _ := doRequest(t, "")
	secondRec, _ := doRequest(t, "")
	firstID := firstRec.Header().Get(echo.HeaderXRequestID)
	secondID := secondRec.Header().Get(echo.HeaderXRequestID)
	if firstID == secondID {
		t.Fatalf("expected different generated request ids, got %q and %q", firstID, secondID)
	}
}

func TestReturnsInternalServerErrorWhenIDGenerationFails(t *testing.T) {
	originalReadRandom := readRandom
	readRandom = func(_ []byte) (int, error) {
		return 0, errors.New("random failed")
	}
	defer func() {
		readRandom = originalReadRandom
	}()

	e := newServer()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
	if rec.Header().Get(echo.HeaderXRequestID) != "" {
		t.Fatal("expected no response request id when generation fails")
	}
}
