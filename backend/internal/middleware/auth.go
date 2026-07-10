package middleware

import (
	"context"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/keto-granola/keto-granola/internal/apperr"
	"github.com/keto-granola/keto-granola/internal/services/auth"
)

func Auth(next echo.HandlerFunc, authProvider auth.AuthProvider) echo.HandlerFunc {
	return func(e echo.Context) error {
		ctx := e.Request().Context()

		authHeader := e.Request().Header.Get("Authorization")

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			slog.Warn("missing or invalid auth token")
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		token := strings.TrimPrefix(authHeader, prefix)

		user, err := authProvider.GetUserByToken(ctx, token)

		if err != nil {
			slog.Warn("get user from token", slog.Any("error", err))
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		if user == nil {
			slog.Warn("user not found from token")
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		if !user.IsAdmin {
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		role := getUserRole(user.IsAdmin)
		ctx = context.WithValue(ctx, auth.RoleContextKey, role)
		ctx = context.WithValue(ctx, auth.UserIDContextKey, user.ID)
		e.SetRequest(e.Request().WithContext(ctx))

		return next(e)
	}
}

func TestAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		// TODO: implement
		return next(e)
	}
}

func getUserRole(isAdmin bool) auth.Role {
	if isAdmin {
		return auth.AdminRole
	}

	return auth.UserRole
}
