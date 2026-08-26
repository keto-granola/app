package mcp

import (
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	inventoryadmin "github.com/keto-granola/keto-granola/internal/inventory/admin"
)

type Server struct {
	streamable *server.StreamableHTTPServer
}

func NewServer(inventoryAdminService *inventoryadmin.Service) *Server {
	s := server.NewMCPServer(
		"keto-granola-store",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	registerInventoryTools(s, inventoryAdminService)

	return &Server{
		streamable: server.NewStreamableHTTPServer(s),
	}
}

func (s *Server) ServeHTTP(c echo.Context) error {
	s.streamable.ServeHTTP(c.Response(), c.Request())

	return nil
}
