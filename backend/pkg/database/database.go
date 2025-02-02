package database

import (
	"database/sql"

	"github.com/hafiztri123/config"
)

func Connect(cfg *config.DatabaseConfig) (*sql.DB, error) {
	
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}