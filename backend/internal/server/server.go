package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/keto-granola/keto-granola/internal/config"
	"github.com/keto-granola/keto-granola/internal/middleware"
	productadmin "github.com/keto-granola/keto-granola/internal/product/admin"
	productweb "github.com/keto-granola/keto-granola/internal/product/web"
	"github.com/keto-granola/keto-granola/internal/server/templates/templatehelpers"
	"github.com/keto-granola/keto-granola/internal/store"
)

const (
	// how long the server will wait to read the entire request after the connection is accepted
	readTimeout = 10 * time.Second
	// how long the server has to write the response after reading the request
	writeTimeout = 10 * time.Second
	// how long to keep a keep-alive connection open waiting for the next request
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 10 * time.Second

	serverRateLimit  = 60
	serverBurstLimit = 120
)

type Dependencies struct {
	Environment config.Environment
	ClientURL   string
	Handlers    *Handlers
	DataStore   *store.Store
}

type Server struct {
	echo *echo.Echo
}

type Handlers struct {
	// api handlers
	ProductAdmin *productadmin.Handler

	// web handlers
	Product *productweb.Handler
}

//go:embed templates
var templateFS embed.FS

func New(ctx context.Context, deps *Dependencies) (*Server, error) {
	e := echo.New()
	NewValidator(e)
	e.HideBanner = true // prevents startup banner from being logged

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{deps.ClientURL},
	}))

	e.Use(middleware.Log)

	// limits each unique IP to 60 requests per minute with a burst of 120.
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStoreWithConfig(
		echoMiddleware.RateLimiterMemoryStoreConfig{
			Rate:      serverRateLimit,
			Burst:     serverBurstLimit,
			ExpiresIn: time.Minute,
		},
	)))

	web := e.Group("")
	api := e.Group(config.APIBasePath)

	apiPublic := api.Group("")
	apiPrivate := api.Group("")

	if deps.Environment == config.EnvironmentTest {
		// TODO: run test middleware
		slog.Info("run test middleware")
	} else {
		// TODO: run auth middleware
		slog.Info("run auth middleware")
	}

	if err := registerRoutes(apiPublic, apiPrivate, web, deps.Handlers, deps.DataStore); err != nil {
		return nil, err
	}

	e.Server.ReadTimeout = readTimeout
	e.Server.WriteTimeout = writeTimeout
	e.Server.IdleTimeout = idleTimeout

	return &Server{
		echo: e,
	}, nil
}

func (s *Server) Start(port string) error {
	if err := s.echo.Start(":" + port); err != nil && err != http.ErrServerClosed {
		slog.Error("start server", slog.Any("error", err))
		return err
	}

	return nil
}

func NewTemplates() (*template.Template, error) {
	tmpl, err := template.New("").Funcs(templatehelpers.FuncMap()).ParseFS(templateFS, "templates/**/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	return tmpl, nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.echo.Shutdown(ctx)
}
