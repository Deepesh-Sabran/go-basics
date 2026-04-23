package models

// request struct (for input)
type CreateUserRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

// model struct (DB)
type User struct {
	ID       int
	Name     string
	Age      int
	Password string
}

// response struct (for output)
type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}