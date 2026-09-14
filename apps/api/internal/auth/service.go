package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cayke1/mydrive-api/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	repo  *UserRepository
	Redis *redis.Client
}

func NewAuthService(repo *UserRepository, redis *redis.Client) *AuthService {
	return &AuthService{repo: repo, Redis: redis}
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

	userId := uuid.New().String()

	sessionToken, err := utils.GenerateToken(input.Email, userId)
	if err != nil {
		return nil, err
	}
	csrfToken, err := utils.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &User{
		ID:           userId,
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

	redisKey := "session:" + sessionToken
	if err := s.Redis.Set(ctx, redisKey, userId, 24*time.Hour).Err(); err != nil {
		log.Printf("[AUTH] Failed to store session in Redis: %v", err)
		return nil, err
	}
	log.Printf("[AUTH] Session stored in Redis with key: %s, value: %s", redisKey, userId)

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

	sessionToken, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return nil, err
	}

	csrfToken, err := utils.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	redisKey := "session:" + sessionToken
	if err := s.Redis.Set(ctx, redisKey, user.ID, 24*time.Hour).Err(); err != nil {
		log.Printf("[AUTH] Failed to store session in Redis: %v", err)
		return nil, err
	}
	log.Printf("[AUTH] Session stored in Redis with key: %s, value: %s", redisKey, user.ID)

	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	user.UpdatedAt = time.Now()

	return user, nil
}
