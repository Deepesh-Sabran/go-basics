package main

import (
	"log"
	"net/http"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/Deepesh-Sabran/go-basics/internal/handlers"
	"github.com/Deepesh-Sabran/go-basics/internal/middleware"
)

func main() {
	config.ConnectDB()
	
	http.HandleFunc("/hello", handlers.HelloHandler)
	http.HandleFunc("POST /login", handlers.Login)
	http.HandleFunc("POST /refresh", handlers.Refresh)
	http.HandleFunc("POST /signup", handlers.CreateUser)
	http.HandleFunc("GET /get-users", middleware.AuthMiddleWare(middleware.RequirePermission("view_user")(handlers.GetUsers)))
	http.HandleFunc("GET /get-user/{id}", middleware.AuthMiddleWare(middleware.RequirePermission("view_user")(handlers.GetUserById)))
	http.HandleFunc("GET /me", middleware.AuthMiddleWare(handlers.GetMe))
	http.HandleFunc("DELETE /delete-users", middleware.AuthMiddleWare(middleware.RequirePermission("delete_user")(handlers.DeleteAllUsers)))
	http.HandleFunc("DELETE /delete-user/{id}", middleware.AuthMiddleWare(middleware.RequirePermission("delete_user")(handlers.DeleteUserById)))
	http.HandleFunc("PATCH /update-user/{id}", middleware.AuthMiddleWare(handlers.UpdateUser))

	log.Println("Server listens on port:8080")
	http.ListenAndServe(":8080", nil)
}
