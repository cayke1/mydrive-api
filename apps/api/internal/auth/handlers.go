package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthController struct {
	service *AuthService
	redis   *redis.Client
}

func NewAuthController(service *AuthService, redis *redis.Client) *AuthController {
	return &AuthController{service: service, redis: redis}
}

func (ac *AuthController) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := ac.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (ac *AuthController) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := ac.service.RegisterUser(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to register: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.SessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    user.CSRFToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (ac *AuthController) RegisterLoginHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := ac.service.Login(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to register: "+err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.SessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    user.CSRFToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.ID)
}

func (ac *AuthController) RegisterLogoutHandler(w http.ResponseWriter, r *http.Request) {
	st, err := r.Cookie("session_token")
	if err == nil && st.Value != "" {
		ac.redis.Del(r.Context(), "session:"+st.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Write([]byte("Session cookie cleared"))
}
