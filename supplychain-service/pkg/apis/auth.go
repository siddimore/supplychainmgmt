package apis

import (
	"fmt"
	"net/http"
	"github.com/dgrijalva/jwt-go"

	"strings"
)

var jwtSecret = []byte("your-secret-key")

func AuthMiddleware(allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("In AuthMiddleWare.")

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		fmt.Println("TokenString", tokenString)
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		fmt.Println("TokenString", tokenString)
		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		fmt.Println("Token", token)

		claims, ok := token.Claims.(jwt.MapClaims)
		for index,claim := range claims {
			fmt.Println("Index: %d, Value: %s\n", index, claim)
		}
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userRole, ok := claims["role"].(string)
		fmt.Println("UserRole", userRole)
		if !ok || !containsRole(userRole, allowedRoles) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func containsRole(role string, roles []string) bool {
	for _, r := range roles {
		if role == r {
			return true
		}
	}
	return false
}
