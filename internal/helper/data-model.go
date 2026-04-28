package helper

import "github.com/Deepesh-Sabran/go-basics/internal/models"

func ToUserResponse(user models.User) models.UserResponse {
	return models.UserResponse{
		ID: 	user.ID,
		Name:	user.Name,
		Age: 	user.Age,
		Role:	user.Role,
	}
}