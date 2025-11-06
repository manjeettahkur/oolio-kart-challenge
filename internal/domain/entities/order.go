package entities

import (
	"errors"
	"fmt"
	"strings"
)

type Order struct {
	ID       string      `json:"id"`
	Total    float64     `json:"total"`
	Discounts float64     `json:"discounts"`
	Items    []OrderItem `json:"items"`
	Products []Product   `json:"products"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (o *OrderItem) Validate() error {
	if strings.TrimSpace(o.ProductID) == "" {
		return errors.New("productId is required")
	}

	if o.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return nil
}

type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items"`
}

func (or *OrderRequest) Validate() error {
	if len(or.Items) == 0 {
		return errors.New("items are required")
	}

	for i, item := range or.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item at index %d: %w", i, err)
		}
	}

	return nil
}
