package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/labstack/echo/v4"
)

func doOrderRequest(t *testing.T, server *echo.Echo, method string, key string) (*httptest.ResponseRecorder, map[string]string) {
	t.Helper()

	req := httptest.NewRequest(method, "/orders", nil)
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	var body map[string]string
	if rec.Body.Len() > 0 && json.Valid(rec.Body.Bytes()) {
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode response body: %v", err)
		}
	}

	return rec, body
}

func TestIdempotencyMiddlewareRequiresKeyForPostRequests(t *testing.T) {
	server := newServer(NewMemoryIdempotencyStore())
	rec, _ := doOrderRequest(t, server, http.MethodPost, "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestIdempotencyMiddlewareCachesSuccessfulPostResponse(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	server := newServer(store)

	firstRec, firstBody := doOrderRequest(t, server, http.MethodPost, "key-1")
	secondRec, secondBody := doOrderRequest(t, server, http.MethodPost, "key-1")

	if firstRec.Code != http.StatusCreated {
		t.Fatalf("expected first status 201, got %d", firstRec.Code)
	}
	if secondRec.Code != http.StatusCreated {
		t.Fatalf("expected replayed status 201, got %d", secondRec.Code)
	}
	if secondRec.Header().Get("X-Order-ID") != "order-123" {
		t.Fatalf("expected cached X-Order-ID header, got %q", secondRec.Header().Get("X-Order-ID"))
	}
	if firstBody["order_id"] != secondBody["order_id"] {
		t.Fatalf("expected cached order_id %q, got %q", firstBody["order_id"], secondBody["order_id"])
	}
}

func TestIdempotencyMiddlewareSeparatesDifferentKeys(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	callCount := atomic.Int32{}
	e := echo.New()
	e.Use(IdempotencyMiddleware(store))
	e.POST("/orders", func(c echo.Context) error {
		count := callCount.Add(1)
		return c.JSON(http.StatusCreated, map[string]int32{"call": count})
	})

	firstReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	firstReq.Header.Set("Idempotency-Key", "key-1")
	firstRec := httptest.NewRecorder()
	e.ServeHTTP(firstRec, firstReq)

	secondReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	secondReq.Header.Set("Idempotency-Key", "key-2")
	secondRec := httptest.NewRecorder()
	e.ServeHTTP(secondRec, secondReq)

	if callCount.Load() != 2 {
		t.Fatalf("expected two handler executions, got %d", callCount.Load())
	}
}

func TestIdempotencyMiddlewareReturnsConflictWhileRequestIsInProgress(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	started := make(chan struct{})
	unblock := make(chan struct{})

	e := echo.New()
	e.Use(IdempotencyMiddleware(store))
	e.POST("/orders", func(c echo.Context) error {
		close(started)
		<-unblock
		return c.JSON(http.StatusCreated, map[string]string{"status": "created"})
	})

	firstReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	firstReq.Header.Set("Idempotency-Key", "same-key")
	firstRec := httptest.NewRecorder()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		e.ServeHTTP(firstRec, firstReq)
	}()

	<-started
	secondRec, _ := doOrderRequest(t, e, http.MethodPost, "same-key")
	if secondRec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", secondRec.Code)
	}

	close(unblock)
	waitGroup.Wait()
}

func TestIdempotencyMiddlewareDoesNotCacheFailedResponses(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	callCount := atomic.Int32{}
	e := echo.New()
	e.Use(IdempotencyMiddleware(store))
	e.POST("/orders", func(c echo.Context) error {
		count := callCount.Add(1)
		if count == 1 {
			return c.JSON(http.StatusInternalServerError, map[string]string{"status": "failed"})
		}
		return c.JSON(http.StatusCreated, map[string]string{"status": "created"})
	})

	firstRec, _ := doOrderRequest(t, e, http.MethodPost, "retry-key")
	secondRec, secondBody := doOrderRequest(t, e, http.MethodPost, "retry-key")

	if firstRec.Code != http.StatusInternalServerError {
		t.Fatalf("expected first status 500, got %d", firstRec.Code)
	}
	if secondRec.Code != http.StatusCreated {
		t.Fatalf("expected second status 201, got %d", secondRec.Code)
	}
	if secondBody["status"] != "created" {
		t.Fatalf("expected created status, got %q", secondBody["status"])
	}
	if callCount.Load() != 2 {
		t.Fatalf("expected two handler executions, got %d", callCount.Load())
	}
}

