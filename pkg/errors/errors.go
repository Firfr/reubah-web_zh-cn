package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode defines the type for error codes
type ErrorCode string

// Error codes
const (
	ErrInvalidFormat     ErrorCode = "INVALID_FORMAT"
	ErrInvalidSize       ErrorCode = "INVALID_SIZE"
	ErrProcessingFailed  ErrorCode = "PROCESSING_FAILED"
	ErrInvalidMIME       ErrorCode = "INVALID_MIME"
	ErrOptimizationFailed ErrorCode = "OPTIMIZATION_FAILED"
	ErrResizeFailed      ErrorCode = "RESIZE_FAILED"
	ErrBackgroundRemoval ErrorCode = "BACKGROUND_REMOVAL_FAILED"
	ErrPDFConversionFailed ErrorCode = "PDF_CONVERSION_FAILED"
)

// AppError represents an application error
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	err     error    // Internal error (not exposed in JSON)
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// New creates a new AppError
func New(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		err:     err,
	}
}

// Response represents the standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *AppError   `json:"error,omitempty"`
}

// SendError sends an error response
func SendError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = New("INTERNAL_ERROR", "An unexpected error occurred", err)
	}

	code := getHTTPCode(appErr.Code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   appErr,
	})
}

func getHTTPCode(code ErrorCode) int {
	switch code {
	case ErrInvalidFormat, ErrInvalidMIME, ErrInvalidSize:
		return http.StatusBadRequest
	case ErrProcessingFailed, ErrOptimizationFailed, ErrResizeFailed, ErrBackgroundRemoval:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
} 