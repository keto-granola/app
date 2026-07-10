package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/keto-granola/keto-granola/internal/services/auth"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("load .env file", slog.Any("error", err))
	}

	var makeAdmin bool
	flag.BoolVar(&makeAdmin, "make-admin", makeAdmin, "Make admin if true/remove admin status if false")
	flag.Parse()

	if err := run(os.Getenv("FIREBASE_CREDS"), makeAdmin); err != nil {
		slog.Error("run failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(creds string, makeAdmin bool) error {
	if creds == "" {
		return fmt.Errorf("FIREBASE_CREDS is not set")
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the email of the user:\n")

	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read email input %v", err)
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email must not be empty")
	}

	fmt.Printf("About to set admin=%v for %s. Continue? [y/n]: ", makeAdmin, email)
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read confirmation: %w", err)
	}
	if strings.TrimSpace(strings.ToLower(confirm)) != "y" {
		return fmt.Errorf("aborted by user")
	}

	ctx := context.Background()

	authProv, err := auth.New(ctx, creds)
	if err != nil {
		return fmt.Errorf("initialise auth provider %v", err)
	}

	userID, err := authProv.GetUserIDByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("fetch user %v", err)
	}
	slog.Info("fetched userID", slog.String("email", email), slog.String("id", userID))

	err = authProv.SetAdminRole(ctx, userID, makeAdmin)
	if err != nil {
		return fmt.Errorf("update user admin status %v", err)
	}

	isAdmin, err := authProv.IsAdminFromUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("verify admin status %v", err)
	}
	if isAdmin != makeAdmin {
		return fmt.Errorf("admin status mismatch after update: want %v, got %v", makeAdmin, isAdmin)
	}

	slog.Info("updated user admin status", slog.String("email", email), slog.Bool("make_admin", makeAdmin))
	return nil
}
