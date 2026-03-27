package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/laflaretapee/go-orders-api/internal/order"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, input order.CreateInput) (order.Order, error) {
	const query = `
		INSERT INTO orders (customer_name, item, quantity, price_cents, status)
		VALUES ($1, $2, $3, $4, 'new')
		RETURNING id, customer_name, item, quantity, price_cents, status, created_at
	`

	var item order.Order
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.CustomerName,
		input.Item,
		input.Quantity,
		input.PriceCents,
	).Scan(
		&item.ID,
		&item.CustomerName,
		&item.Item,
		&item.Quantity,
		&item.PriceCents,
		&item.Status,
		&item.CreatedAt,
	)
	if err != nil {
		return order.Order{}, err
	}

	return item, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (order.Order, error) {
	const query = `
		SELECT id, customer_name, item, quantity, price_cents, status, created_at
		FROM orders
		WHERE id = $1
	`

	var item order.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.CustomerName,
		&item.Item,
		&item.Quantity,
		&item.PriceCents,
		&item.Status,
		&item.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return order.Order{}, order.ErrNotFound
	}
	if err != nil {
		return order.Order{}, err
	}

	return item, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]order.Order, error) {
	const query = `
		SELECT id, customer_name, item, quantity, price_cents, status, created_at
		FROM orders
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]order.Order, 0)
	for rows.Next() {
		var item order.Order
		if err := rows.Scan(
			&item.ID,
			&item.CustomerName,
			&item.Item,
			&item.Quantity,
			&item.PriceCents,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
