package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) APIError {
	return APIError{
		Code:    code,
		Type:    http.StatusText(code),
		Message: message,
	}
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

var (
	// Product errors
	ErrProductNotFound  = errors.New("product not found")
	ErrInvalidProductID = errors.New("invalid product ID")

	// Order errors
	ErrInvalidOrderRequest = errors.New("invalid order request")
	ErrEmptyOrderItems     = errors.New("order must contain at least one item")
	ErrInvalidQuantity     = errors.New("quantity must be greater than 0")
	ErrInvalidProductRef   = errors.New("invalid product reference in order")

	// Promo code errors
	ErrInvalidPromoCode  = errors.New("invalid promo code")
	ErrPromoCodeTooShort = errors.New("promo code must be at least 8 characters")
	ErrPromoCodeTooLong  = errors.New("promo code must be at most 10 characters")
	ErrPromoCodeNotFound = errors.New("promo code not found in sufficient databases")

	// Authentication errors
	ErrUnauthorized  = errors.New("unauthorized")
	ErrMissingAPIKey = errors.New("missing API key")
	ErrInvalidAPIKey = errors.New("invalid API key")

	// Validation errors
	ErrValidationFailed = errors.New("validation failed")
	ErrRequiredField    = errors.New("required field missing")
	ErrInvalidFormat    = errors.New("invalid format")
	ErrInvalidJSON      = errors.New("invalid JSON format")
	ErrDuplicateItem    = errors.New("duplicate item")
	ErrExceedsLimit     = errors.New("exceeds allowed limit")

	// Internal errors
	ErrInternalServer = errors.New("internal server error")
)

func MapErrorToAPIError(err error) APIError {
	if err == nil {
		return APIError{}
	}

	if apiErr, ok := err.(APIError); ok {
		return apiErr
	}

	switch {
	case errors.Is(err, ErrInvalidProductID):
		return NewAPIError(http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrProductNotFound):
		return NewAPIError(http.StatusNotFound, err.Error())

	case errors.Is(err, ErrInvalidOrderRequest),
		errors.Is(err, ErrEmptyOrderItems),
		errors.Is(err, ErrInvalidQuantity),
		errors.Is(err, ErrInvalidProductRef),
		errors.Is(err, ErrInvalidJSON),
		errors.Is(err, ErrRequiredField),
		errors.Is(err, ErrInvalidFormat):
		return NewAPIError(http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrUnauthorized),
		errors.Is(err, ErrMissingAPIKey),
		errors.Is(err, ErrInvalidAPIKey):
		return NewAPIError(http.StatusUnauthorized, err.Error())

	case errors.Is(err, ErrInvalidPromoCode),
		errors.Is(err, ErrPromoCodeTooShort),
		errors.Is(err, ErrPromoCodeTooLong),
		errors.Is(err, ErrPromoCodeNotFound):
		return NewAPIError(http.StatusUnprocessableEntity, err.Error())

	case errors.Is(err, ErrValidationFailed),
		errors.Is(err, ErrDuplicateItem),
		errors.Is(err, ErrExceedsLimit):
		return NewAPIError(http.StatusBadRequest, err.Error())

	default:
		return NewAPIError(http.StatusInternalServerError, "internal server error")
	}
}

func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
