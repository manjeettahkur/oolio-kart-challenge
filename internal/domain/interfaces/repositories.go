package interfaces

import (
	"context"
	"ooliokartchallenge/internal/domain/entities"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]entities.Product, error)
	GetByID(ctx context.Context, id string) (*entities.Product, error)
}

type PromoRepository interface {
	ValidateCode(ctx context.Context, code string) (bool, error)
}
