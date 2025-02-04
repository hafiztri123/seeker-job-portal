package ports

import (
	"github.com/google/uuid"
	"github.com/hafiztri123/internal/core/domain"
)


type UpdateProfileRequest struct {
	Fullname *string `json:"full_name,omitempty" validate:"required"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
	About *string `json:"about,omitempty" validate:"omitempty,max=500"`
}

type UserService interface {
	UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*domain.User, error)
	GetProfile(userID uuid.UUID) (*domain.User, error)
}