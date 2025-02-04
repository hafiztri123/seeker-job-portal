package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hafiztri123/config"
	"github.com/hafiztri123/internal/core/ports"
	"github.com/hafiztri123/internal/core/services"
	"github.com/hafiztri123/internal/handlers"
	"github.com/hafiztri123/internal/middleware"
	"github.com/hafiztri123/internal/repositories/postgres"
	"github.com/hafiztri123/migrations"
	"github.com/hafiztri123/pkg/database"
	"github.com/hafiztri123/pkg/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	redisClient "github.com/redis/go-redis/v9"
)

const (
	BASE_URL = "/api/v1"
)

func main() {
	err := godotenv.Load("/home/hafizh/seeker.com/app.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.Connect(config.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	migrator := migrations.NewMigrator(db)

	err = migrator.Run()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	redisClient := redis.NewRedisClient(config.Redis)
	defer redisClient.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	authHandler, authService := authHandlerInit(db, redisClient, config)
	profileHandler := profileHandlerInit(db)
	setupRoutes(app, authHandler, profileHandler, authService)
	log.Fatal(app.Listen(":8080"))

}

func authHandlerInit(db *sql.DB, redisClient *redisClient.Client, config *config.Config) (*handlers.AuthHandler, ports.AuthService) {
	authRepo := postgres.NewUserRepository(db)
	authService := services.NewAuthService(authRepo, redisClient, config)
	return handlers.NewAuthHandler(authService), authService
}

func profileHandlerInit(db *sql.DB) *handlers.ProfileHandler {
	userRepo := postgres.NewUserRepository(db)
	userService := services.NewProfileService(userRepo)
	return handlers.NewProfileHandler(userService)
}

func setupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, profileHandler *handlers.ProfileHandler, authService ports.AuthService) {
	api := app.Group(BASE_URL)

	// Public routes
	app.Get("/health", handlers.HealthCheck)
	auth := api.Group("/auth")
	auth.Post("/login", middleware.ValidateBody(&handlers.LoginRequest{}) ,authHandler.Login)
	auth.Post("/register", middleware.ValidateBody(&handlers.RegisterRequest{}), authHandler.Register)
	auth.Post("/refresh",middleware.ValidateBody(&handlers.RefreshRequest{}) ,authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("/user", middleware.AuthMiddleware(authService))
	protected.Get("/profile", profileHandler.GetProfile)
	protected.Put("/profile", profileHandler.UpdateProfile)
}
