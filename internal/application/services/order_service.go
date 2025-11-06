package services

import (
	"context"
	"fmt"
	"ooliokartchallenge/internal/domain/entities"
	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/internal/domain/interfaces"
	"strings"
	"time"
)

type OrderService struct {
	productRepo  interfaces.ProductRepository
	promoService interfaces.PromoService
}

func NewOrderService(productRepo interfaces.ProductRepository, promoService interfaces.PromoService) interfaces.OrderService {
	return &OrderService{
		productRepo:  productRepo,
		promoService: promoService,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, req entities.OrderRequest) (*entities.Order, error) {


	if err := s.validateOrderRequest(req); err != nil {
		return nil, err
	}

	orderProducts, totalAmount, err := s.validateAndCalculateItems(ctx, req.Items)

	if err != nil {
		return nil, err
	}

	discountAmount, err := s.applyPromoCodeDiscount(ctx, req.CouponCode, totalAmount)

	if err != nil {
		return nil, err
	}

	finalTotal := totalAmount - discountAmount

	orderID := fmt.Sprintf("order_%d", time.Now().UnixNano())

	order := &entities.Order{
		ID:       orderID,
		Total:    finalTotal,
		Discounts: discountAmount,
		Items:    req.Items,
		Products: orderProducts,
	}

	return order, nil
}

func (s *OrderService) validateOrderRequest(req entities.OrderRequest) error {

	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errors.ErrInvalidOrderRequest, err)
	}

	if req.CouponCode != "" {
		if err := s.validateCouponCodeFormat(req.CouponCode); err != nil {
			return fmt.Errorf("%w: %v", errors.ErrInvalidPromoCode, err)
		}
	}

	return nil
}

func (s *OrderService) validateCouponCodeFormat(couponCode string) error {

	trimmed := strings.TrimSpace(couponCode)
	if trimmed == "" {
		return fmt.Errorf("coupon code cannot be empty or whitespace only")
	}

	if len(trimmed) < 8 {
		return errors.ErrPromoCodeTooShort
	}
	if len(trimmed) > 10 {
		return errors.ErrPromoCodeTooLong
	}

	for _, char := range trimmed {
		if char < 'A' || char > 'Z' {
			return fmt.Errorf("promo code must contain only uppercase letters (no numbers or special characters)")
		}
	}

	return nil
}

func (s *OrderService) validateAndCalculateItems(ctx context.Context, items []entities.OrderItem) ([]entities.Product, float64, error) {

	var orderProducts []entities.Product

	var totalAmount float64

	for idx, item := range items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)

		if err != nil {
			return nil, 0, fmt.Errorf("%w: item %d not exits", errors.ErrInvalidProductID, idx)
		}

		orderProducts = append(orderProducts, *product)

		itemAmount := product.Price * float64(item.Quantity)
		totalAmount += itemAmount

	}

	return orderProducts, totalAmount, nil
}

func (s *OrderService) applyPromoCodeDiscount(ctx context.Context, couponCode string, totalAmount float64) (float64, error) {

	if strings.TrimSpace(couponCode) == "" {
		return 0, nil
	}

	isValid, err := s.promoService.ValidatePromoCode(ctx, couponCode)
	if err != nil {
		return 0, fmt.Errorf("failed to validate promo code: %w", err)
	}

	if !isValid {
		return 0, fmt.Errorf("%w: code '%s' is not valid", errors.ErrInvalidPromoCode, couponCode)
	}

	discountAmount := totalAmount * 0.10

	return discountAmount, nil
}
