package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hafiztri123/internal/handlers"
	"github.com/hafiztri123/internal/middleware"
	"github.com/hafiztri123/migrations"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
func main() {
	err := godotenv.Load("/home/hafizh/seeker.com/app.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New()) //ALLOW ALL ORIGINS. TODO: CUSTOMIZE LATER

	app.Get("/health", handlers.HealthCheck)
	log.Fatal(app.Listen(":8080"))

}