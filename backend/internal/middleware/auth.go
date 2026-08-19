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
			slog.Error("auth: missing/bad Bearer prefix", "header", authHeader)
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		token := strings.TrimPrefix(authHeader, prefix)

		user, err := authProvider.GetUserFromToken(ctx, token)

		if err != nil {
			slog.Error("auth: token verification failed", "err", err)
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		if user == nil {
			slog.Error("auth: user is nil")
			return apperr.ToHTTPError(apperr.Unauthorised("Middleware.Auth", "unauthorised"))
		}

		if !user.IsAdmin {
			slog.Error("auth: user not admin", "user_id", user.ID)
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
