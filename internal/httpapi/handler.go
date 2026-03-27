package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/laflaretapee/go-orders-api/internal/order"
)

type OrderService interface {
	Create(ctx context.Context, input order.CreateInput) (order.Order, error)
	GetByID(ctx context.Context, id int64) (order.Order, error)
	List(ctx context.Context) ([]order.Order, error)
}

type Handler struct {
	service OrderService
	mux     *http.ServeMux
}

func NewHandler(service OrderService) *Handler {
	handler := &Handler{
		service: service,
		mux:     http.NewServeMux(),
	}

	handler.routes()

	return handler
}

func (h *Handler) Router() http.Handler {
	return h.mux
}

func (h *Handler) routes() {
	h.mux.HandleFunc("POST /orders", h.createOrder)
	h.mux.HandleFunc("GET /orders", h.listOrders)
	h.mux.HandleFunc("GET /orders/{id}", h.getOrder)
}

func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	items, err := h.service.List(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"orders": items,
	})
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	item, err := h.service.GetByID(ctx, id)
	if err != nil {
		status, message := mapError(err)
		writeError(w, status, message)
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input order.CreateInput
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	created, err := h.service.Create(ctx, input)
	if err != nil {
		status, message := mapError(err)
		writeError(w, status, message)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, order.ErrInvalidCustomerName),
		errors.Is(err, order.ErrInvalidItem),
		errors.Is(err, order.ErrInvalidQuantity),
		errors.Is(err, order.ErrInvalidPriceCents):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, order.ErrNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
