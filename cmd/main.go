package main

import (
	"fmt"
	"net/http"
	"github.com/Deepesh-Sabran/go-basics/internal/handlers"
)

func main() {
	http.HandleFunc("/hello", handlers.HelloHandler)
	http.HandleFunc("/user", handlers.CreateUser)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Server listens on port:8080")
}
