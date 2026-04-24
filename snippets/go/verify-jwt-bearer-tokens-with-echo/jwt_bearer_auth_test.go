package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const testJWTSecret = "test-secret"

func signedToken(t *testing.T, secretKey string, claims jwt.MapClaims, method jwt.SigningMethod) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return tokenString
}

func performProfileRequest(t *testing.T, authorizationHeader string) (*httptest.ResponseRecorder, map[string]string) {
	t.Helper()

	e := newServer(testJWTSecret)
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	if authorizationHeader != "" {
		req.Header.Set(echo.HeaderAuthorization, authorizationHeader)
	}
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	var body map[string]string
	if rec.Body.Len() > 0 && rec.Code == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode response body: %v", err)
		}
	}

	return rec, body
}

func TestJWTBearerAuthAllowsValidToken(t *testing.T) {
	token := signedToken(
		t,
		testJWTSecret,
		jwt.MapClaims{
			"sub":  "user-123",
			"role": "admin",
			"exp":  time.Now().Add(10 * time.Minute).Unix(),
		},
		jwt.SigningMethodHS256,
	)

	rec, body := performProfileRequest(t, "Bearer "+token)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if body["sub"] != "user-123" {
		t.Fatalf("expected sub user-123, got %q", body["sub"])
	}
	if body["role"] != "admin" {
		t.Fatalf("expected role admin, got %q", body["role"])
	}
}

func TestJWTBearerAuthRejectsMissingOrMalformedAuthorizationHeader(t *testing.T) {
	tests := []struct {
		name                string
		authorizationHeader string
	}{
		{name: "missing header", authorizationHeader: ""},
		{name: "basic scheme", authorizationHeader: "Basic abc.def.ghi"},
		{name: "missing token", authorizationHeader: "Bearer"},
		{name: "whitespace token", authorizationHeader: "Bearer    "},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			rec, _ := performProfileRequest(t, testCase.authorizationHeader)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("expected status 401, got %d", rec.Code)
			}
		})
	}
}

func TestJWTBearerAuthRejectsInvalidTokenFormat(t *testing.T) {
	rec, _ := performProfileRequest(t, "Bearer not-a-jwt")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestJWTBearerAuthRejectsTokenWithInvalidSignature(t *testing.T) {
	token := signedToken(
		t,
		"wrong-secret",
		jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(10 * time.Minute).Unix(),
		},
		jwt.SigningMethodHS256,
	)

	rec, _ := performProfileRequest(t, "Bearer "+token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestJWTBearerAuthRejectsExpiredToken(t *testing.T) {
	token := signedToken(
		t,
		testJWTSecret,
		jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(-10 * time.Minute).Unix(),
		},
		jwt.SigningMethodHS256,
	)

	rec, _ := performProfileRequest(t, "Bearer "+token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestJWTBearerAuthRejectsTokenWithoutExpiration(t *testing.T) {
	token := signedToken(
		t,
		testJWTSecret,
		jwt.MapClaims{
			"sub": "user-123",
		},
		jwt.SigningMethodHS256,
	)

	rec, _ := performProfileRequest(t, "Bearer "+token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestJWTBearerAuthRejectsUnexpectedSigningMethod(t *testing.T) {
	token := signedToken(
		t,
		testJWTSecret,
		jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(10 * time.Minute).Unix(),
		},
		jwt.SigningMethodHS384,
	)

	rec, _ := performProfileRequest(t, "Bearer "+token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}
