package httpapi_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/laflaretapee/go-orders-api/internal/httpapi"
	"github.com/laflaretapee/go-orders-api/internal/order"
)

type fakeOrderService struct {
	orders []order.Order
}

func (s *fakeOrderService) Create(_ context.Context, input order.CreateInput) (order.Order, error) {
	if input.CustomerName == "" {
		return order.Order{}, order.ErrInvalidCustomerName
	}

	item := order.Order{
		ID:           int64(len(s.orders) + 1),
		CustomerName: input.CustomerName,
		Item:         input.Item,
		Quantity:     input.Quantity,
		PriceCents:   input.PriceCents,
		Status:       "new",
		CreatedAt:    time.Now().UTC(),
	}

	s.orders = append(s.orders, item)

	return item, nil
}

func (s *fakeOrderService) GetByID(_ context.Context, id int64) (order.Order, error) {
	for _, item := range s.orders {
		if item.ID == id {
			return item, nil
		}
	}

	return order.Order{}, order.ErrNotFound
}

func (s *fakeOrderService) List(_ context.Context) ([]order.Order, error) {
	result := make([]order.Order, len(s.orders))
	copy(result, s.orders)

	return result, nil
}

func TestCreateOrder(t *testing.T) {
	service := &fakeOrderService{}
	handler := httpapi.NewHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{"customer_name":"Alice","item":"Keyboard","quantity":2,"price_cents":129900}`))
	rec := httptest.NewRecorder()

	handler.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
}

func TestGetOrderNotFound(t *testing.T) {
	service := &fakeOrderService{}
	handler := httpapi.NewHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/orders/42", nil)
	rec := httptest.NewRecorder()

	handler.Router().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
