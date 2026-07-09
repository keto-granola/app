package auth

import (
	"context"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Token string
type Role string

const adminRole Role = "admin"

type Auth struct {
	client *auth.Client
}

type AuthProvider interface {
	GetUserFromIDToken(ctx context.Context, token string) (*User, error)
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

func (a *Auth) GetUserFromIDToken(ctx context.Context, accessToken string) (*User, error) {
	token, err := a.client.VerifyIDToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	userRecord, err := a.client.GetUser(ctx, token.UID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:      userRecord.UID,
		IsAdmin: isAdmin(token),
	}, nil
}

func isAdmin(token *auth.Token) bool {
	adminValue, ok := token.Claims[string(adminRole)]
	if !ok {
		return false
	}

	isAdmin, ok := adminValue.(bool)
	return ok && isAdmin
}
