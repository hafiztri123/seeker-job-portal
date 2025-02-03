package ports

import (
	"errors"

	"github.com/hafiztri123/internal/core/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)



type AuthService interface {
	Login(email, password string) (*domain.TokenPair, error)
	ValidateToken(token string) (*domain.User, error)
	Register(user *domain.User) error
	RefreshToken(token string) (*domain.TokenPair, error)
	RevokeRefreshToken(token string) error 
}

type UserRepository interface{ 
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Create(user *domain.User) error
	Delete(id string) error
	Update(user *domain.User) error
}