func TestIdempotencyMiddlewareBypassesNonPostRequests(t *testing.T) {
	server := newServer(NewMemoryIdempotencyStore())
	rec, _ := doOrderRequest(t, server, http.MethodGet, "")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", rec.Body.String())
	}
	if rec.Header().Get("Idempotency-Key") != "" {
		t.Fatal("expected middleware to ignore idempotency for GET requests")
	}
}

type loadAfterTryStartStore struct {
	loadCount int32
}

func (store *loadAfterTryStartStore) Load(key string) (CachedResponse, bool) {
	callNumber := atomic.AddInt32(&store.loadCount, 1)
	if callNumber == 2 {
		headers := make(http.Header)
		headers.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		headers.Set("X-Order-ID", "order-123")

		return CachedResponse{
			StatusCode: http.StatusCreated,
			Headers:    headers,
			Body:       []byte("{\"order_id\":\"order-123\",\"status\":\"created\"}\n"),
		}, true
	}

	return CachedResponse{}, false
}

func (store *loadAfterTryStartStore) TryStart(key string) bool {
	return false
}

func (store *loadAfterTryStartStore) Save(key string, response CachedResponse) {
}

func (store *loadAfterTryStartStore) Delete(key string) {
}

func TestIdempotencyMiddlewareReplaysCompletedResponseWhenStateChangesBeforeTryStart(t *testing.T) {
	server := newServer(&loadAfterTryStartStore{})
	rec, body := doOrderRequest(t, server, http.MethodPost, "key-1")

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	if rec.Header().Get("X-Order-ID") != "order-123" {
		t.Fatalf("expected cached X-Order-ID header, got %q", rec.Header().Get("X-Order-ID"))
	}
	if body["order_id"] != "order-123" {
		t.Fatalf("expected cached order_id, got %q", body["order_id"])
	}
}

func TestIdempotencyMiddlewareSeparatesDifferentUsersWithSameKey(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	callCount := atomic.Int32{}
	e := echo.New()
	e.Use(IdempotencyMiddleware(store))
	e.POST("/orders", func(c echo.Context) error {
		count := callCount.Add(1)
		return c.JSON(http.StatusCreated, map[string]string{"call": strconv.FormatInt(int64(count), 10)})
	})

	firstReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	firstReq.Header.Set("Idempotency-Key", "same-key")
	firstReq.Header.Set("X-User-ID", "user-a")
	firstRec := httptest.NewRecorder()
	e.ServeHTTP(firstRec, firstReq)

	secondReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	secondReq.Header.Set("Idempotency-Key", "same-key")
	secondReq.Header.Set("X-User-ID", "user-b")
	secondRec := httptest.NewRecorder()
	e.ServeHTTP(secondRec, secondReq)

	if firstRec.Code != http.StatusCreated || secondRec.Code != http.StatusCreated {
		t.Fatalf("expected both requests to succeed, got %d and %d", firstRec.Code, secondRec.Code)
	}
	if callCount.Load() != 2 {
		t.Fatalf("expected two handler executions, got %d", callCount.Load())
	}
}

func TestIdempotencyMiddlewareNormalizesAuthorizationScope(t *testing.T) {
	store := NewMemoryIdempotencyStore()
	callCount := atomic.Int32{}
	e := echo.New()
	e.Use(IdempotencyMiddleware(store))
	e.POST("/orders", func(c echo.Context) error {
		count := callCount.Add(1)
		return c.JSON(http.StatusCreated, map[string]string{"call": strconv.FormatInt(int64(count), 10)})
	})

	firstReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	firstReq.Header.Set("Idempotency-Key", "same-key")
	firstReq.Header.Set(echo.HeaderAuthorization, "Bearer shared-token")
	firstRec := httptest.NewRecorder()
	e.ServeHTTP(firstRec, firstReq)

	secondReq := httptest.NewRequest(http.MethodPost, "/orders", nil)
	secondReq.Header.Set("Idempotency-Key", "same-key")
	secondReq.Header.Set(echo.HeaderAuthorization, "bearer shared-token")
	secondRec := httptest.NewRecorder()
	e.ServeHTTP(secondRec, secondReq)

	if firstRec.Code != http.StatusCreated || secondRec.Code != http.StatusCreated {
		t.Fatalf("expected both requests to succeed, got %d and %d", firstRec.Code, secondRec.Code)
	}
	if callCount.Load() != 1 {
		t.Fatalf("expected one handler execution, got %d", callCount.Load())
	}
}
