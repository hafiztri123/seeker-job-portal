// internal/repositories/postgres/job.go
package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hafiztri123/internal/core/domain"
)

type JobRepository struct {
   db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
   return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *domain.Job) error {
   query := `
       INSERT INTO jobs (
           id, company_id, title, description, 
           salary_min, salary_max, location_type,
           location, employment_type, deadline, status
       ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
   `
   locationJSON, err := json.Marshal(job.Location)
   if err != nil {
       return err
   }

   _, err = r.db.Exec(
       query,
       job.ID,
       job.CompanyID,
       job.Title,
       job.Description,
       job.SalaryMin,
       job.SalaryMax,
       job.LocationType,
       locationJSON,
       job.EmploymentType,
       job.Deadline,
       job.Status,
   )
   return err
}

func (r *JobRepository) GetByID(id uuid.UUID) (*domain.Job, error) {
   query := `
       SELECT 
           id, company_id, title, description,
           salary_min, salary_max, location_type,
           location, employment_type, deadline, status,
           created_at, updated_at
       FROM jobs
       WHERE id = $1 AND deleted_at IS NULL
   `

   job := &domain.Job{}
   var locationJSON []byte

   err := r.db.QueryRow(query, id).Scan(
       &job.ID,
       &job.CompanyID,
       &job.Title,
       &job.Description,
       &job.SalaryMin,
       &job.SalaryMax,
       &job.LocationType,
       &locationJSON,
       &job.EmploymentType,
       &job.Deadline,
       &job.Status,
       &job.CreatedAt,
       &job.UpdatedAt,
   )

   if err == sql.ErrNoRows {
       return nil, err
   }
   if err != nil {
       return nil, err
   }

   if err := json.Unmarshal(locationJSON, &job.Location); err != nil {
       return nil, err
   }

   return job, nil
}


func (r *JobRepository) Update(job *domain.Job) error {
	query := `
		UPDATE jobs 
		SET title = $1, description = $2, 
			salary_min = $3, salary_max = $4,
			location_type = $5, location = $6,
			employment_type = $7, deadline = $8,
			status = $9
		WHERE id = $10 AND deleted_at IS NULL
	`
	locationJSON, err := json.Marshal(job.Location)
	if err != nil {
		return err
	}
 
	result, err := r.db.Exec(
		query,
		job.Title,
		job.Description,
		job.SalaryMin,
		job.SalaryMax,
		job.LocationType,
		locationJSON,
		job.EmploymentType,
		job.Deadline,
		job.Status,
		job.ID,
	)
	if err != nil {
		return err
	}
 
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrJobNotFound
	}
	return nil
 }
 
 func (r *JobRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE jobs 
		SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
 
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrJobNotFound
	}
	return nil
 }
 
 func (r *JobRepository) List(filter JobFilter) ([]*domain.Job, error) {
	query := `
		SELECT 
			id, company_id, title, description,
			salary_min, salary_max, location_type,
			location, employment_type, deadline, status,
			created_at, updated_at
		FROM jobs
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	if filter.CompanyID != uuid.Nil {
		query += " AND company_id = $1"
		args = append(args, filter.CompanyID)
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, filter.Status)
	}
	query += " ORDER BY created_at DESC"
 
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
 
	var jobs []*domain.Job
	for rows.Next() {
		job := &domain.Job{}
		var locationJSON []byte
		err := rows.Scan(
			&job.ID,
			&job.CompanyID,
			&job.Title,
			&job.Description,
			&job.SalaryMin,
			&job.SalaryMax,
			&job.LocationType,
			&locationJSON,
			&job.EmploymentType,
			&job.Deadline,
			&job.Status,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
 
		if err := json.Unmarshal(locationJSON, &job.Location); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
 }
 
 type JobFilter struct {
	CompanyID uuid.UUID
	Status    string
 }