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
	RoleID	 int		// for db relation between users -> roles table
	Role 	 string		// for response / JWT use
}

// response struct (for output)
type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	RoleID int  `json:"role_id"`
}

// paginated response struct
type PaginatedUserResponse struct {
	Count 	int				`json:"count"`
	Data	[]UserResponse	`json:"data"`
	Page	int				`json:"page"`
	Limit	int				`json:"limit"`
}

type TokenClaims struct {
    UserId      int      `json:"user_id"`
    Name        string   `json:"name"`
    jwt.RegisteredClaims
}

type EmailJob struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Retries int	  `json:"retries"`
}