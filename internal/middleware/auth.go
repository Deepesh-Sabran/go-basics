package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("VENGEANCE") // letter move to environment variable

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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			// Security check
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("❌ Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract claims and check if token is valid or not
		claims, ok:= token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("❌ Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// extract data from claims
		userId:= claims["user_id"]
		name:= claims["name"]
		role := claims["role"]
		permissions:= claims["permissions"]

		// add data to request context
		ctx:= context.WithValue(r.Context(), "user_id", userId)
		ctx = context.WithValue(ctx, "name", name)
		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "permissions", permissions)

		// If everything is fine → go to handler {{ passing updated value with context }}
		next(w, r.WithContext(ctx))
	}
}