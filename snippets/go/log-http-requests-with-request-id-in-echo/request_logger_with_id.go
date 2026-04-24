package main

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func RequestLoggerWithRequestID(logger *slog.Logger) echo.MiddlewareFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startedAt := time.Now()
			err := next(c)

			requestID := requestIDFromHeaders(c)
			statusCode := responseStatusCode(c, err)
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}

			// Flow: run handler -> derive request metadata -> emit one structured access log -> return original error
			logger.Info(
				"request completed",
				slog.String("request_id", requestID),
				slog.String("method", c.Request().Method),
				slog.String("path", path),
				slog.Int("status", statusCode),
				slog.Int64("duration_ms", time.Since(startedAt).Milliseconds()),
				slog.String("client_ip", c.RealIP()),
			)

			return err
		}
	}
}

func requestIDFromHeaders(c echo.Context) string {
	requestID := strings.TrimSpace(c.Response().Header().Get(echo.HeaderXRequestID))
	if requestID != "" {
		return requestID
	}
	return strings.TrimSpace(c.Request().Header.Get(echo.HeaderXRequestID))
}

func responseStatusCode(c echo.Context, err error) int {
	if err != nil {
		if c.Response().Committed && c.Response().Status != 0 {
			return c.Response().Status
		}

		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			return httpError.Code
		}

		return http.StatusInternalServerError
	}

	if c.Response().Status != 0 {
		return c.Response().Status
	}

	return http.StatusOK
}

func newServer(logger *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Use(RequestLoggerWithRequestID(logger))
	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	})
	return e
}

func main() {
	logger := slog.Default()
	e := newServer(logger)
	e.Logger.Fatal(e.Start(":8080"))
}
