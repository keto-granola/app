package admin

import (
	"context"
)

type Service struct {
	store Repository
}

func NewService(store Repository) *Service {
	return &Service{store: store}
}

func (s *Service) GetLowStock(ctx context.Context) ([]LowStockItem, error) {
	products, err := s.store.GetLowStock(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}
