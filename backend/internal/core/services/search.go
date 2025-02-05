package services

import (
	"context"
	"fmt"

	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
	"github.com/hafiztri123/internal/repositories/elasticsearch"
	"github.com/hafiztri123/internal/repositories/postgres"
)

type searchService struct {
	searchRepo elasticsearch.SearchRepository
	jobRepo   postgres.JobRepository
 }
 
 func NewSearchService(searchRepo elasticsearch.SearchRepository, jobRepo postgres.JobRepository) ports.SearchService {
	return &searchService{
		searchRepo: searchRepo,
		jobRepo:    jobRepo,
	}
 }
 
 func (s *searchService) SearchJobs(ctx context.Context, query ports.JobSearchQuery) ([]ports.JobSearchResult, error) {
	esQuery := elasticsearch.JobSearchQuery{
		Keyword:      query.Keyword,
		LocationType: query.LocationType,
		Skills:      query.Skills,
		SalaryMin:   query.SalaryMin,
		SalaryMax:   query.SalaryMax,
		Page:        query.Page,
		PageSize:    query.PageSize,
	}
 
	jobs, err := s.searchRepo.SearchJobs(ctx, esQuery)
	if err != nil {
		return nil, err
	}
 
	results := make([]ports.JobSearchResult, len(jobs))
	for i, job := range jobs {
		results[i] = mapToSearchResult(job)
	}
 
	return results, nil
 }
 
 func (s *searchService) IndexJob(ctx context.Context, job *domain.Job) error {
	doc := mapToJobDocument(job)
	return s.searchRepo.IndexJob(ctx, doc)
 }
 
 func mapToJobDocument(job *domain.Job) *elasticsearch.JobDocument {
	return &elasticsearch.JobDocument{
		ID:            job.ID.String(),
		CompanyID:     job.CompanyID.String(),
		Title:         job.Title,
		Description:   job.Description,
		SalaryMin:     job.SalaryMin,
		SalaryMax:     job.SalaryMax,
		LocationType:  job.LocationType,
		Location:      job.Location,
		Status:        job.Status,
		CreatedAt:     job.CreatedAt,
	}
 }

 func mapToSearchResult(doc elasticsearch.JobDocument) ports.JobSearchResult {
	return ports.JobSearchResult{
		ID:           doc.ID,
		CompanyID:    doc.CompanyID,
		Title:        doc.Title,
		Description:  doc.Description,
		SalaryRange:  fmt.Sprintf("$%d - $%d", doc.SalaryMin, doc.SalaryMax),
		LocationType: doc.LocationType,
		Skills:       doc.Skills,
		CreatedAt:    doc.CreatedAt,
	}
 }
 