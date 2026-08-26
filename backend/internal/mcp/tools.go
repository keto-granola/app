package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/keto-granola/keto-granola/internal/inventory/admin"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerInventoryTools(server *server.MCPServer, service *admin.Service) {
	handler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		products, err := service.GetLowStock(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(productsFrom(products)), nil
	}

	server.AddTool(
		mcp.NewTool("get_low_stock_products",
			mcp.WithDescription("Returns products where stock is below restock threshold. Use when the user asks about restocking or low inventory."),
		),
		handler,
	)
}

func productsFrom(products []admin.LowStockItem) string {
	if len(products) == 0 {
		return `{"message": "No products currently below restock threshold."}`
	}

	data, err := json.Marshal(products)
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to format results: %s"}`, err)
	}

	return string(data)
}
