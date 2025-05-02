package middleware

import (
	"backend-challenge/utils"
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
			token, err := utils.ValidateToken(parts[1], secret)
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(jwt.MapClaims)
			userId := claims["user_id"]
			ctx := context.WithValue(r.Context(), "user_id", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
