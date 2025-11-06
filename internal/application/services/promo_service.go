package services

import (
	"context"
	"fmt"
	"ooliokartchallenge/internal/domain/interfaces"
)

type PromoService struct {
	promoRepo interfaces.PromoRepository
}

func NewPromoService(promoRepo interfaces.PromoRepository) interfaces.PromoService {
	return &PromoService{
		promoRepo: promoRepo,
	}
}

func (s *PromoService) ValidatePromoCode(ctx context.Context, code string) (bool, error) {

	if len(code) < 8 {
		return false, nil
	}
	if len(code) > 10 {
		return false, nil
	}

	exists, err := s.promoRepo.ValidateCode(ctx, code)
	if err != nil {
		return false, fmt.Errorf("failed to validate promo code: %w", err)
	}

	return exists, nil
}
