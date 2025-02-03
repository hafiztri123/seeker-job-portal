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
	"github.com/hafiztri123/internal/core/ports"
	"github.com/hafiztri123/internal/core/services"
	"github.com/hafiztri123/internal/handlers"
	"github.com/hafiztri123/internal/middleware"
	"github.com/hafiztri123/internal/repositories/postgres"
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

	authHandler, authService := authHandlerInit(db)
	profileHandler := profileHandlerInit(db)

	profileRoutes(app, profileHandler, authService)
	healthRoutes(app)
	authRoutes(app, authHandler)



	log.Fatal(app.Listen(":8080"))

}

const (
	BASE_URL = "/api/v1"

)

func healthRoutes(app *fiber.App)  {
	app.Get("/health", handlers.HealthCheck)
}

func authRoutes(app *fiber.App, handler *handlers.AuthHandler)  {
	auth := app.Group(BASE_URL + "/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
}


func authHandlerInit(db *sql.DB) (*handlers.AuthHandler, ports.AuthService) {
	authRepo := postgres.NewUserRepository(db)
	authService := services.NewAuthService(authRepo)
	return handlers.NewAuthHandler(authService), authService
}

func profileHandlerInit(db *sql.DB) *handlers.ProfileHandler {
	profileRepo := postgres.NewUserRepository(db)
	profileService := services.NewProfileService(profileRepo)
	return handlers.NewProfileHandler(profileService)
}

func profileRoutes(app *fiber.App, handler *handlers.ProfileHandler, authService ports.AuthService) {
	profile := app.Group(BASE_URL + "/user/profile", middleware.AuthMiddleware(authService))
	profile.Get("/", handler.GetProfile)
	profile.Put("/", handler.UpdateProfile)
}

