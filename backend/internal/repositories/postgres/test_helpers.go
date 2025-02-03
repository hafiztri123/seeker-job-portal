package postgres

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

func runTestMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE users (
			id UUID PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			hashed_password VARCHAR(255) NOT NULL,
			full_name VARCHAR(255) NOT NULL,
			deleted_at TIMESTAMP
		)
	`)

	return err
}

func cleanupTestDB(db *sql.DB, t *testing.T) {
	_, err := db.Exec(`DROP TABLE IF EXISTS users`)
	assert.NoError(t, err)
}

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://hafizh:Sudarmi12@localhost:5432/seeker?sslmode=disable")
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	cleanupTestDB(db, t)

	err = runTestMigrations(db)
	require.NoError(t, err)

	return db
}
