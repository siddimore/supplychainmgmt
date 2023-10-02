package apis

import (
	//"encoding/json"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key")

// AuthMiddleware is a middleware for JWT authentication and role-based authorization.
func AuthMiddleware(allowedRoles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok || !containsRole(userRole, allowedRoles) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func containsRole(role string, roles []string) bool {
	for _, r := range roles {
		if role == r {
			return true
		}
	}
	return false
}
