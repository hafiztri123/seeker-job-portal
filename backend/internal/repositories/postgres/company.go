package postgres

import (
	"database/sql"
	"fmt"

	"github.com/hafiztri123/internal/core/domain"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *domain.Company) error {
	query := `
		INSERT INTO companies (
		id,
		email, 
		hashed_password, 
		name
		)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(
		query,
		company.Id,
		company.Email,
		company.HashedPassword,
		company.Name,
	)

	return err
 }

func (r *CompanyRepository) Update(company *domain.Company) error {
	query := `
		UPDATE companies
		SET name = $1, phone_number = $2, location = $3, about = $4, business_category = $5, company_size = $6, profile_picture = $7  ,updated_at = now()
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(
		query,
		company.Name,
		company.PhoneNumber,
		company.Location,
		company.About,
		company.BussinessCategory,
		company.Size,
		company.Id,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("company not found")
	}

	return nil

}

func(r *CompanyRepository) GetAllCompany() ([]*domain.Company, error) {
	query := `
		SELECT id, email, hashed_password, name, phone_number, location, about, business_category, company_size, created_at, updated_at, deleted_at
		FROM companies
		WHERE deleted_at IS NULL
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var companies []*domain.Company
	for rows.Next() {
		company := &domain.Company{}
		if err := rows.Scan(
			&company.Id,
			&company.Email,
			&company.HashedPassword,
			&company.Name,
			&company.PhoneNumber,
			&company.Location,
			&company.About,
			&company.BussinessCategory,
			&company.Size,
			&company.CreatedAt,
			&company.UpdatedAt,
			&company.DeletedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	return companies, nil
}

func(r *CompanyRepository) FindByID(id string) (*domain.Company, error) {
	query := `
		SELECT id, email, hashed_password, name, phone_number, location, about, business_category, company_size, created_at, updated_at, deleted_at
		FROM companies
		WHERE id = $1 AND deleted_at IS NULL
	`

	company :=&domain.Company{}
	err := r.db.QueryRow(query, id).Scan(
		&company.Id,
		&company.Email,
		&company.HashedPassword,
		&company.Name,
		&company.PhoneNumber,
		&company.Location,
		&company.About,
		&company.BussinessCategory,
		&company.Size,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return company, err
}	

func(r *CompanyRepository) FindByEmail(email string) (*domain.Company, error) {
	query := `
		SELECT id, email, hashed_password, name, phone_number, location, about, business_category, company_size, created_at, updated_at, deleted_at
		FROM companies
		WHERE email = $1 AND deleted_at IS NULL
	`

	company :=&domain.Company{}
	err := r.db.QueryRow(query, email).Scan(
		&company.Id,
		&company.Email,
		&company.HashedPassword,
		&company.Name,
		&company.PhoneNumber,
		&company.Location,
		&company.About,
		&company.BussinessCategory,
		&company.Size,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return company, err
}


