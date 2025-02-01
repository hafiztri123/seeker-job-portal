package migrations

import (
	"database/sql"

	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "postgres://hafizh:Sudarmi12@localhost:5432/seeker?sslmode=disable")

	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func TestMigrator(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	migrator := NewMigrator(db)

	t.Run("CreateMigrationsTable", func(t *testing.T) {
		err := migrator.CreateMigrationsTable()
		assert.NoError(t, err)

		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'migrations'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)
	
	})

	t.Run("GetAppliedMigrations", func(t *testing.T) {
		_, err := db.Exec(
			"INSERT INTO migrations (name) VALUES ($1)",
			"test_migration.sql",
		)

		assert.NoError(t, err)
		migrations, err := migrator.GetAppliedMigrations()
		assert.NoError(t, err)
		// assert.Len(t, migrations, 1)
		assert.Equal(t, "test_migration.sql", migrations[0].Name)
	})

	t.Run("ApplyMigration", func(t *testing.T) {
		testMigration := `
			CREATE TABLE test_table (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				name TEXT
			);
		` 

		err := migrator.ApplyMigration("test_table.sql", testMigration)
		assert.NoError(t, err)

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", "test_table.sql").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT from information_schema.tables
				WHERE table_name = 'test_table'
			)
		`).Scan(&exists)
		assert.NoError(t, err)
		assert.True(t, exists)

		t.Run("RollbackMigration", func(t *testing.T) {
			testRollback := `DROP TABLE test_table;`
			err := migrator.RollbackMigration("test_table.sql", testRollback)
			assert.NoError(t, err)

			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1",
				"test_table.sql").Scan(&count)
			assert.NoError(t, err)
			assert.Equal(t, 0, count)

			var exists bool
			err = db.QueryRow(`
				SELECT EXISTS (
					SELECT from information_schema.tables
					WHERE table_name = 'test_table'
				)
			`).Scan(&exists)
			assert.NoError(t, err)
			assert.False(t, exists)
		})

	})


}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS test_table, migrations`)
	assert.NoError(t, err)
}

