package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mati/go-ticket/internal/auth"
	"github.com/mati/go-ticket/internal/domain"
)

type UserServiceInterface interface {
	RegisterUser(ctx context.Context, email, password string) error
	LoginUser(ctx context.Context, email, password string) (string, error)
}

type UserService struct {
	userRepository domain.UserRepository
	jwtService     *auth.JWTService
}

func NewUserService(userRepository domain.UserRepository, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, email, password string) error {
	if err := domain.ValidatePassword(password); err != nil {
		return err
	}
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	newUser, err := domain.NewUser(uuid.New(), email, passwordHash, domain.UserRoleUser)
	if err != nil {
		return err
	}
	return s.userRepository.CreateUser(ctx, newUser)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	userFromDB, err := s.userRepository.GetUserByEmail(ctx, email)

	hashToVerify := "$2a$10$dummyhashfortimingattackprotection1234567890123456"
	if err == nil {
		hashToVerify = userFromDB.PasswordHash()
	}

	verifyErr := auth.VerifyPassword(hashToVerify, password)

	if err != nil || verifyErr != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.jwtService.GenerateToken(userFromDB)
	if err != nil {
		return "", err
	}
	return token, nil
}
