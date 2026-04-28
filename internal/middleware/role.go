package middleware

import (
	"log"
	"net/http"
)

func RequirePermission(permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			perms, ok:= r.Context().Value("permissions").([]interface{})
			if !ok {
				log.Println("you are forbidden to do this action")
				http.Error(w, "User Forbidden", http.StatusForbidden)
				return
			}

			hasPermission:= false

			for _, p:= range perms {
				if p.(string) == permission {
					hasPermission = true
					break
				}
			} 
			
			if !hasPermission {
				log.Println("you are forbidden to do this action")
				http.Error(w, "User Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}