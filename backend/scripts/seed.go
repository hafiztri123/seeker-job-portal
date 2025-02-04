package seed

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedUsers() error {
	users := []struct {
		email string
		password string
		fullname string
	}{
		{"Kb6aK@example.com", "password", "John Doe"},
		{"mKb6aK@example.com", "password", "Jane Doe"},
		{"3Nl7X@example.com", "password", "Bob Smith"},
	}

	for _, user := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		_, err := s.db.Exec(`
			INSERT INTO users (id, email, hashed_password, full_name)
			VALUES (uuid_generate_v4(), $1, $2, $3)
			ON CONFLICT (email) DO NOTHING
		`, user.email, string(hashedPassword), user.fullname)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) SeedCompanies() error {
	companies := []struct {
		name string
		email string
		password string
		about string
	}{
		{"Acme Inc.", "Kb6aK@example.com", "password", "We make things"},
		{"XYZ Corp.", "mKb6aK@example.com", "password", "We sell things"},
		{"ABC Corp.", "3Nl7X@example.com", "password", "We buy things"},
		{"Limbus Company", "mKb6aK@example.com", "password", "We sell things"},
		{"ABC Company", "3Nl7X@example.com", "password", "We buy things"},
	}

	for _, company := range companies {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(company.password), bcrypt.DefaultCost)
		_, err := s.db.Exec(`
			INSERT INTO companies (id, name, email, hashed_password, about)
			VALUES (uuid_generate_v4(), $1, $2, $3, $4)
			ON CONFLICT (email) DO NOTHING
		`, company.name, company.email, string(hashedPassword), company.about)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) SeedJobs() error {
	var companyIDs []string
	cursor, err := s.db.Query("SELECT id FROM companies")
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		var id string
		if err := cursor.Scan(&id); err != nil {
			return err
		}
		companyIDs = append(companyIDs, id)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	jobs := []struct {
		title string
		description string
		salaryMin int
		salaryMax int
		locationType string
		employmentType string
		status string
		deadline string
	}{
		{
			title: "Software Engineer",
			description: "We are looking for a talented software engineer to join our team.",
			salaryMin: 50000,
			salaryMax: 80000,
			locationType: "remote",
			employmentType: "full_time",
			status: "open",
			deadline: "2023-06-30",
		},
		{
			title: "Data Analyst",
			description: "We are looking for a talented data analyst to join our team.",
			salaryMin: 40000,
			salaryMax: 60000,
			locationType: "onsite",
			employmentType: "full_time",
			status: "closed",
			deadline: "2023-06-30",
		},
		{
			title: "Sales Representative",
			description: "We are looking for a talented sales representative to join our team.",
			salaryMin: 40000,
			salaryMax: 60000,
			locationType: "onsite",
			employmentType: "full_time",
			status: "open",
			deadline: "2023-06-30",
		},
		{
			title: "Product Manager",
			description: "We are looking for a talented product manager to join our team.",
			salaryMin: 60000,
			salaryMax: 90000,
			locationType: "onsite",
			employmentType: "full_time",
			status: "open",
			deadline: "2023-06-30",
		},
		{
			title: "UX Designer",
			description: "We are looking for a talented UX designer to join our team.",
			salaryMin: 50000,
			salaryMax: 80000,
			locationType: "onsite",
			employmentType: "full_time",
			status: "closed",
			deadline: "2023-06-30",
		},
	}

	for i, j := range jobs {
		_, err := s.db.Exec(
			`INSERT INTO jobs (
				id,
				company_id,
				title, 
				description, 
				salary_min, 
				salary_max, 
				location_type, 
				employment_type, 
				status, 
				deadline
			) VALUES
			(uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			companyIDs[i%len(companyIDs)],
			j.title, 
			j.description, 
			j.salaryMin, 
			j.salaryMax, 
			j.locationType, 
			j.employmentType, 
			j.status, 
			j.deadline,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) SeedSkills() error {
	skills := []struct {
		name string
		category string
	}{
		{"Python", "Programming Language"},
		{"Java", "Programming Language"},
		{"JavaScript", "Programming Language"},
		{"React", "Frontend Framework"},
		{"Django", "Web Framework"},
	}

	for _, skill := range skills {
		_, err := s.db.Exec(`
			INSERT INTO skills (id, name, category)
			VALUES (uuid_generate_v4(), $1, $2)
			ON CONFLICT (name) DO NOTHING
		`, skill.name, skill.category)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) SeedJobApplications() error {
	var userID, jobID string
	err := s.db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
	if err != nil {
		return err
	}

	err = s.db.QueryRow("SELECT id FROM jobs LIMIT 1").Scan(&jobID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO job_applications(
			id, user_id, job_id,
			status, cover_letter
		)
		VALUES (
			uuid_generate_v4(), $1, $2, 
			'pending', 'I am excited to apply for this job!'
		)
		ON CONFLICT (user_id, job_id) DO NOTHING
	`, userID, jobID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Seeder) SeedAll() error {
    if err := s.SeedUsers(); err != nil {
        return err
    }
    if err := s.SeedCompanies(); err != nil {
        return err
    }
    if err := s.SeedSkills(); err != nil {
        return err
    }
    if err := s.SeedJobs(); err != nil {
        return err
    }
    return s.SeedJobApplications()
}