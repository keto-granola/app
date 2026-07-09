package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/keto-granola/keto-granola/internal/services/auth"
)

func Auth(next echo.HandlerFunc, authProvider auth.AuthProvider) echo.HandlerFunc {
	// TODO: read access token from authorization header
	return func(c echo.Context) error {
		return next(c)
	}
}

func TestAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: implement
		return next(c)
	}
}
