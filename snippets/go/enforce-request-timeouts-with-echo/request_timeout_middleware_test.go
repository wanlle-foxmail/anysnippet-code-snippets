package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func newTestServer(timeout time.Duration, handler echo.HandlerFunc) *echo.Echo {
	e := echo.New()
	e.Use(RequestTimeout(timeout))
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

func decodeMessage(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()

	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	return body["message"]
}

func TestRequestTimeoutReturnsBufferedResponseWhenHandlerFinishesInTime(t *testing.T) {
	server := newTestServer(50*time.Millisecond, func(c echo.Context) error {
		time.Sleep(5 * time.Millisecond)
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})

	rec := doRequest(t, server, "/work")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if decodeMessage(t, rec) != "ok" {
		t.Fatalf("expected message ok, got %q", decodeMessage(t, rec))
	}
}

func TestRequestTimeoutReturns503AndDiscardsLateWrites(t *testing.T) {
	server := newTestServer(10*time.Millisecond, func(c echo.Context) error {
		time.Sleep(30 * time.Millisecond)
		return c.JSON(http.StatusOK, map[string]string{"message": "late"})
	})

	rec := doRequest(t, server, "/work")
	time.Sleep(40 * time.Millisecond)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}
	if got := decodeMessage(t, rec); got != "request timed out" {
		t.Fatalf("expected timeout message, got %q", got)
	}
	if rec.Body.String() == "late" || rec.Body.String() == "{\"message\":\"late\"}\n" {
		t.Fatal("expected late handler writes to be discarded")
	}
}

func TestRequestTimeoutCancelsRequestContext(t *testing.T) {
	observedCancellation := make(chan struct{}, 1)
	server := newTestServer(10*time.Millisecond, func(c echo.Context) error {
		<-c.Request().Context().Done()
		observedCancellation <- struct{}{}
		return nil
	})

	rec := doRequest(t, server, "/work")

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}

	select {
	case <-observedCancellation:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected handler to observe request context cancellation")
	}
}

func TestRequestTimeoutAppliesPerRequestDeadlinesIndependently(t *testing.T) {
	server := newTestServer(20*time.Millisecond, func(c echo.Context) error {
		if c.QueryParam("speed") == "slow" {
			time.Sleep(40 * time.Millisecond)
			return c.JSON(http.StatusOK, map[string]string{"message": "slow"})
		}
		time.Sleep(5 * time.Millisecond)
		return c.JSON(http.StatusOK, map[string]string{"message": "fast"})
	})

	results := make(map[string]*httptest.ResponseRecorder)
	var mutex sync.Mutex
	var waitGroup sync.WaitGroup

	for _, speed := range []string{"fast", "slow"} {
		speed := speed
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			rec := doRequest(t, server, "/work?speed="+speed)
			mutex.Lock()
			results[speed] = rec
			mutex.Unlock()
		}()
	}

	waitGroup.Wait()

	if results["fast"].Code != http.StatusOK {
		t.Fatalf("expected fast request status 200, got %d", results["fast"].Code)
	}
	if results["slow"].Code != http.StatusServiceUnavailable {
		t.Fatalf("expected slow request status 503, got %d", results["slow"].Code)
	}
}

func TestRequestTimeoutRejectsNonPositiveTimeout(t *testing.T) {
	server := newTestServer(0, func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	rec := doRequest(t, server, "/work")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestRequestTimeoutPreservesHandlerErrors(t *testing.T) {
	server := newTestServer(50*time.Millisecond, func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	})

	rec := doRequest(t, server, "/work")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("expected error body from handler error")
	}
}
