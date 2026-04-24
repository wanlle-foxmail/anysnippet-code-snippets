package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var timeNow = time.Now

type clientBucket struct {
	Tokens     int
	LastRefill time.Time
}

func TokenBucketRateLimit(capacity int, refillInterval time.Duration) echo.MiddlewareFunc {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}
	if refillInterval <= 0 {
		panic("refill interval must be greater than 0")
	}

	var mutex sync.Mutex
	buckets := make(map[string]clientBucket)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Flow:
			//   read the client IP bucket
			//      |
			//      +-> token available -> consume it -> call next
			//      `-> empty bucket -> set Retry-After -> return 429
			allowed, retryAfter := takeClientToken(c.RealIP(), capacity, refillInterval, buckets, &mutex)
			if !allowed {
				c.Response().Header().Set("Retry-After", strconv.Itoa(retryAfterSeconds(retryAfter)))
				return echo.NewHTTPError(http.StatusTooManyRequests, "rate limit exceeded")
			}

			return next(c)
		}
	}
}

func takeClientToken(clientIP string, capacity int, refillInterval time.Duration, buckets map[string]clientBucket, mutex *sync.Mutex) (bool, time.Duration) {
	if strings.TrimSpace(clientIP) == "" {
		clientIP = "unknown"
	}

	now := timeNow()

	mutex.Lock()
	defer mutex.Unlock()

	bucket, exists := buckets[clientIP]
	if !exists {
		buckets[clientIP] = clientBucket{Tokens: capacity - 1, LastRefill: now}
		return true, 0
	}

	bucket = refillTokens(bucket, now, capacity, refillInterval)
	if bucket.Tokens == 0 {
		buckets[clientIP] = bucket
		nextRefillAt := bucket.LastRefill.Add(refillInterval)
		return false, nextRefillAt.Sub(now)
	}

	bucket.Tokens--
	buckets[clientIP] = bucket
	return true, 0
}

func refillTokens(bucket clientBucket, now time.Time, capacity int, refillInterval time.Duration) clientBucket {
	if bucket.LastRefill.IsZero() {
		bucket.LastRefill = now
		if bucket.Tokens == 0 {
			bucket.Tokens = capacity
		}
		return bucket
	}

	elapsed := now.Sub(bucket.LastRefill)
	if elapsed < refillInterval {
		return bucket
	}

	refillCount := int(elapsed / refillInterval)
	bucket.Tokens = minInt(capacity, bucket.Tokens+refillCount)
	bucket.LastRefill = bucket.LastRefill.Add(time.Duration(refillCount) * refillInterval)
	return bucket
}

func retryAfterSeconds(retryAfter time.Duration) int {
	if retryAfter <= 0 {
		return 1
	}
	seconds := int((retryAfter + time.Second - 1) / time.Second)
	if seconds < 1 {
		return 1
	}
	return seconds
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

func newServer(capacity int, refillInterval time.Duration) *echo.Echo {
	e := echo.New()
	e.Use(TokenBucketRateLimit(capacity, refillInterval))
	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})
	return e
}

func main() {
	e := newServer(2, time.Second)
	e.Logger.Fatal(e.Start(":8080"))
}