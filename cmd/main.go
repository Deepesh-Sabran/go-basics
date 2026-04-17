package main

import (
	"log"
	"net/http"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/Deepesh-Sabran/go-basics/internal/handlers"
)

func main() {
	config.ConnectDB()
	
	http.HandleFunc("/hello", handlers.HelloHandler)
	http.HandleFunc("POST /create-user", handlers.CreateUser)
	http.HandleFunc("GET /get-users", handlers.GetUsers)
	http.HandleFunc("GET /get-user/{id}", handlers.GetUserById)
	http.HandleFunc("DELETE /delete-users", handlers.DeleteAllUsers)
	http.HandleFunc("DELETE /delete-user/{id}", handlers.DeleteUserById)
	http.HandleFunc("PATCH /update-user/{id}", handlers.UpdateUser)

	log.Println("Server listens on port:8080")
	http.ListenAndServe(":8080", nil)
}
