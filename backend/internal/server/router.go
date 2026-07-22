package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/keto-granola/keto-granola/internal/store"
	"github.com/keto-granola/keto-granola/internal/webassets"
)

const pingTimeout = 5 * time.Second

func registerRoutes(apiPublic, apiPrivate, web *echo.Group, handlers *Handlers, dataStore *store.Store) error {
	registerHealthEndpoint(apiPublic, dataStore)

	if err := registerAssetRoutes(web); err != nil {
		return err
	}

	registerAPIRoutes(apiPrivate, handlers)

	registerWebRoutes(web, handlers)

	return nil
}

func registerHealthEndpoint(api *echo.Group, dataStore *store.Store) {
	api.GET("/health", func(e echo.Context) error {
		dbStatus := "ok"
		httpStatus := http.StatusOK

		pingCtx, cancel := context.WithTimeout(e.Request().Context(), pingTimeout)
		defer cancel()

		err := dataStore.PingDB(pingCtx)
		if err != nil {
			httpStatus = http.StatusServiceUnavailable
			dbStatus = "unreachable"
		}

		return e.JSON(httpStatus, map[string]string{
			"status": "ok",
			"db":     dbStatus,
		})
	})
}

func registerAssetRoutes(web *echo.Group) error {
	handler, err := webassets.AssetsHandler()
	if err != nil {
		return fmt.Errorf("set up assets handler: %w", err)
	}

	web.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", handler)))

	return nil
}

func registerAPIRoutes(apiPrivate *echo.Group, handlers *Handlers) {
	registerProductAdminRoutes(apiPrivate, handlers)
}

func registerProductAdminRoutes(apiPrivate *echo.Group, handlers *Handlers) {
	products := apiPrivate.Group("/admin/products")

	products.POST("", Handle(handlers.ProductAdmin.CreateProduct, http.StatusCreated))
}

func registerWebRoutes(web *echo.Group, handlers *Handlers) {
	web.GET("products/:id", handlers.Product.GetProductPage)

	web.GET("/admin", handlers.Admin.ServeShell)
	web.GET("/admin/*", handlers.Admin.ServeShell)
}
