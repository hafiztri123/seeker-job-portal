package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)



type LocationType string
type EmploymentType string

var (
	ErrJobNotFound = errors.New("job not found")
)

const (
	remote LocationType = "remote"
	onsite LocationType = "onsite"
	hybrid LocationType = "hybrid"
)

const (
	// [full_time, part_time, internship, contract]
	fullTime EmploymentType = "full_time"
	partTime EmploymentType = "part_time"
	internship EmploymentType = "internship"
	contract EmploymentType = "contract"
)

type Job struct {
	ID uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"company_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	SalaryMin int `json:"salary_min"`
	SalaryMax int `json:"salary_max"`
	LocationType LocationType `json:"location_type"`
	Location json.RawMessage `json:"location"`
	EmploymentType EmploymentType `json:"employment_type"`
	Deadline time.Time `json:"deadline"`
	Status string `json:"status"`
	Skills []string `json:"skills"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}