package ports

import (
	"github.com/google/uuid"
	"github.com/hafiztri123/internal/core/domain"
)


type UpdateProfileRequest struct {
	Fullname *string `json:"full_name,omitempty" validate:"required"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,e164"`
	About *string `json:"about,omitempty" validate:"omitempty,max=500"`
	Location *map[string]interface{} `json:"location,omitempty" validate:"omitempty"`
	ProfilePicture *string `json:"profile_picture,omitempty" validate:"omitempty,uuid4"`
}

type UserService interface {
	UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*domain.User, error)
	GetProfile(userID uuid.UUID) (*domain.User, error)
}