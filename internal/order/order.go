package order

import "time"

type Order struct {
	ID           int64     `json:"id"`
	CustomerName string    `json:"customer_name"`
	Item         string    `json:"item"`
	Quantity     int       `json:"quantity"`
	PriceCents   int64     `json:"price_cents"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateInput struct {
	CustomerName string `json:"customer_name"`
	Item         string `json:"item"`
	Quantity     int    `json:"quantity"`
	PriceCents   int64  `json:"price_cents"`
}
