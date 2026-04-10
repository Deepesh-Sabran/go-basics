package main

import (
	"net/http"

	"github.com/Deepesh-Sabran/go-basics/internal/handlers"
)

func main() {
	http.HandleFunc("/hello", handlers.HelloHandler)

	http.ListenAndServe(":8080", nil)
}
