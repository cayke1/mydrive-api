package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/cayke1/mydrive-api/internal/config"
	"github.com/cayke1/mydrive-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const EmailKey contextKey = "email"
const UserIdKey contextKey = "id"

var publicRoutes = map[string]bool{
	"/health":          true,
	"/auth/login":      true,
	"/auth/register":   true,
	"/auth/logout":     true,
}

func isPublicRoute(path string) bool {
	return publicRoutes[path]
}

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st, err := r.Cookie("session_token")
		if err != nil {
			log.Printf("[AUTH] No session_token cookie found: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if st.Value == "" {
			log.Printf("[AUTH] session_token is empty")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(st.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Load().JwtSecret), nil
		})
		if err != nil || !token.Valid || claims.Email == "" {
			log.Printf("[AUTH] Invalid token: %v, Valid: %v, Email: %s", err, token.Valid, claims.Email)
			http.Error(w, "Invalid token "+err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("[AUTH] Token valid for user: %s", claims.Email)
		ctx := context.WithValue(r.Context(), EmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserIdKey, claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		Authorize(next).ServeHTTP(w, r)
	})
}
