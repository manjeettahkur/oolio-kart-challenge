package middleware

import (
	"encoding/json"
	"net/http"
	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/pkg/logger"
)

const (
	APIKeyHeader = "api_key"
	ValidAPIKey  = "apitest"
)

type AuthMiddleware struct {
	logger *logger.Logger
}

func NewAuthMiddleware(logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
	}
}

func (m *AuthMiddleware) RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(APIKeyHeader)

		if apiKey == "" {
			m.handleAuthError(w, r, errors.ErrMissingAPIKey)
			return
		}

		if apiKey != ValidAPIKey {
			m.handleAuthError(w, r, errors.ErrInvalidAPIKey)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) handleAuthError(w http.ResponseWriter, r *http.Request, err error) {
	apiError := errors.MapErrorToAPIError(err)

	contextLogger := m.logger.WithContext(r.Context())
	contextLogger.Warn("Authentication failed",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_ip", r.RemoteAddr,
		"user_agent", r.Header.Get("User-Agent"),
		"status_code", apiError.Code,
		"error_type", apiError.Type,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiError.Code)

	errorResponse := errors.ErrorResponse{
		Error: apiError,
	}

	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		contextLogger.Error("Failed to encode auth error response", "encode_error", encodeErr.Error())
	}
}
