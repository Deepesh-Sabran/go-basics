package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Deepesh-Sabran/go-basics/internal/config"
	"github.com/Deepesh-Sabran/go-basics/internal/models"
	repo "github.com/Deepesh-Sabran/go-basics/internal/repository"
	"github.com/Deepesh-Sabran/go-basics/internal/services"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("❌ Token is missing")
			http.Error(w, "Missing Token", http.StatusUnauthorized)
			return
		}

		// Validate format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("❌ Invalid token format")
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse & validate token
		token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {

			// Security check
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return config.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			log.Println("❌ Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract claims and check if token is valid or not
		claims, ok := token.Claims.(*models.TokenClaims)
		if !ok {
			log.Println("❌ Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := repo.GetUserAuthInfo(claims.UserId)
		if err != nil {
			log.Println("❌ User not found 😞 --Auth")
			http.Error(w, "User not found 😞 --Auth", http.StatusBadRequest)
			return
		}

		permissions, err := services.GetPermissionsByRole(user.RoleID)
		if err != nil {
			log.Println("❌ No Permission found")
			http.Error(w, "No Permission Found", http.StatusBadRequest)
			return
		}

		// add data to request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		ctx = context.WithValue(ctx, "name", claims.Name)
		ctx = context.WithValue(ctx, "permissions", permissions)

		// If everything is fine → go to handler {{ passing updated value with context }}
		next(w, r.WithContext(ctx))
	}
}