package mcp

import (
	"crypto/subtle"

	"github.com/keto-granola/keto-granola/internal/apperr"
	"github.com/labstack/echo/v4"
)

func RequireMCPToken(token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			expected := "Bearer " + token

			if subtle.ConstantTimeCompare([]byte(auth), []byte(expected)) != 1 {
				return apperr.ToHTTPError(apperr.Unauthorised("Middleware.MCP", "unauthorised"))
			}

			return next(c)
		}
	}
}
