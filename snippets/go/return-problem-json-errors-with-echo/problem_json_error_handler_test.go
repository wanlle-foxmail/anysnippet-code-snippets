package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func decodeProblemDetails(t *testing.T, recorder *httptest.ResponseRecorder) ProblemDetails {
	t.Helper()

	var problem ProblemDetails
	if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
		t.Fatalf("decode problem details: %v", err)
	}
	return problem
}

func TestProblemJSONErrorHandlerFormatsEchoHTTPError(t *testing.T) {
	server := newServer()
	req := httptest.NewRequest(http.MethodGet, "/orders/123", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)
	problem := decodeProblemDetails(t, rec)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
	if got := rec.Header().Get(echo.HeaderContentType); got != problemJSONContentType {
		t.Fatalf("expected %q, got %q", problemJSONContentType, got)
	}
	if problem.Type != "about:blank" {
		t.Fatalf("expected type about:blank, got %q", problem.Type)
	}
	if problem.Title != http.StatusText(http.StatusNotFound) {
		t.Fatalf("expected title %q, got %q", http.StatusText(http.StatusNotFound), problem.Title)
	}
	if problem.Detail != "order not found" {
		t.Fatalf("expected detail %q, got %q", "order not found", problem.Detail)
	}
	if problem.Instance != "/orders/123" {
		t.Fatalf("expected instance /orders/123, got %q", problem.Instance)
	}
}

func TestProblemJSONErrorHandlerFormatsGenericErrorsAs500(t *testing.T) {
	server := newServer()
	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)
	problem := decodeProblemDetails(t, rec)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
	if problem.Detail != http.StatusText(http.StatusInternalServerError) {
		t.Fatalf("expected detail %q, got %q", http.StatusText(http.StatusInternalServerError), problem.Detail)
	}
	if problem.Instance != "/boom" {
		t.Fatalf("expected instance /boom, got %q", problem.Instance)
	}
}

func TestProblemJSONErrorHandlerFallsBackToStatusTextForNonStringHTTPErrorMessages(t *testing.T) {
	server := newServer()
	req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)
	problem := decodeProblemDetails(t, rec)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if problem.Detail != http.StatusText(http.StatusBadRequest) {
		t.Fatalf("expected detail %q, got %q", http.StatusText(http.StatusBadRequest), problem.Detail)
	}
}

func TestProblemJSONErrorHandlerDoesNotRewriteCommittedResponses(t *testing.T) {
	server := newServer()
	req := httptest.NewRequest(http.MethodGet, "/accepted", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", rec.Code)
	}
	if got := rec.Header().Get(echo.HeaderContentType); got == problemJSONContentType {
		t.Fatalf("did not expect committed response content type %q", got)
	}
	if rec.Body.String() != "accepted\n" {
		t.Fatalf("expected original committed body, got %q", rec.Body.String())
	}
}

func TestWriteProblemJSONFillsDefaultTypeAndTitle(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/orders/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := WriteProblemJSON(c, ProblemDetails{Status: http.StatusConflict, Detail: "duplicate order"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	problem := decodeProblemDetails(t, rec)
	if problem.Type != "about:blank" {
		t.Fatalf("expected type about:blank, got %q", problem.Type)
	}
	if problem.Title != http.StatusText(http.StatusConflict) {
		t.Fatalf("expected title %q, got %q", http.StatusText(http.StatusConflict), problem.Title)
	}
	if problem.Detail != "duplicate order" {
		t.Fatalf("expected detail duplicate order, got %q", problem.Detail)
	}
}

func TestWriteProblemJSONRejectsNilContext(t *testing.T) {
	err := WriteProblemJSON(nil, ProblemDetails{Status: http.StatusBadRequest})
	if err == nil {
		t.Fatal("expected nil context error")
	}
}