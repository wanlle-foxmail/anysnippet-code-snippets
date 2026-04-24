package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func doRateLimitedRequest(t *testing.T, server http.Handler, clientIP string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	if clientIP != "" {
		req.Header.Set("X-Forwarded-For", clientIP)
	}
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	return rec
}

func TestTokenBucketRateLimitAllowsRequestsWithinCapacity(t *testing.T) {
	originalTimeNow := timeNow
	timeNow = func() time.Time { return time.Unix(1_700_000_000, 0) }
	defer func() {
		timeNow = originalTimeNow
	}()

	server := newServer(2, time.Second)
	firstRec := doRateLimitedRequest(t, server, "203.0.113.10")
	secondRec := doRateLimitedRequest(t, server, "203.0.113.10")

	if firstRec.Code != http.StatusOK {
		t.Fatalf("expected first status 200, got %d", firstRec.Code)
	}
	if secondRec.Code != http.StatusOK {
		t.Fatalf("expected second status 200, got %d", secondRec.Code)
	}
}

func TestTokenBucketRateLimitReturns429WhenBucketIsEmpty(t *testing.T) {
	originalTimeNow := timeNow
	now := time.Unix(1_700_000_000, 0)
	timeNow = func() time.Time { return now }
	defer func() {
		timeNow = originalTimeNow
	}()

	server := newServer(1, 2*time.Second)
	firstRec := doRateLimitedRequest(t, server, "203.0.113.10")
	secondRec := doRateLimitedRequest(t, server, "203.0.113.10")

	if firstRec.Code != http.StatusOK {
		t.Fatalf("expected first status 200, got %d", firstRec.Code)
	}
	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second status 429, got %d", secondRec.Code)
	}
	if secondRec.Header().Get("Retry-After") != "2" {
		t.Fatalf("expected Retry-After 2, got %q", secondRec.Header().Get("Retry-After"))
	}
}

func TestTokenBucketRateLimitRefillsTokensAfterTheInterval(t *testing.T) {
	originalTimeNow := timeNow
	currentTime := time.Unix(1_700_000_000, 0)
	timeNow = func() time.Time { return currentTime }
	defer func() {
		timeNow = originalTimeNow
	}()

	server := newServer(1, time.Second)
	firstRec := doRateLimitedRequest(t, server, "203.0.113.10")
	blockedRec := doRateLimitedRequest(t, server, "203.0.113.10")
	currentTime = currentTime.Add(time.Second)
	refilledRec := doRateLimitedRequest(t, server, "203.0.113.10")

	if firstRec.Code != http.StatusOK {
		t.Fatalf("expected first status 200, got %d", firstRec.Code)
	}
	if blockedRec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected blocked status 429, got %d", blockedRec.Code)
	}
	if refilledRec.Code != http.StatusOK {
		t.Fatalf("expected refilled status 200, got %d", refilledRec.Code)
	}
}

func TestTokenBucketRateLimitSeparatesBucketsPerClientIP(t *testing.T) {
	originalTimeNow := timeNow
	timeNow = func() time.Time { return time.Unix(1_700_000_000, 0) }
	defer func() {
		timeNow = originalTimeNow
	}()

	server := newServer(1, time.Second)
	firstClientRec := doRateLimitedRequest(t, server, "203.0.113.10")
	secondClientRec := doRateLimitedRequest(t, server, "203.0.113.11")

	if firstClientRec.Code != http.StatusOK {
		t.Fatalf("expected first client status 200, got %d", firstClientRec.Code)
	}
	if secondClientRec.Code != http.StatusOK {
		t.Fatalf("expected second client status 200, got %d", secondClientRec.Code)
	}
}

func TestTokenBucketRateLimitUsesUnknownBucketWhenRealIPIsMissing(t *testing.T) {
	originalTimeNow := timeNow
	timeNow = func() time.Time { return time.Unix(1_700_000_000, 0) }
	defer func() {
		timeNow = originalTimeNow
	}()

	server := newServer(1, time.Second)
	firstRec := doRateLimitedRequest(t, server, "")
	secondRec := doRateLimitedRequest(t, server, "")

	if firstRec.Code != http.StatusOK {
		t.Fatalf("expected first status 200, got %d", firstRec.Code)
	}
	if secondRec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second status 429, got %d", secondRec.Code)
	}
}

func TestTokenBucketRateLimitPanicsOnInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		capacity       int
		refillInterval time.Duration
	}{
		{name: "invalid capacity", capacity: 0, refillInterval: time.Second},
		{name: "invalid refill interval", capacity: 1, refillInterval: 0},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			defer func() {
				if recoveredValue := recover(); recoveredValue == nil {
					t.Fatal("expected panic for invalid configuration")
				}
			}()

			_ = newServer(testCase.capacity, testCase.refillInterval)
		})
	}
}