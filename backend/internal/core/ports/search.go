package ports

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hafiztri123/internal/core/domain"
)


type JobSearchQuery struct {
	Keyword string
	LocationType string
	Skills []string
	SalaryMin *int
	SalaryMax *int
	Page int
	PageSize int
}

type JobSearchResult struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	SalaryRange   string    `json:"salary_range"`
	LocationType  domain.LocationType    `json:"location_type"`
	Location      json.RawMessage `json:"location"`
	EmploymentType domain.EmploymentType `json:"employment_type"`
	Status        string    `json:"status"`
	Skills        []string  `json:"skills"`
	CreatedAt     time.Time `json:"created_at"`
 }

type SearchService interface {
	SearchJobs(ctx context.Context, query JobSearchQuery) ([]JobSearchResult, error)
	IndexJob(ctx context.Context, job *domain.Job) error
}