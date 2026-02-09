package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mati/go-ticket/internal/domain"
)

type UserRepository struct {
	Queries *Queries
}

func NewUserRepository(queries *Queries) *UserRepository {
	return &UserRepository{
		Queries: queries,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	params := CreateUserParams{
		ID:           pgtype.UUID{Bytes: user.ID(), Valid: true},
		Email:        user.Email(),
		PasswordHash: user.PasswordHash(),
		Role:         UserRole(user.Role()),
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	_, err := ur.Queries.CreateUser(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrUserEmailAlreadyExists
		}
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	return nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := ur.Queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email from database: %w", err)
	}

	return domain.NewUserFromPersistence(user.ID.Bytes, user.Email, user.PasswordHash, domain.UserRole(user.Role), user.CreatedAt.Time, user.UpdatedAt.Time)
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := ur.Queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID from database: %w", err)
	}

	return domain.NewUserFromPersistence(user.ID.Bytes, user.Email, user.PasswordHash, domain.UserRole(user.Role), user.CreatedAt.Time, user.UpdatedAt.Time)
}
