package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const JWTClaimsContextKey = "jwt_claims"

func JWTBearerAuth(secretKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Flow:
			//   read Authorization header
			//      |
			//      +-> missing or malformed bearer token -> return 401
			//      `-> parse HS256 token -> validate exp -> store claims -> call next
			tokenString, err := extractBearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid bearer token")
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(secretKey), nil
				},
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			)
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid bearer token")
			}

			expirationTime, err := claims.GetExpirationTime()
			if err != nil || expirationTime == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid bearer token")
			}

			c.Set(JWTClaimsContextKey, claims)
			return next(c)
		}
	}
}

func extractBearerToken(authorizationHeader string) (string, error) {
	parts := strings.Fields(strings.TrimSpace(authorizationHeader))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization header must use bearer scheme")
	}
	if parts[1] == "" {
		return "", fmt.Errorf("bearer token is required")
	}
	return parts[1], nil
}

func profileHandler(c echo.Context) error {
	claims, ok := c.Get(JWTClaimsContextKey).(jwt.MapClaims)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "missing jwt claims")
	}

	subject, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)

	return c.JSON(http.StatusOK, map[string]string{
		"sub":  subject,
		"role": role,
	})
}

func newServer(secretKey string) *echo.Echo {
	e := echo.New()
	e.Use(JWTBearerAuth(secretKey))
	e.GET("/profile", profileHandler)
	return e
}

func makeExampleToken(secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  "demo-user",
		"role": "admin",
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func main() {
	secretKey := "demo-secret"
	token, err := makeExampleToken(secretKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Use this Bearer token for GET /profile: %s", token)

	e := newServer(secretKey)
	e.Logger.Fatal(e.Start(":8080"))
}
