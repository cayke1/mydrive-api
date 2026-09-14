package auth

import "net/http"

func (ac *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", ac.RegisterUserHandler)
	mux.HandleFunc("POST /auth/login", ac.RegisterLoginHandler)
	mux.HandleFunc("POST /auth/logout", ac.RegisterLogoutHandler)
	mux.Handle("GET /auth", Authorize(ac.redis, http.HandlerFunc(ac.GetUsersHandler)))
}
