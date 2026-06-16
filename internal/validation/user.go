package validation

import (
	"strings"

	appErrors "github.com/Deepesh-Sabran/go-basics/internal/errors"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
)

func ValidateCreateUser(req models.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return appErrors.BadRequest("Name is required")
	}

	if req.Age < 0 {
		return appErrors.BadRequest("Age must be greater than 0")
	}

	if len(req.Password) < 8 {
		return appErrors.BadRequest("Password must be at least 8 characters")
	}

	return nil
}

func ValidateLogin(req models.LoginRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return appErrors.BadRequest("Name is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return appErrors.BadRequest("Password is required")
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
	if req.Name == nil && req.Age == nil && req.Password == nil {
		return appErrors.BadRequest("at least one field is required")
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return appErrors.BadRequest("name cannot be empty")
	}

	if req.Age != nil && *req.Age <= 0 {
		return appErrors.BadRequest("age must be greater than 0")
	}
	
	if req.Password != nil && len(*req.Password) < 8 {
		return appErrors.BadRequest("password must be at least 8 characters")
	}

	return nil
}