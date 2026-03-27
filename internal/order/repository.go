package order

import "context"

type Repository interface {
	Create(ctx context.Context, input CreateInput) (Order, error)
	GetByID(ctx context.Context, id int64) (Order, error)
	List(ctx context.Context) ([]Order, error)
}
