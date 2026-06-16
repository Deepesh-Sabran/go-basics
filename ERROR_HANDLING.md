# Application Typed Error Handling

This document explains the production-grade typed error handling system implemented in this application.

## Overview

The error handling system provides:

- **Type-safe errors** with structured error codes and HTTP status codes
- **Error wrapping** for maintaining error context throughout the stack
- **Validation errors** with field-level details
- **Consistent HTTP responses** with standardized error format
- **Error checking utilities** for safe error type assertions

## Architecture

### Core Components

#### 1. **AppError** (`internal/errors/app_error.go`)

The central error type that implements the `error` interface:

```go
type AppError struct {
    Code       ErrorCode              // Machine-readable error code
    Message    string                 // User-friendly error message
    StatusCode int                    // HTTP status code
    Details    map[string]interface{} // Additional error details
    Err        error                  // Wrapped error for debugging
}
```

**Features:**

- Implements `error` interface with `Error()` method
- Implements `Unwrap()` for error chain inspection (Go 1.13+)
- `WithDetails()` method to add field-specific information
- Structured JSON serialization for API responses

#### 2. **Error Codes** (`ErrorCode` type)

Pre-defined error codes for consistency:

```go
CodeBadRequest          = "BAD_REQUEST"              // 400
CodeNotFound            = "NOT_FOUND"                // 404
CodeUnauthorized        = "UNAUTHORIZED"             // 401
CodeForbidden           = "FORBIDDEN"                // 403
CodeConflict            = "CONFLICT"                 // 409
CodeValidationError     = "VALIDATION_ERROR"         // 422
CodeInternalServerError = "INTERNAL_SERVER_ERROR"    // 500
```

#### 3. **Error Constructors**

Two variations for each error type:

```go
// Basic error (no wrapped error)
NotFound(message string) *AppError

// With wrapped error for debugging
NotFoundWithErr(message string, err error) *AppError
```

#### 4. **Response Handler** (`pkg/response/error.go`)

Converts AppErrors to consistent HTTP responses:

```go
HandleError(w http.ResponseWriter, err error)
HandleSuccess(w http.ResponseWriter, statusCode int, data interface{})
```

## Usage Examples

### 1. Basic Error Creation

```go
// Simple error
return appErrors.NotFound("User not found")

// With wrapped internal error
return appErrors.NotFoundWithErr("User not found", dbErr)
```

### 2. Validation Errors with Details

```go
errors := make(map[string]interface{})
errors["email"] = "Invalid email format"
errors["age"] = "Age must be greater than 18"

return appErrors.ValidationErrorWithDetails("validation failed", errors)
```

### 3. Adding Context with Details

```go
err := someOperation()
if err != nil {
    return appErrors.ConflictWithErr("username already taken", err).WithDetails(map[string]interface{}{
        "field": "username",
        "value": "john_doe",
    })
}
```

### 4. In Handlers

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := services.GetUserById(id)
    if err != nil {
        // AppError is automatically handled
        response.HandleError(w, err)
        return
    }

    response.HandleSuccess(w, http.StatusOK, user)
}
```

### 5. In Services

```go
func GetUserById(userId int) (*User, error) {
    user, err := repo.GetUserById(userId)
    if err != nil {
        log.Println("DB error:", err)
        // Wrap DB error with application error
        return nil, appErrors.NotFoundWithErr("user not found", err)
    }

    if user == nil {
        return nil, appErrors.NotFound("User not found")
    }

    return user, nil
}
```

## Error Checking Utilities

### Type Assertion with Safety

```go
// Method 1: Using errors.As (recommended)
appErr, ok := errors.AsAppError(err)
if ok {
    statusCode := appErr.StatusCode
    code := appErr.Code
}

