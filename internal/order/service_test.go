package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/laflaretapee/go-orders-api/internal/order"
)

type memoryRepository struct {
	orders []order.Order
	nextID int64
}

func (r *memoryRepository) Create(_ context.Context, input order.CreateInput) (order.Order, error) {
	r.nextID++

	item := order.Order{
		ID:           r.nextID,
		CustomerName: input.CustomerName,
		Item:         input.Item,
		Quantity:     input.Quantity,
		PriceCents:   input.PriceCents,
		Status:       "new",
		CreatedAt:    time.Now().UTC(),
	}

	r.orders = append(r.orders, item)

	return item, nil
}

func (r *memoryRepository) GetByID(_ context.Context, id int64) (order.Order, error) {
	for _, item := range r.orders {
		if item.ID == id {
			return item, nil
		}
	}

	return order.Order{}, order.ErrNotFound
}

func (r *memoryRepository) List(_ context.Context) ([]order.Order, error) {
	result := make([]order.Order, len(r.orders))
	copy(result, r.orders)

	return result, nil
}

func TestCreateOrder(t *testing.T) {
	repo := &memoryRepository{}
	service := order.NewService(repo)

	created, err := service.Create(context.Background(), order.CreateInput{
		CustomerName: "  Ivan Petrov  ",
		Item:         "Laptop",
		Quantity:     1,
		PriceCents:   7999000,
	})
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}

	if created.CustomerName != "Ivan Petrov" {
		t.Fatalf("expected trimmed customer name, got %q", created.CustomerName)
	}
}

func TestCreateOrderRejectsInvalidInput(t *testing.T) {
	repo := &memoryRepository{}
	service := order.NewService(repo)

	_, err := service.Create(context.Background(), order.CreateInput{
		CustomerName: "",
		Item:         "Laptop",
		Quantity:     1,
		PriceCents:   100,
	})
	if err != order.ErrInvalidCustomerName {
		t.Fatalf("expected ErrInvalidCustomerName, got %v", err)
	}
}
