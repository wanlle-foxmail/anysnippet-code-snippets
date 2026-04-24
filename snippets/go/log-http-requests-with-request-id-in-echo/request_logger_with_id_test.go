package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func newBufferLogger(buffer *bytes.Buffer) *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return attr
		},
	}
	return slog.New(slog.NewJSONHandler(buffer, handlerOptions))
}

func newLoggingServer(logger *slog.Logger, handler echo.HandlerFunc) *echo.Echo {
	e := echo.New()
	e.Use(RequestLoggerWithRequestID(logger))
	e.GET("/work", handler)
	return e
}

func doRequest(t *testing.T, server *echo.Echo, target string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	return rec
}

func decodeSingleLogLine(t *testing.T, buffer *bytes.Buffer) map[string]any {
	t.Helper()

	var logLine map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buffer.Bytes()), &logLine); err != nil {
		t.Fatalf("decode log line: %v", err)
	}
	return logLine
}

func TestRequestLoggerWithRequestIDLogsSuccessfulRequest(t *testing.T) {
	var buffer bytes.Buffer
	server := newLoggingServer(newBufferLogger(&buffer), func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/work", nil)
	req.Header.Set(echo.HeaderXRequestID, "incoming-id")
	req.Header.Set(echo.HeaderXRealIP, "203.0.113.10")
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	logLine := decodeSingleLogLine(t, &buffer)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if logLine["request_id"] != "incoming-id" {
		t.Fatalf("expected request_id incoming-id, got %#v", logLine["request_id"])
	}
	if logLine["method"] != http.MethodGet {
		t.Fatalf("expected method GET, got %#v", logLine["method"])
	}
	if logLine["path"] != "/work" {
		t.Fatalf("expected path /work, got %#v", logLine["path"])
	}
	if int(logLine["status"].(float64)) != http.StatusOK {
		t.Fatalf("expected status 200 in log, got %#v", logLine["status"])
	}
	if int(logLine["duration_ms"].(float64)) < 0 {
		t.Fatalf("expected non-negative duration, got %#v", logLine["duration_ms"])
	}
	if logLine["client_ip"] != "203.0.113.10" {
		t.Fatalf("expected client_ip 203.0.113.10, got %#v", logLine["client_ip"])
	}
}

func TestRequestLoggerWithRequestIDPrefersResponseHeaderRequestID(t *testing.T) {
	var buffer bytes.Buffer
	server := newLoggingServer(newBufferLogger(&buffer), func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderXRequestID, "generated-id")
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/work", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	logLine := decodeSingleLogLine(t, &buffer)

	if logLine["request_id"] != "generated-id" {
		t.Fatalf("expected request_id generated-id, got %#v", logLine["request_id"])
	}
	if int(logLine["status"].(float64)) != http.StatusNoContent {
		t.Fatalf("expected status 204 in log, got %#v", logLine["status"])
	}
}

func TestRequestLoggerWithRequestIDLogsEmptyRequestIDWhenMissing(t *testing.T) {
	var buffer bytes.Buffer
	server := newLoggingServer(newBufferLogger(&buffer), func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	rec := doRequest(t, server, "/work")
	logLine := decodeSingleLogLine(t, &buffer)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if logLine["request_id"] != "" {
		t.Fatalf("expected empty request_id, got %#v", logLine["request_id"])
	}
}

func TestRequestLoggerWithRequestIDLogsHTTPErrorStatus(t *testing.T) {
	var buffer bytes.Buffer
	server := newLoggingServer(newBufferLogger(&buffer), func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	})

	rec := doRequest(t, server, "/work")
	logLine := decodeSingleLogLine(t, &buffer)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if int(logLine["status"].(float64)) != http.StatusBadRequest {
		t.Fatalf("expected logged status 400, got %#v", logLine["status"])
	}
}

func TestRequestLoggerWithRequestIDLogsGenericErrorAsInternalServerError(t *testing.T) {
	var buffer bytes.Buffer
	server := newLoggingServer(newBufferLogger(&buffer), func(c echo.Context) error {
		return errors.New("boom")
	})

	rec := doRequest(t, server, "/work")
	logLine := decodeSingleLogLine(t, &buffer)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
	if int(logLine["status"].(float64)) != http.StatusInternalServerError {
		t.Fatalf("expected logged status 500, got %#v", logLine["status"])
	}
}
