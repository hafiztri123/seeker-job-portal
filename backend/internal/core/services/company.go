package services

import (
	"github.com/google/uuid"
	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
	"github.com/hafiztri123/internal/repositories/postgres"
)

type CompanyService struct {
	repo *postgres.CompanyRepository
}

func NewCompanyService(repo *postgres.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}


func (s *CompanyService) Register(req *ports.CompanyRegisterRequest) error {
	company, err := domain.NewCompany(req.Name, req.Email, req.Password)
	if err != nil {
		return err
	}

	if err := s.repo.Create(company); err != nil {
		return err
	}

	return nil
}

func (s *CompanyService) UpdateProfile(id string, req ports.CompanyUpdateProfileRequest) (*domain.Company, error) {
	company, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.About != nil {
		company.About = req.About
	}

	if req.BussinessCategory != nil {
		company.BussinessCategory = req.BussinessCategory
	}

	if req.CompanySize != nil {
		company.Size = req.CompanySize
	}

	if req.PhoneNumber != nil {
		company.PhoneNumber = req.PhoneNumber
	}

	if req.Location != nil {
		company.Location = req.Location
	}

	if req.Name != nil {
		company.Name = *req.Name
	}


	if err := s.repo.Update(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyService) GetAllCompany() ([]*domain.Company, error) {
	return s.repo.GetAllCompany()
}

func(s *CompanyService) GetCompany(userID uuid.UUID) (*domain.Company, error) {
	return s.repo.FindByID(userID.String())
}

