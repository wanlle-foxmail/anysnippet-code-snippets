package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const RoleHeader = "X-User-Role"

func RBACMiddleware(allowedRolesByRoute map[string][]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			allowedRoles, ok := allowedRolesByRoute[c.Path()]
			if !ok {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "route is not allowed"})
			}

			role := strings.ToLower(strings.TrimSpace(c.Request().Header.Get(RoleHeader)))
			if role == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "role is required"})
			}

			for _, allowedRole := range allowedRoles {
				if role == strings.ToLower(strings.TrimSpace(allowedRole)) {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{"message": "role is not allowed"})
		}
	}
}

func NewServer() *echo.Echo {
	e := echo.New()
	e.Use(RBACMiddleware(map[string][]string{
		"/admin":   {"admin"},
		"/reports": {"admin", "analyst"},
	}))
	e.GET("/admin", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "admin ok"})
	})
	e.GET("/reports", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "reports ok"})
	})
	return e
}

func main() {
	e := NewServer()
	e.Logger.Fatal(e.Start(":8080"))
}