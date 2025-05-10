package service

import (
	"context"
	"errors"

	"github.com/DexScen/WebBook/backend/auth/internal/domain"
	e "github.com/DexScen/WebBook/backend/auth/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	Register(ctx context.Context, user *domain.User) error
	GetPassword(ctx context.Context, login string) (string, error)
	GetRole(ctx context.Context, login string) (string, error)
	UserExists(ctx context.Context, login, email string) (bool, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
}

type Users struct {
	repo UsersRepository
}

func NewUsers(repo UsersRepository) *Users {
	return &Users{
		repo: repo,
	}
}

func (u *Users) LogIn(ctx context.Context, login, password string) (string, error) {
	//todo check if can login
	passwordHash, err := u.repo.GetPassword(ctx, login)

	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			return "", e.ErrUserNotFound
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", e.ErrWrongPassword
	}
	return u.repo.GetRole(ctx, login)
}

func (u *Users) Register(ctx context.Context, user *domain.User) error {
	exists, err := u.repo.UserExists(ctx, user.Login, user.Email)
	if exists {
		return e.ErrUserExists
	}
	if err != nil {
		return err
	}

	//костыль, обычно делают мапперы структуры на структуру поискать в инете
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return u.repo.Register(ctx, &domain.User{
		Login:    user.Login,
		Name:     user.Name,
		Password: string(hash),
		Email:    user.Email,
		Role:     user.Role,
	})
	// check unique login + check ajax l8r
}

// GetByLogin returns a user by login
func (u *Users) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	// Implement logic to get user by login
	// This requires a new repository method
	return u.repo.GetByLogin(ctx, login)
}
