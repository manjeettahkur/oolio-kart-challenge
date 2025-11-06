package handlers

import (
	"encoding/json"
	"net/http"
	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/internal/domain/interfaces"
	"ooliokartchallenge/pkg/logger"
)

type ProductHandler struct {
	productService interfaces.ProductService
	logger         *logger.Logger
}

func NewProductHandler(productService interfaces.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         log,
	}
}

// ListProducts handles GET /product requests to return all products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	products, err := h.productService.ListProducts(ctx)
	if err != nil {
		HandleError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(products); err != nil {
		HandleError(w, r, err, h.logger)
		return
	}
}

// GetProduct handles GET /product/{productId} requests to return a specific product
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productID := r.PathValue("id")

	if productID == "" {
		HandleError(w, r, errors.ErrInvalidProductID, h.logger)
		return
	}

	product, err := h.productService.GetProduct(ctx, productID)
	if err != nil {
		HandleError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(product); err != nil {
		HandleError(w, r, err, h.logger)
		return
	}
}
