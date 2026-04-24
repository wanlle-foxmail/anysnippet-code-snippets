package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"crypto/sha256"
	"encoding/hex"
	"github.com/labstack/echo/v4"
)

type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type IdempotencyStore interface {
	Load(key string) (CachedResponse, bool)
	TryStart(key string) bool
	Save(key string, response CachedResponse)
	Delete(key string)
}

type memoryEntry struct {
	response   CachedResponse
	inProgress bool
	completed  bool
}

type MemoryIdempotencyStore struct {
	mutex   sync.Mutex
	entries map[string]memoryEntry
}

func NewMemoryIdempotencyStore() *MemoryIdempotencyStore {
	return &MemoryIdempotencyStore{entries: make(map[string]memoryEntry)}
}

func (store *MemoryIdempotencyStore) Load(key string) (CachedResponse, bool) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	entry, ok := store.entries[key]
	if !ok || !entry.completed {
		return CachedResponse{}, false
	}

	return cloneCachedResponse(entry.response), true
}

func (store *MemoryIdempotencyStore) TryStart(key string) bool {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	entry, exists := store.entries[key]
	if exists && (entry.inProgress || entry.completed) {
		return false
	}

	store.entries[key] = memoryEntry{inProgress: true}
	return true
}

func (store *MemoryIdempotencyStore) Save(key string, response CachedResponse) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.entries[key] = memoryEntry{
		response:  cloneCachedResponse(response),
		completed: true,
	}
}

func (store *MemoryIdempotencyStore) Delete(key string) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	delete(store.entries, key)
}

func IdempotencyMiddleware(store IdempotencyStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if store == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "idempotency store is required")
			}
			if c.Request().Method != http.MethodPost {
				return next(c)
			}

			idempotencyKey := strings.TrimSpace(c.Request().Header.Get("Idempotency-Key"))
			if idempotencyKey == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Idempotency-Key header is required")
			}
			scopedKey := scopedIdempotencyKey(c, idempotencyKey)

			if cachedResponse, ok := store.Load(scopedKey); ok {
				return writeCachedResponse(c.Response().Writer, cachedResponse)
			}
			if !store.TryStart(scopedKey) {
				if cachedResponse, ok := store.Load(scopedKey); ok {
					return writeCachedResponse(c.Response().Writer, cachedResponse)
				}
				return echo.NewHTTPError(http.StatusConflict, "request already in progress")
			}

			// Flow:
			//   require Idempotency-Key on POST
			//      |
			//      +-> completed key -> replay cached response
			//      +-> in-progress key -> return 409
			//      `-> new key -> run handler -> cache successful response -> return it
			response := c.Response()
			originalWriter := response.Writer
			recorder := httptest.NewRecorder()
			response.Writer = recorder

			err := next(c)
			response.Writer = originalWriter
			if err != nil {
				store.Delete(scopedKey)
				return err
			}

			capturedResponse := cachedResponseFromRecorder(recorder)
			if capturedResponse.StatusCode >= http.StatusOK && capturedResponse.StatusCode < http.StatusMultipleChoices {
				store.Save(scopedKey, capturedResponse)
			} else {
				store.Delete(scopedKey)
			}

			if err := writeCachedResponse(originalWriter, capturedResponse); err != nil {
				return fmt.Errorf("write captured response: %w", err)
			}
			return nil
		}
	}
}

func scopedIdempotencyKey(c echo.Context, idempotencyKey string) string {
	requestPath := c.Request().URL.EscapedPath()
	if requestPath == "" {
		requestPath = "/"
	}

	callerScope := requestCallerScope(c)
	return c.Request().Method + ":" + requestPath + ":" + callerScope + ":" + idempotencyKey
}

func requestCallerScope(c echo.Context) string {
	userID := strings.TrimSpace(c.Request().Header.Get("X-User-ID"))
	if userID != "" {
		return "user:" + userID
	}

	authorizationHeader := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
	if authorizationHeader != "" {
		return "auth:" + authorizationScopeHash(authorizationHeader)
	}

	return "anonymous"
}

func authorizationScopeHash(authorizationHeader string) string {
	parts := strings.Fields(authorizationHeader)
	normalizedValue := authorizationHeader
	if len(parts) > 0 {
		normalizedValue = strings.ToLower(parts[0])
		if len(parts) > 1 {
			normalizedValue += " " + strings.Join(parts[1:], " ")
		}
	}

	hash := sha256.Sum256([]byte(normalizedValue))
	return hex.EncodeToString(hash[:16])
}
func cachedResponseFromRecorder(recorder *httptest.ResponseRecorder) CachedResponse {
	statusCode := recorder.Code
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	headers := make(http.Header, len(recorder.Header()))
	for key, values := range recorder.Header() {
		headers[key] = append([]string(nil), values...)
	}

	body := append([]byte(nil), recorder.Body.Bytes()...)
	return CachedResponse{StatusCode: statusCode, Headers: headers, Body: body}
}

func cloneCachedResponse(response CachedResponse) CachedResponse {
	headers := make(http.Header, len(response.Headers))
	for key, values := range response.Headers {
		headers[key] = append([]string(nil), values...)
	}
	body := append([]byte(nil), response.Body...)
	return CachedResponse{StatusCode: response.StatusCode, Headers: headers, Body: body}
}

func writeCachedResponse(writer http.ResponseWriter, response CachedResponse) error {
	if writer == nil {
		return errors.New("writer is required")
	}
	for key, values := range response.Headers {
		writer.Header()[key] = append([]string(nil), values...)
	}
	writer.WriteHeader(response.StatusCode)
	if len(response.Body) == 0 {
		return nil
	}
	_, err := writer.Write(response.Body)
	return err
}

func newServer(store IdempotencyStore) *echo.Echo {
	e := echo.New()
	e.Use(IdempotencyMiddleware(store))
	e.POST("/orders", func(c echo.Context) error {
		c.Response().Header().Set("X-Order-ID", "order-123")
		return c.JSON(http.StatusCreated, map[string]string{
			"order_id": "order-123",
			"status":   "created",
		})
	})
	e.GET("/orders", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	return e
}

func main() {
	store := NewMemoryIdempotencyStore()
	e := newServer(store)
	e.Logger.Fatal(e.Start(":8080"))
}
