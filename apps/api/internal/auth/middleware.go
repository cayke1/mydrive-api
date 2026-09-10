package auth

import (
	"context"
	"net/http"

	"github.com/cayke1/mydrive-api/internal/config"
	"github.com/cayke1/mydrive-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const EmailKey contextKey = "email"

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if st.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(st.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Load().JwtSecret), nil
		})
		if err != nil || !token.Valid || claims.Email == "" {
			http.Error(w, "Invalid token "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), EmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
