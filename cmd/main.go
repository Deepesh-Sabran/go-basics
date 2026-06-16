package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/Deepesh-Sabran/go-basics/internal/handlers"
	"github.com/Deepesh-Sabran/go-basics/internal/middleware"
	"github.com/Deepesh-Sabran/go-basics/internal/workers"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()

	ctx, stop:= signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	http.HandleFunc("/hello", handlers.HelloHandler)
	http.HandleFunc("POST /login", handlers.Login)
	http.HandleFunc("POST /refresh", handlers.Refresh)
	http.HandleFunc("POST /signup", handlers.CreateUser)
	http.HandleFunc("POST /logout", handlers.Logout)
	http.HandleFunc("GET /get-users", middleware.AuthMiddleWare(handlers.GetUsers))
	http.HandleFunc("GET /get-user/{id}", middleware.AuthMiddleWare(middleware.RequireOwnershipOrPermission("view_user")(handlers.GetUserById)))
	http.HandleFunc("GET /me", middleware.AuthMiddleWare(handlers.GetMe))
	http.HandleFunc("DELETE /delete-users", middleware.AuthMiddleWare(middleware.RequireOwnershipOrPermission("delete_user")(handlers.DeleteAllUsers)))
	http.HandleFunc("DELETE /delete-user/{id}", middleware.AuthMiddleWare(middleware.RequireOwnershipOrPermission("delete_user")(handlers.DeleteUserById)))
	http.HandleFunc("PATCH /update-user/{id}", middleware.AuthMiddleWare(middleware.RequireOwnershipOrPermission("update_user")(handlers.UpdateUser)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	workers.RecoverEmailJobs()
	go workers.StartEmailWorker(ctx)

	go func() {
		<-ctx.Done()

		log.Println("shutting down...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Println("server shutdown error:", err)
		}
	}()

	log.Println("Server listens on port:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
