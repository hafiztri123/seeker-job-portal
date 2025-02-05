package ports

import (
	"encoding/json"

	"github.com/hafiztri123/internal/core/domain"
)


type CompanyRegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type CompanyUpdateProfileRequest struct {
	Name *string `json:"name,omitempty"`
	About *string `json:"about,omitempty"`
	BussinessCategory *string `json:"bussiness_category,omitempty"`
	CompanySize *string `json:"company_size,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Location *json.RawMessage `json:"location,omitempty"`
}


type CompanyService interface {
	Register(req *CompanyRegisterRequest) error
	UpdateProfile(id string, req CompanyUpdateProfileRequest) (*domain.Company, error)
	GetAllCompany() ([]*domain.Company, error)
	GetCompany() (*domain.Company, error)
}