package users

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, u *User, pw string) error
	Login(ctx context.Context, email, pw string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmailOrPhone(ctx context.Context, identifier string) (*User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SearchUsers(ctx context.Context, query string) ([]*User, error)
}

type userService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, u *User, pw string) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	// check email exists
	_, err := s.repo.GetByEmail(ctx, u.Email)
	if err == nil {
		return ErrEmailExists
	}

	u.PasswordHash = HashPassword(pw)
	u.Role = "member"
	u.SecurityLabel = "public" // Default MAC label
	u.AccountType = "basic"    // Default account type
	u.Attributes = map[string]any{}

	return s.repo.Create(ctx, u)
}

func (s *userService) Login(ctx context.Context, email, pw string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !VerifyPassword(u.PasswordHash, pw) {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) GetByEmailOrPhone(ctx context.Context, identifier string) (*User, error) {
	return s.repo.GetByEmailOrPhone(ctx, identifier)
}

func (s *userService) UpdateUser(ctx context.Context, u *User) error {
	return s.repo.Update(ctx, u)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) SearchUsers(ctx context.Context, query string) ([]*User, error) {
	return s.repo.Search(ctx, query)
}