// Method 2: Check without assertion
if errors.IsAppError(err) {
    // Handle AppError specifically
}
```

### Error Unwrapping

```go
// Get wrapped error
if err != nil {
    wrappedErr := errors.Unwrap(err)
    if wrappedErr != nil {
        // Handle underlying error
    }
}
```

## Response Format

### Success Response

```json
{
    "data": {...},
    "status": 200
}
```

### Error Response

```json
{
  "code": "VALIDATION_ERROR",
  "message": "validation failed",
  "details": {
    "email": "Invalid email format",
    "age": "Age must be greater than 18"
  }
}
```

### Validation Error with Details

```json
{
  "code": "CONFLICT",
  "message": "username already taken",
  "details": {
    "field": "username",
    "value": "john_doe"
  }
}
```

## Best Practices

### 1. **Use Typed Errors Consistently**

- Always return `*AppError` from business logic
- Never return `errors.New()` from services

```go
// ❌ Don't do this
if err != nil {
    return errors.New("something went wrong")
}

// ✅ Do this
if err != nil {
    return appErrors.InternalServerError("something went wrong")
}
```

### 2. **Wrap External Errors**

- Preserve error context by wrapping
- Use `*WithErr()` constructors for debugging

```go
// ❌ Don't lose error context
dbErr := db.Query()
if dbErr != nil {
    return appErrors.InternalServerError("database error")
}

// ✅ Wrap for debugging
dbErr := db.Query()
if dbErr != nil {
    return appErrors.InternalServerErrorWithErr("database error", dbErr)
}
```

### 3. **Use Validation Errors for Input Validation**

- Collect all validation errors
- Return them with field details

```go
errors := make(map[string]interface{})

if len(req.Email) == 0 {
    errors["email"] = "Email is required"
}

if req.Age < 18 {
    errors["age"] = "Age must be at least 18"
}

if len(errors) > 0 {
    return appErrors.ValidationErrorWithDetails("validation failed", errors)
}
```

### 4. **Consistent Error Handling in Handlers**

- Use `response.HandleError()` for error responses
- Use `response.HandleSuccess()` for success responses

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := services.CreateUser(req)
    if err != nil {
        response.HandleError(w, err)
        return
    }

    response.HandleSuccess(w, http.StatusCreated, user)
}
```

### 5. **Log Internal Errors**

- Log wrapped errors for debugging
- Don't expose internal error details to clients

```go
user, err := services.GetUserById(id)
if err != nil {
    log.Printf("Failed to fetch user: %v", err) // Log with details
    response.HandleError(w, err)                 // Send user-friendly error
    return
}
```

## Migration Guide

### From Old Pattern

```go
// Old
http.Error(w, "invalid input", http.StatusBadRequest)

// New
response.HandleError(w, appErrors.BadRequest("invalid input"))
```

### From errors.New()

```go
// Old
return errors.New("user not found")

// New
return appErrors.NotFound("user not found")
```

### From String Matching

```go
// Old
if err.Error() == "user not found" {
    http.Error(w, "User not found", http.StatusNotFound)
}

// New
appErr, ok := errors.AsAppError(err)
if ok && appErr.Code == errors.CodeNotFound {
    response.HandleError(w, err)
}
```

## Testing

### Asserting Error Types

```go
func TestGetUserNotFound(t *testing.T) {
    user, err := services.GetUserById(999)

    appErr, ok := errors.AsAppError(err)
    if !ok {
        t.Fatal("expected AppError")
    }

    if appErr.Code != errors.CodeNotFound {
        t.Errorf("expected NOT_FOUND, got %s", appErr.Code)
    }
}
```

### Testing with Details

```go
func TestValidationError(t *testing.T) {
    err := validation.ValidateCreateUser(invalidReq)

    appErr, ok := errors.AsAppError(err)
    if !ok {
        t.Fatal("expected AppError")
    }

    if appErr.Code != errors.CodeValidationError {
        t.Fatal("expected validation error")
    }

    if details, ok := appErr.Details["email"]; !ok {
        t.Error("expected email in details")
    }
}
```

## Summary

This typed error system provides:
✅ Type-safe error handling
✅ Consistent error responses
✅ Error context preservation through wrapping
✅ Field-level validation error details
✅ Production-ready implementation
✅ Easy integration with existing code
