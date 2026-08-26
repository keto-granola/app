package store

import (
	"context"

	"github.com/keto-granola/keto-granola/internal/apperr"
	"github.com/keto-granola/keto-granola/internal/inventory/admin"
	"github.com/keto-granola/keto-granola/internal/store"
	"github.com/keto-granola/keto-granola/internal/store/db/generated"
	"github.com/keto-granola/keto-granola/internal/store/db/utils"
)

type Store struct {
	queries *generated.Queries
}

func New(queries *generated.Queries) *Store {
	return &Store{queries: queries}
}

func (s *Store) GetLowStock(ctx context.Context) ([]admin.LowStockItem, error) {
	rows, err := store.ExecWithResult(ctx, func() ([]generated.GetLowStockRow, error) {
		return s.queries.GetLowStock(ctx)
	})

	if err != nil {
		return nil, apperr.Internal("Store.GetLowStock", err)
	}

	products := make([]admin.LowStockItem, 0, len(rows))

	for _, row := range rows {
		products = append(products, lowStockFrom(&row))
	}

	return products, nil
}

func lowStockFrom(row *generated.GetLowStockRow) admin.LowStockItem {
	return admin.LowStockItem{
		ProductID:   utils.UUIDFrom(row.ProductID),
		ProductName: row.ProductName,
		Quantity:    row.Quantity,
		Threshold:   row.RestockThreshold,
	}
}
