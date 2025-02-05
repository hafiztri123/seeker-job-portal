package services

import (
	"github.com/google/uuid"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
)

type profileService struct {
	userRepo ports.UserRepository
}

func NewProfileService(userRepo ports.UserRepository) ports.UserService {
	return &profileService{
		userRepo: userRepo,
	}
}
func(s *profileService) UpdateProfile(userID uuid.UUID, req ports.UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID.String())
	if err != nil {
		return nil, err
	}

	if req.Fullname != nil {
		user.FullName = *req.Fullname
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}

	if req.About != nil {
		user.About = *req.About
	}

	if req.Location != nil {
		user.Location = *req.Location
	}

	if req.ProfilePicture != nil {
		user.ProfilePicture = *req.ProfilePicture
	}



	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *profileService) GetProfile(userID uuid.UUID) (*domain.User, error) {
	return s.userRepo.FindByID(userID.String())
}