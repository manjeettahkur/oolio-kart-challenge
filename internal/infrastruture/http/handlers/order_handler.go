package handlers

import (
	"encoding/json"
	"net/http"
	"ooliokartchallenge/internal/domain/entities"
	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/internal/domain/interfaces"
	"ooliokartchallenge/pkg/logger"
)

type OrderHandler struct {
	orderService interfaces.OrderService
	logger       *logger.Logger
}

func NewOrderHandler(orderService interfaces.OrderService, log *logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       log,
	}
}

// PlaceOrder handles POST /order requests to create a new order
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var orderRequest entities.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		HandleError(w, r, errors.ErrInvalidJSON, h.logger)
		return
	}

	order, err := h.orderService.PlaceOrder(ctx, orderRequest)
	if err != nil {
		HandleError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		HandleError(w, r, err, h.logger)
		return
	}
}
