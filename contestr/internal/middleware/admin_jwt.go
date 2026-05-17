package middleware

import (
	"contestr/internal/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AdminJWT(authService *auth.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if !strings.HasPrefix(path, "/api/admin") {
				return next(c)
			}
			if c.Request().Method == http.MethodPost && path == "/api/admin/auth/login" {
				return next(c)
			}

			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "missing authorization header",
				})
			}

			claims, err := authService.ParseToken(header)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid or expired token",
				})
			}

			c.Set(auth.ContextUsernameKey, claims.Username)
			return next(c)
		}
	}
}
