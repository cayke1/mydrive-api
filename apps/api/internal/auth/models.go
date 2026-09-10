package auth

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	SessionToken string    `json:"session_token"`
	CSRFToken    string    `json:"csrf_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	ID           string    `json:"id"`
	SessionToken string    `json:"session_token"`
	CSRFToken    string    `json:"csrf_token"`
	UpdatedAt    time.Time `json:"updated_at"`
}
