package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"time"
)

//go:embed sql/*.sql
var migrationFiles embed.FS

type Migration struct {
    ID        string    `db:"id"`
    Name      string    `db:"name"`
    Applied   time.Time `db:"applied"`
}

type Migrator struct {
    db *sql.DB
}

func NewMigrator(db *sql.DB) Migrator {
    return Migrator{db: db}
}

func (m *Migrator) CreateMigrationsTable() error {
    //uuid-ossp for uuid_generate_v4()
    query := `
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE IF NOT EXISTS migrations (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
            name VARCHAR(255) NOT NULL UNIQUE,
            applied TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `
    _, err := m.db.Exec(query)
    return err
}

func (m *Migrator) GetAppliedMigrations() ([]Migration, error) {
    var migrations []Migration
    cursor, err := m.db.Query("SELECT id, name, applied FROM migrations ORDER BY applied")
    if err != nil {
        return nil, err
    }

    for cursor.Next() {
        var migration Migration
        err = cursor.Scan(&migration.ID, &migration.Name, &migration.Applied)
        if err != nil {
            return nil, err
        }
        migrations = append(migrations, migration)
    }
    return migrations, cursor.Err()
}

func (m *Migrator) ApplyMigration(name string, content string) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    if _, err := tx.Exec(content); err != nil {
        return fmt.Errorf("failed to apply migration %s: %w", name, err)
    }

    if _, err := tx.Exec(
        "INSERT INTO migrations (name) VALUES ($1)",
        name,
    ) ; err != nil {
        return fmt.Errorf("failed to record migration %s: %w", name, err)
    }

    return tx.Commit()
}

func (m *Migrator) RollbackMigration(name string, content string) error {
    tx, err := m.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    if _, err := tx.Exec(content); err != nil {
        return fmt.Errorf("failed to rollback migration %s: %w", name, err)
    }

    if _, err := tx.Exec(
        "DELETE FROM migrations WHERE name = $1",
        name,
    ) ; err != nil {
        return fmt.Errorf("failed to record migration %s: %w", name, err)
    }

    return tx.Commit()
}

func (m *Migrator) Run() error {
    if err := m.CreateMigrationsTable(); err != nil {
        return err
    }

    applied, err := m.GetAppliedMigrations()
    if err != nil {
        return err
    }

    appliedMap := make(map[string]bool)
    for _, m := range applied {
        appliedMap[m.Name] = true
    }

    entries, err := migrationFiles.ReadDir("sql")
    if err != nil {
        return err
    }

    for _, entry := range entries {
        if appliedMap[entry.Name()] {
            continue
        }

        content, err := migrationFiles.ReadFile(fmt.Sprintf("sql/%s", entry.Name()))
        if err != nil {
            return err
        }

        if err := m.ApplyMigration(entry.Name(), string(content)); err != nil {
            return err
        }
    }

    return nil
}