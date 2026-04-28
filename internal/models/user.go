package models

import "github.com/golang-jwt/jwt/v5"

// request struct (for input)
type CreateUserRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Role 	 string	`json:"role"`
}

// model struct (DB)
type User struct {
	ID       int
	Name     string
	Age      int
	Password string
	Role 	 string `json:"role"`
}

// response struct (for output)
type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Role string	`json:"role"`
}

// paginated response struct
type PaginatedUserResponse struct {
	Count 	int				`json:"count"`
	Data	[]UserResponse	`json:"data"`
	Page	int				`json:"page"`
	Limit	int				`json:"limit"`
}

type AccessClaims struct {
    UserId      int     `json:"user_id"`
    Name        string   `json:"name"`
    Role        string   `json:"role"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}