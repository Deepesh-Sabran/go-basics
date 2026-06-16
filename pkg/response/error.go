package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Deepesh-Sabran/go-basics/internal/errors"
)

// ErrorResponse represents the error response structure sent to clients
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HandleError handles errors and writes appropriate HTTP response
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "unknown error occurred",
		})
		return
	}

	// Try to cast to AppError
	appErr, ok := errors.AsAppError(err)
	if !ok {
		// Not an AppError, treat as internal server error
		log.Printf("❌ Unexpected error type: %T: %v", err, err)
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, ErrorResponse{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		})
		return
	}

	// Write the app error
	w.WriteHeader(appErr.StatusCode)
	writeJSON(w, ErrorResponse{
		Code:    string(appErr.Code),
		Message: appErr.Message,
		Details: appErr.Details,
	})
}

// HandleSuccess writes a successful JSON response
func HandleSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeJSON is a helper function to write JSON with proper headers
func writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}