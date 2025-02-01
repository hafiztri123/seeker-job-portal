package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/hafiztri123/migrations"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrator := migrations.NewMigrator(db)
	err = migrator.Run()
	if err != nil {
		log.Fatal(err)
	}

}