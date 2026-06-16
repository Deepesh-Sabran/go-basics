package models

import "github.com/golang-jwt/jwt/v5"

// request struct (for input)
type CreateUserRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Age      *int    `json:"age,omitempty"`
	Password *string `json:"password,omitempty"`
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

func (r UpdateUserRequest) ToUpdatesMap() map[string]interface{} {
	updates:= make(map[string]interface{})

	if r.Name != nil {
		updates["name"] = *r.Name
	}

	if r.Age != nil {
		updates["age"] = *r.Age
	}

	if r.Password != nil {
		updates["password"] = *r.Password
	}

	return updates
}