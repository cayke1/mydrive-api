package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]User, error) {
	query := `SELECT id, email, created_at, updated_at
			FROM users ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, email, password_hash, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	var user User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

