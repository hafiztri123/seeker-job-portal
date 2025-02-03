package postgres

import (
	"database/sql"

	"github.com/hafiztri123/internal/core/domain"
	"github.com/hafiztri123/internal/core/ports"
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

func(r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET full_name= $1, phone_number = $2, about = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, user.FullName, user.PhoneNumber, user.About, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ports.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(id string) error {
	_, err := r.db.Exec(`UPDATE users SET deleted_at = now() WHERE id = $1`, id)
	return err
}