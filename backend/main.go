package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/keto-granola/keto-granola/internal/config"
	"github.com/keto-granola/keto-granola/internal/product"
	productadmin "github.com/keto-granola/keto-granola/internal/product/admin"
	productstore "github.com/keto-granola/keto-granola/internal/product/store"
	"github.com/keto-granola/keto-granola/internal/product/web"
	"github.com/keto-granola/keto-granola/internal/server"
	"github.com/keto-granola/keto-granola/internal/store"
	"github.com/keto-granola/keto-granola/internal/webassets"
)

func main() {
	if err := run(); err != nil {
		slog.Error("run", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.ParseEnv()
	if err != nil {
		return fmt.Errorf("parse env vars %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	dataStore, err := store.New(ctx, cfg.DbURL)
	if err != nil {
		return fmt.Errorf("create store %w", err)
	}
	defer dataStore.Close()

	assetsLoader, err := webassets.New()
	if err != nil {
		return fmt.Errorf("init asset loader: %w", err)
	}

	templates, err := server.NewTemplates()
	if err != nil {
		return fmt.Errorf("create templates: %w", err)
	}

	handlers := composeHandlers(
		dataStore,
		templates,
		assetsLoader,
		cfg.ClientURL,
		cfg.Environment,
	)

	serverDeps := &server.Dependencies{
		Environment: cfg.Environment,
		ClientURL:   cfg.ClientURL,
		Handlers:    handlers,
		DataStore:   dataStore,
	}

	echo, err := server.New(ctx, serverDeps)
	if err != nil {
		return err
	}

	serverStartErr := make(chan error, 1)

	go func() {
		serverStartErr <- echo.Start(cfg.Port)
	}()

	select {
	case <-ctx.Done():
		slog.Info("context cancelled")
	case err = <-serverStartErr:
	}

	if shutdownErr := echo.Stop(); shutdownErr != nil {
		slog.Error("server shutdown", slog.Any("error", shutdownErr))
	}

	return err
}

func composeHandlers(
	db *store.Store,
	tmpl *template.Template,
	assetsLoader *webassets.Loader,
	clientURL string,
	env config.Environment,
) *server.Handlers {
	productStore := productstore.New(db.Queries)
	prodService := product.NewService(productStore)
	prodAdminService := productadmin.NewService(productStore)

	return &server.Handlers{
		ProductAdmin: productadmin.NewHandler(prodAdminService),
		Product:      web.NewHandler(prodService, assetsLoader, tmpl, clientURL, env),
	}
}
