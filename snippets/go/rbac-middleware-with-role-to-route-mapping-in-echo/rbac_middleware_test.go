package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRBACMiddlewareAllowsMappedRole(t *testing.T) {
	response := performRoleRequest(http.MethodGet, "/admin", "admin")
	assertJSONMessage(t, response, http.StatusOK, "admin ok")
}

func TestRBACMiddlewareRejectsMissingRole(t *testing.T) {
	response := performRoleRequest(http.MethodGet, "/admin", "")
	assertJSONMessage(t, response, http.StatusUnauthorized, "role is required")
}

func TestRBACMiddlewareRejectsForbiddenRole(t *testing.T) {
	response := performRoleRequest(http.MethodGet, "/admin", "viewer")
	assertJSONMessage(t, response, http.StatusForbidden, "role is not allowed")
}

func TestRBACMiddlewareAllowsAnyMappedRoleForSharedRoute(t *testing.T) {
	response := performRoleRequest(http.MethodGet, "/reports", "analyst")
	assertJSONMessage(t, response, http.StatusOK, "reports ok")
}

func TestRBACMiddlewareRejectsUnmappedRoute(t *testing.T) {
	e := echoTestServerWithUnmappedRoute()
	request := httptest.NewRequest(http.MethodGet, "/billing", nil)
	request.Header.Set(RoleHeader, "admin")
	response := httptest.NewRecorder()

	e.ServeHTTP(response, request)
	assertJSONMessage(t, response, http.StatusForbidden, "route is not allowed")
}

func TestRBACMiddlewareNormalizesRoleHeaderWhitespaceAndCase(t *testing.T) {
	response := performRoleRequest(http.MethodGet, "/reports", " Analyst ")
	assertJSONMessage(t, response, http.StatusOK, "reports ok")
}

func performRoleRequest(method string, path string, role string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)
	if role != "" {
		request.Header.Set(RoleHeader, role)
	}
	response := httptest.NewRecorder()

	NewServer().ServeHTTP(response, request)
	return response
}

func echoTestServerWithUnmappedRoute() *echo.Echo {
	e := echo.New()
	e.Use(RBACMiddleware(map[string][]string{
		"/admin": {"admin"},
	}))
	e.GET("/billing", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "billing ok"})
	})
	return e
}

func assertJSONMessage(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantMessage string) {
	t.Helper()
	if response.Code != wantStatus {
		t.Fatalf("expected %d, got %d", wantStatus, response.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload["message"] != wantMessage {
		t.Fatalf("got %+v; want message %q", payload, wantMessage)
	}
}