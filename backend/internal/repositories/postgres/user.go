package postgres

import (
	"database/sql"

	"github.com/hafiztri123/internal/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	
	query := `
		INSERT INTO users (id, email, hashed_password, full_name)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.HashedPassword,
		user.FullName,
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, hashed_password, full_name 
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL
	`

	user :=&domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return user, err
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	query := `
		SELECT id, email, hashed_password, full_name
		FROM users
		WHERE id = $1 ANd deleted_at IS NULL
	`

	user :=&domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.FullName,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}

	return user, err
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec(`UPDATE users SET deleted_at = now() WHERE id = $1`, id)
	return err
}