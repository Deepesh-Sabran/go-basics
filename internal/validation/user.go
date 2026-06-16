package validation

import (
	"strings"

	appErrors "github.com/Deepesh-Sabran/go-basics/internal/errors"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
)

func ValidateCreateUser(req models.CreateUserRequest) error {
	errors := make(map[string]interface{})

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "Name is required"
	}

	if req.Age <= 0 {
		errors["age"] = "Age must be greater than 0"
	}

	if len(req.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters"
	}

	if len(errors) > 0 {
		return appErrors.ValidationErrorWithDetails("validation failed", errors)
	}

	return nil
}

func ValidateLogin(req models.LoginRequest) error {
	errors := make(map[string]interface{})

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "Name is required"
	}

	if strings.TrimSpace(req.Password) == "" {
		errors["password"] = "Password is required"
	}

	if len(errors) > 0 {
		return appErrors.ValidationErrorWithDetails("validation failed", errors)
	}

	return nil
}

func ValidateRefresh(req models.RefreshRequest) error {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return appErrors.Unauthorized("Refresh Token is required")
	}

	return nil
}

func ValidateUpdateUser(req models.UpdateUserRequest) error {
	errors := make(map[string]interface{})

	if req.Name == nil && req.Age == nil && req.Password == nil {
		return appErrors.BadRequest("at least one field is required")
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		errors["name"] = "name cannot be empty"
	}

	if req.Age != nil && *req.Age <= 0 {
		errors["age"] = "age must be greater than 0"
	}

	if req.Password != nil && len(*req.Password) < 8 {
		errors["password"] = "password must be at least 8 characters"
	}

	if len(errors) > 0 {
		return appErrors.ValidationErrorWithDetails("validation failed", errors)
	}

	return nil
}