package middleware

import (
	"log"
	"net/http"
	"strconv"
)

func RequireOwnershipOrPermission(permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract userId from token
			tokenUserId, ok := r.Context().Value("user_id").(int)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Extract path id (Optional)
			idStr := r.PathValue("id")
			
			isOwner := false
			if idStr != "" {
				idInt, err := strconv.Atoi(idStr)
				if err == nil {
					// Ownership check only if ID exists in URL
					isOwner = (tokenUserId == idInt)
				}
			}

			// 3. Check permission
			perms, _ := r.Context().Value("permissions").([]string)
			hasPermission := false
			for _, p := range perms {
				if p == permission {
					hasPermission = true
					break
				}
			}

			// 4. Logic: Allow if user IS the owner OR has the admin permission
			if !isOwner && !hasPermission {
				log.Println("🚫 Forbidden: Not owner and lacks permission:", permission)
				http.Error(w, "🚫 Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}
