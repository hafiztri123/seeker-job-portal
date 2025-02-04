package domain

import "github.com/google/uuid"

type LocationType string

const (
	remote LocationType = "remote"
	onsite LocationType = "onsite"
	hybrid LocationType = "hybrid"
)

type Job struct {
	ID uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"company_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	SalaryMin int `json:"salary_min"`
	SalaryMax int `json:"salary_max"`
	LocationType LocationType `json:"location_type"`
	

}