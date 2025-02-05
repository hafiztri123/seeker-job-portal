package ports

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/internal/core/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type AuthService interface {
	Login(email, password string) (*domain.TokenPair, error)
	ValidateToken(token string) (jwt.MapClaims, error)
	Register(user *domain.User) error
	RefreshToken(token string) (*domain.TokenPair, error)
	RevokeRefreshToken(token string) error
	GetUser(id string) (*domain.User, error)
	GetCompany(id string) (*domain.Company, error)
}

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Create(user *domain.User) error
	Delete(id string) error
	Update(user *domain.User) error
}
