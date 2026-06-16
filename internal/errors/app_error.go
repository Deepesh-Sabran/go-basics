package errors

import (
	"errors"
	"fmt"
)

// ErrorCode represents the application error code
type ErrorCode string

const (
	CodeBadRequest          ErrorCode = "BAD_REQUEST"
	CodeNotFound            ErrorCode = "NOT_FOUND"
	CodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	CodeForbidden           ErrorCode = "FORBIDDEN"
	CodeConflict            ErrorCode = "CONFLICT"
	CodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeValidationError     ErrorCode = "VALIDATION_ERROR"
	CodeUnprocessableEntity ErrorCode = "UNPROCESSABLE_ENTITY"
)

// AppError represents an application error with status code and error code
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status_code"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"` // wrapped error for internal use
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       CodeBadRequest,
		Message:    message,
		StatusCode: 400,
	}
}

// BadRequestWithErr creates a 400 Bad Request error with wrapped error
func BadRequestWithErr(message string, err error) *AppError {
	return &AppError{
		Code:       CodeBadRequest,
		Message:    message,
		StatusCode: 400,
		Err:        err,
	}
}

// NotFound creates a 404 Not Found error
func NotFound(message string) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    message,
		StatusCode: 404,
	}
}

// NotFoundWithErr creates a 404 Not Found error with wrapped error
func NotFoundWithErr(message string, err error) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    message,
		StatusCode: 404,
		Err:        err,
	}
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: 401,
	}
}

// UnauthorizedWithErr creates a 401 Unauthorized error with wrapped error
func UnauthorizedWithErr(message string, err error) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: 401,
		Err:        err,
	}
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: 403,
	}
}

// ForbiddenWithErr creates a 403 Forbidden error with wrapped error
func ForbiddenWithErr(message string, err error) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: 403,
		Err:        err,
	}
}

// Conflict creates a 409 Conflict error
func Conflict(message string) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: 409,
	}
}

// ConflictWithErr creates a 409 Conflict error with wrapped error
func ConflictWithErr(message string, err error) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: 409,
		Err:        err,
	}
}

// ValidationError creates a 422 Unprocessable Entity error for validation failures
func ValidationError(message string) *AppError {
	return &AppError{
		Code:       CodeValidationError,
		Message:    message,
		StatusCode: 422,
	}
}

// ValidationErrorWithDetails creates a validation error with field details
func ValidationErrorWithDetails(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       CodeValidationError,
		Message:    message,
		StatusCode: 422,
		Details:    details,
	}
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(message string) *AppError {
	return &AppError{
		Code:       CodeInternalServerError,
		Message:    message,
		StatusCode: 500,
	}
}

// InternalServerErrorWithErr creates a 500 Internal Server Error with wrapped error
func InternalServerErrorWithErr(message string, err error) *AppError {
	return &AppError{
		Code:       CodeInternalServerError,
		Message:    message,
		StatusCode: 500,
		Err:        err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converts an error to AppError if possible
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}