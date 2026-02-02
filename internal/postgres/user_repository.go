package postgres

import (
	"context"

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
	return nil
}
