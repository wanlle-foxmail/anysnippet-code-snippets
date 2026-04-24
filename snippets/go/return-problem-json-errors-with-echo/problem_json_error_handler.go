package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const problemJSONContentType = "application/problem+json"

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func WriteProblemJSON(c echo.Context, problem ProblemDetails) error {
	if c == nil {
		return errors.New("context is required")
	}
	if problem.Status == 0 {
		problem.Status = http.StatusInternalServerError
	}
	if problem.Type == "" {
		problem.Type = "about:blank"
	}
	if problem.Title == "" {
		problem.Title = http.StatusText(problem.Status)
	}

	payload, err := json.Marshal(problem)
	if err != nil {
		return fmt.Errorf("marshal problem details: %w", err)
	}

	response := c.Response()
	response.Header().Set(echo.HeaderContentType, problemJSONContentType)
	response.WriteHeader(problem.Status)
	if _, err := response.Write(append(payload, '\n')); err != nil {
		return fmt.Errorf("write problem details: %w", err)
	}
	return nil

}

func ProblemJSONErrorHandler(err error, c echo.Context) {
	if c == nil || c.Response().Committed {
		return
	}

	problem := problemFromError(c, err)
	if writeErr := WriteProblemJSON(c, problem); writeErr != nil {
		c.Echo().Logger.Error(writeErr)
	}
}

func problemFromError(c echo.Context, err error) ProblemDetails {
	statusCode := http.StatusInternalServerError
	detail := http.StatusText(http.StatusInternalServerError)

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		statusCode = httpError.Code
		detail = problemDetailFromHTTPError(httpError)
	}

	instance := ""
	if c != nil && c.Request() != nil && c.Request().URL != nil {
		instance = c.Request().URL.Path
	}

	return ProblemDetails{
		Type:     "about:blank",
		Title:    http.StatusText(statusCode),
		Status:   statusCode,
		Detail:   detail,
		Instance: instance,
	}
}

func problemDetailFromHTTPError(httpError *echo.HTTPError) string {
	if httpError == nil {
		return http.StatusText(http.StatusInternalServerError)
	}

	switch message := httpError.Message.(type) {
	case string:
		trimmedMessage := strings.TrimSpace(message)
		if trimmedMessage != "" {
			return trimmedMessage
		}
	case error:
		trimmedMessage := strings.TrimSpace(message.Error())
		if trimmedMessage != "" {
			return trimmedMessage
		}
	}

	return http.StatusText(httpError.Code)

}

func newServer() *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = ProblemJSONErrorHandler
	e.GET("/orders/:orderID", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusNotFound, "order not found")
	})
	e.GET("/boom", func(c echo.Context) error {
		return errors.New("database unavailable")
	})
	e.GET("/invalid", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"field": "name"})
	})
	e.GET("/accepted", func(c echo.Context) error {
		c.Response().WriteHeader(http.StatusAccepted)
		_, _ = c.Response().Write([]byte("accepted\n"))
		return echo.NewHTTPError(http.StatusBadRequest, "ignored")
	})
	return e
}

func main() {
	e := newServer()
	e.Logger.Fatal(e.Start(":8080"))
}