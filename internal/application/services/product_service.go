package services

import (
	"context"
	"fmt"
	"ooliokartchallenge/internal/domain/entities"
	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/internal/domain/interfaces"
	"strconv"
)

type ProductService struct {
	productRepo interfaces.ProductRepository
}

func NewProductService(productRepo interfaces.ProductRepository) interfaces.ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]entities.Product, error) {
	products, err := s.productRepo.GetAll(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}

	return products, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*entities.Product, error) {
	if _, err := strconv.Atoi(id); err != nil {
		return nil, errors.ErrInvalidProductID
	}
	product, err := s.productRepo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve product: %w", err)
	}

	return product, nil
}
