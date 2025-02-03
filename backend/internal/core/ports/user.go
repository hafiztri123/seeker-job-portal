package ports

import (
	"github.com/google/uuid"
	"github.com/hafiztri123/internal/core/domain"
)


type UpdateProfileRequest struct {
	Fullname string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	About string `json:"about"`
}

type UserService interface {
	UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*domain.User, error)
	GetProfile(userID uuid.UUID) (*domain.User, error)
}