package order

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidCustomerName = errors.New("customer_name is required")
	ErrInvalidItem         = errors.New("item is required")
	ErrInvalidQuantity     = errors.New("quantity must be greater than 0")
	ErrInvalidPriceCents   = errors.New("price_cents must be greater than 0")
	ErrNotFound            = errors.New("order not found")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Order, error) {
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.Item = strings.TrimSpace(input.Item)

	switch {
	case input.CustomerName == "":
		return Order{}, ErrInvalidCustomerName
	case input.Item == "":
		return Order{}, ErrInvalidItem
	case input.Quantity <= 0:
		return Order{}, ErrInvalidQuantity
	case input.PriceCents <= 0:
		return Order{}, ErrInvalidPriceCents
	default:
		return s.repo.Create(ctx, input)
	}
}

func (s *Service) GetByID(ctx context.Context, id int64) (Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Order, error) {
	return s.repo.List(ctx)
}
