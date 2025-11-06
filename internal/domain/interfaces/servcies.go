package interfaces

import (
	"context"
	"ooliokartchallenge/internal/domain/entities"
)

type ProductService interface {
	ListProducts(ctx context.Context) ([]entities.Product, error)
	GetProduct(ctx context.Context, id string) (*entities.Product, error)
}

type OrderService interface {
	PlaceOrder(ctx context.Context, req entities.OrderRequest) (*entities.Order, error)
}

type PromoService interface {
	ValidatePromoCode(ctx context.Context, code string) (bool, error)
}
