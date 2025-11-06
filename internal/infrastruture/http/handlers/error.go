package handlers

import (
	"encoding/json"
	"net/http"

	"ooliokartchallenge/internal/domain/errors"
	"ooliokartchallenge/pkg/logger"
)


func HandleError(w http.ResponseWriter, r *http.Request, err error, log *logger.Logger) {
	apiError := errors.MapErrorToAPIError(err)

	contextLogger := log.WithContext(r.Context())
	contextLogger.Error("Request failed",
		"error", err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", apiError.Code,
		"error_type", apiError.Type,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiError.Code)

	errorResponse := errors.ErrorResponse{
		Error: apiError,
	}

	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		contextLogger.Error("Failed to encode error response", "encode_error", encodeErr.Error())
	}
}
