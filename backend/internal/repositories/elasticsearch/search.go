package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/hafiztri123/config"
	"github.com/hafiztri123/internal/core/domain"
)

type SearchRepository struct{
	client *elasticsearch.Client
	index string
}

func NewSearchRepository(cfg *config.ElasticSearchConfig) (*SearchRepository, error) {
	client, err :=elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	})

	return &SearchRepository{
		client: client,
		index: cfg.Index,
	}, err
}


type JobDocument struct {
	ID            string    `json:"id"`
	CompanyID     string    `json:"company_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	SalaryMin     int       `json:"salary_min"`
	SalaryMax     int       `json:"salary_max"`
	LocationType  domain.LocationType    `json:"location_type"`
	Location      json.RawMessage `json:"location"`
	EmploymentType domain.EmploymentType    `json:"employment_type"`
	Status        string    `json:"status"`
	Skills        []string  `json:"skills"`
	CreatedAt     time.Time    `json:"created_at"`
 }

type JobSearchQuery struct {
	Keyword string
	LocationType string
	Skills []string
	SalaryMin *int
	SalaryMax *int
	Page int
	PageSize int
}

 func (r *SearchRepository) IndexJob(ctx context.Context, job *JobDocument) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = r.client.Index(
		r.index,
		bytes.NewReader(payload),
		r.client.Index.WithDocumentID(job.ID),
		r.client.Index.WithContext(ctx),
	)
	return err
 }

 func buildFilters(query JobSearchQuery) []map[string]interface{} {
	var filters []map[string]interface{}

	if query.LocationType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"location_type": query.LocationType,
			},
		})
	}

	if len(query.Skills) > 0 {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"skills": query.Skills,
			},
		})
	}

	if query.SalaryMin != nil {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"salary_min": map[string]interface{}{
					"gte": *query.SalaryMin,
				},
			},
		})
	}

	if query.SalaryMax != nil {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"salary_max": map[string]interface{}{
					"lte": *query.SalaryMax,
				},
			},
		})
	}

	return filters
 }

 func (r *SearchRepository) SearchJobs(ctx context.Context, query JobSearchQuery) ([]JobDocument, error) {
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  query.Keyword,
							"fields": []string{"title^2", "description", "skills"},
						},
					},
				},
				"filter": buildFilters(query),
			},
		},
		"sort": []map[string]interface{}{
			{"created_at": "desc"},
		},
		"from": (query.Page - 1) * query.PageSize,
		"size": query.PageSize,
	}
 
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}
 
	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(r.index),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
 
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
 
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	jobs := make([]JobDocument, len(hits))
 
	for i, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		
		if err := json.Unmarshal(sourceBytes, &jobs[i]); err != nil {
			return nil, err
		}
	}
 
	return jobs, nil
 }