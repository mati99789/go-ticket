package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleUser      UserRole = "user"
	UserRoleAdmin     UserRole = "admin"
	UserRoleOrganizer UserRole = "organizer"
)

type User struct {
	id           uuid.UUID
	email        string
	passwordHash string
	role         UserRole
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(id uuid.UUID, email string, password string, role UserRole) (*User, error) {
	if id == uuid.Nil {
		return nil, ErrUserIDNil
	}

	if email == "" {
		return nil, ErrUserEmailEmpty
	}
	if password == "" {
		return nil, ErrUserPasswordEmpty
	}
	if role != UserRoleUser && role != UserRoleAdmin && role != UserRoleOrganizer {
		return nil, ErrUserRoleInvalid
	}

	if !strings.Contains(email, "@") {
		return nil, ErrUserEmailInvalid
	}

	return &User{
		id:           id,
		email:        email,
		passwordHash: password,
		role:         role,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}, nil
}

func (u *User) UpdateEmail(email string) error {
	if email == "" {
		return ErrUserEmailEmpty
	}

	if !strings.Contains(email, "@") {
		return ErrUserEmailInvalid
	}
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

func (u *User) UpdatePassword(password string) error {
	if password == "" {
		return ErrUserPasswordEmpty
	}

	u.passwordHash = password
	u.updatedAt = time.Now()
	return nil
}

func (u *User) UpdateRole(role UserRole) error {
	if role != UserRoleUser && role != UserRoleAdmin && role != UserRoleOrganizer {
		return ErrUserRoleInvalid
	}
	u.role = role
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) Role() UserRole {
	return u.role
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
