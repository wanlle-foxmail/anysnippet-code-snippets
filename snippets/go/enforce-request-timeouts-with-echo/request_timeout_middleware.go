package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
)

type timeoutBufferWriter struct {
	recorder *httptest.ResponseRecorder
	timedOut atomic.Bool
}

func newTimeoutBufferWriter() *timeoutBufferWriter {
	return &timeoutBufferWriter{recorder: httptest.NewRecorder()}
}

func (writer *timeoutBufferWriter) Header() http.Header {
	return writer.recorder.Header()
}

func (writer *timeoutBufferWriter) Write(buffer []byte) (int, error) {
	if writer.timedOut.Load() {
		return len(buffer), nil
	}
	return writer.recorder.Write(buffer)
}

func (writer *timeoutBufferWriter) WriteHeader(statusCode int) {
	if writer.timedOut.Load() {
		return
	}
	writer.recorder.WriteHeader(statusCode)
}

func RequestTimeout(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if timeout <= 0 {
				return echo.NewHTTPError(http.StatusInternalServerError, "timeout must be greater than 0")
			}

			// Flow:
			//   add request deadline to context
			//      |
			//      +-> handler finishes in time -> flush buffered response
			//      `-> deadline fires first -> return 503 and discard late writes
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			response := c.Response()
			originalWriter := response.Writer
			timeoutWriter := newTimeoutBufferWriter()
			response.Writer = timeoutWriter
			c.SetRequest(c.Request().WithContext(ctx))

			resultCh := make(chan error, 1)
			go func() {
				resultCh <- next(c)
			}()

			select {
			case err := <-resultCh:
				response.Writer = originalWriter
				if err != nil {
					return err
				}
				if err := flushBufferedResponse(response, timeoutWriter); err != nil {
					return err
				}
				return nil
			case <-ctx.Done():
				timeoutWriter.timedOut.Store(true)
				return writeTimeoutResponse(originalWriter)
			}
		}
	}
}

func flushBufferedResponse(response *echo.Response, timeoutWriter *timeoutBufferWriter) error {
	for key, values := range timeoutWriter.recorder.Header() {
		response.Header()[key] = append([]string(nil), values...)
	}

	statusCode := timeoutWriter.recorder.Code
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	response.WriteHeader(statusCode)

	if timeoutWriter.recorder.Body.Len() == 0 {
		return nil
	}

	if _, err := response.Write(timeoutWriter.recorder.Body.Bytes()); err != nil {
		return fmt.Errorf("write buffered response: %w", err)
	}

	return nil
}

func writeTimeoutResponse(originalWriter http.ResponseWriter) error {
	payload := []byte("{\"message\":\"request timed out\"}\n")
	originalWriter.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	originalWriter.WriteHeader(http.StatusServiceUnavailable)
	_, err := originalWriter.Write(payload)
	if err != nil {
		return fmt.Errorf("write timeout response: %w", err)
	}
	return nil
}

func helloHandler(c echo.Context) error {
	select {
	case <-time.After(10 * time.Millisecond):
		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	case <-c.Request().Context().Done():
		return nil
	}
}

func newServer(timeout time.Duration) *echo.Echo {
	e := echo.New()
	e.Use(RequestTimeout(timeout))
	e.GET("/hello", helloHandler)
	return e
}

func main() {
	e := newServer(100 * time.Millisecond)
	e.Logger.Fatal(e.Start(":8080"))
}
