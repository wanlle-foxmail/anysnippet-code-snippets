package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const RequestIDContextKey = "request_id"

var readRandom = rand.Read

func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Flow:
			//   read X-Request-Id header
			//      |
			//      +-> present -> use it
			//      `-> missing -> generate one or return error
			//   store it in context and response header -> call next
			requestID := strings.TrimSpace(c.Request().Header.Get(echo.HeaderXRequestID))
			if requestID == "" {
				generatedRequestID, err := generateRequestID()
				if err != nil {
					return fmt.Errorf("generate request id: %w", err)
				}
				requestID = generatedRequestID
			}

			c.Set(RequestIDContextKey, requestID)
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)

			return next(c)
		}
	}
}

func generateRequestID() (string, error) {
	buffer := make([]byte, 16)
	if _, err := readRandom(buffer); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return hex.EncodeToString(buffer), nil
}

func helloHandler(c echo.Context) error {
	requestID, ok := c.Get(RequestIDContextKey).(string)
	if !ok || requestID == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "missing request id")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":    "ok",
		"request_id": requestID,
	})
}

func newServer() *echo.Echo {
	e := echo.New()
	e.Use(RequestIDMiddleware())
	e.GET("/hello", helloHandler)
	return e
}

func main() {
	e := newServer()
	e.Logger.Fatal(e.Start(":8080"))
}
