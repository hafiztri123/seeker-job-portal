package ports

import (
	"errors"

	"github.com/hafiztri123/internal/core/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type AuthService interface {
	Login(email, password string) (string, error)
	ValidateToken(token string) (*domain.User, error)
	Register(user *domain.User) error
}

type UserRepository interface{ 
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Create(user *domain.User) error
	Delete(id string) error
}

