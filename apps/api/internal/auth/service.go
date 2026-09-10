package auth

import (
	"context"
	"errors"
	"time"

	"github.com/cayke1/mydrive-api/internal/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	repo *UserRepository
}

func NewAuthService(repo *UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GetUsers(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *AuthService) RegisterUser(ctx context.Context, input CreateUserInput) (*User, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if len(input.Email) > 100 {
		return nil, errors.New("user email exceeds maximum length of 100")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	userExists, _ := s.repo.GetByEmail(ctx, input.Email)

	if userExists != nil {
		return nil, errors.New("email already in use")
	}
	passwordHash, err := utils.HashPassword(input.Password)

	if err != nil {
		return nil, err
	}

	sessionToken, err := utils.GenerateToken(input.Email)
	if err != nil {
		return nil, err
	}
	csrfToken, err := utils.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &User{
		ID:           uuid.New().String(),
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = s.repo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input CreateUserInput) (*User, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if len(input.Email) > 100 {
		return nil, errors.New("user email exceeds maximum length of 100")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	sessionToken, err := utils.GenerateToken(user.Email)
	if err != nil {
		return nil, err
	}

	csrfToken, err := utils.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	query := &UpdateUserInput{
		ID:           user.ID,
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		UpdatedAt:    time.Now(),
	}

	updatedUser, err := s.repo.UpdateSession(ctx, query)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
