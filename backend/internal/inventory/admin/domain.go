package admin

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out mocks/mock.go -pkg mocks . Repository

type Repository interface {
	GetLowStock(ctx context.Context) ([]LowStockItem, error)
}

type LowStockItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int32     `json:"quantity"`
	Threshold   int32     `json:"threshold"`
}
