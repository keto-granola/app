package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Token string
type Role string

type ContextKey string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"

	UserIDContextKey ContextKey = "userID"
	RoleContextKey   ContextKey = "role"
)

type Auth struct {
	client *auth.Client
}

type AuthProvider interface {
	GetUserFromToken(ctx context.Context, token string) (*User, error)
}

type User struct {
	ID      string
	IsAdmin bool
}

func New(ctx context.Context, credentials string) (*Auth, error) {
	opts := option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(credentials))

	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &Auth{
		client: client,
	}, nil
}

func (a *Auth) GetUserFromToken(ctx context.Context, accessToken string) (*User, error) {
	token, err := a.client.VerifyIDToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("verify id token: %w", err)
	}

	isAdmin, _ := token.Claims[string(AdminRole)].(bool)

	return &User{
		ID:      token.UID,
		IsAdmin: isAdmin,
	}, nil
}

func (a *Auth) GetUserIDByEmail(ctx context.Context, email string) (userID string, err error) {
	user, err := a.client.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	return user.UID, nil
}

func (a *Auth) SetAdminRole(ctx context.Context, userID string, setAdmin bool) error {
	return a.client.SetCustomUserClaims(ctx, userID, map[string]any{
		"admin": setAdmin,
	})
}

func (a *Auth) IsAdminFromUserID(ctx context.Context, userID string) (bool, error) {
	user, err := a.client.GetUser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("get user: %w", err)
	}

	admin, _ := user.CustomClaims["admin"].(bool)

	return admin, nil
}
