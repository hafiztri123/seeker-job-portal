package main

import (
    "database/sql"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/hafiztri123/config"
    "github.com/hafiztri123/internal/core/ports"
    "github.com/hafiztri123/internal/core/services"
    "github.com/hafiztri123/internal/handlers"
    "github.com/hafiztri123/internal/middleware"
    "github.com/hafiztri123/internal/repositories/elasticsearch"
    "github.com/hafiztri123/internal/repositories/postgres"
    "github.com/hafiztri123/migrations"
    "github.com/hafiztri123/pkg/database"
    "github.com/hafiztri123/pkg/redis"
    seed "github.com/hafiztri123/scripts"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    redisClient "github.com/redis/go-redis/v9"
)

const BASE_URL = "/api/v1"

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler,
    })

    app.Use(recover.New())
    app.Use(logger.New())
    app.Use(cors.New())

    cfg, db, redisClient := initializeApp()
    defer db.Close()
    defer redisClient.Close()

    setupDatabase(db)
    setupRoutes(app, db, redisClient, cfg)

    log.Fatal(app.Listen(":8080"))
}

// Initialize application dependencies
func initializeApp() (*config.Config, *sql.DB, *redisClient.Client) {
    if err := godotenv.Load("/app/app.env"); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    redisClient := redis.NewRedisClient(cfg.Redis)

    return cfg, db, redisClient
}

// Setup database migrations and seeders
func setupDatabase(db *sql.DB) {
    migrator := migrations.NewMigrator(db)
    if err := migrator.Run(); err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }

    if os.Getenv("APP_ENV") == "development" {
        seeder := seed.NewSeeder(db)
        if err := seeder.SeedAll(); err != nil {
            log.Fatalf("Error seeding database: %v", err)
        }
    }
}

// Setup routes for the application
func setupRoutes(app *fiber.App, db *sql.DB, redisClient *redisClient.Client, cfg *config.Config) {
    authHandler, authService := authHandlerInit(db, redisClient, cfg)
    profileHandler := profileHandlerInit(db)
    searchHandler := searchHandlerInit(db, cfg)
    companyHandler := companyHandlerInit(db)

    api := app.Group(BASE_URL)

    // Public routes
    app.Get("/health", handlers.HealthCheck)

    // Auth routes
    auth := api.Group("/auth")
    auth.Post("/login", middleware.ValidateBody(&handlers.LoginRequest{}), authHandler.Login)
    auth.Post("/register/user", middleware.ValidateBody(&handlers.RegisterRequest{}), authHandler.Register)
    auth.Post("/register/company", middleware.ValidateBody(&ports.CompanyRegisterRequest{}), companyHandler.Register)
    auth.Post("/refresh", middleware.ValidateBody(&handlers.RefreshRequest{}), authHandler.RefreshToken)

    // Search routes
    jobs := api.Group("/jobs")
    jobs.Get("/search", searchHandler.SearchJobs)

    // Protected user routes
    userProtected := api.Group("/user", middleware.AuthMiddleware(authService))
    userProtected.Get("/profile", profileHandler.GetProfile)
    userProtected.Put("/profile", middleware.ValidateBody(&ports.UpdateProfileRequest{}), profileHandler.UpdateProfile)

    // Protected company routes
    companyProtected := api.Group("/company", middleware.AuthMiddleware(authService))
    companyProtected.Get("/profile", companyHandler.GetCompany)
    companyProtected.Put("/profile", middleware.ValidateBody(&ports.CompanyUpdateProfileRequest{}), companyHandler.UpdateProfile)
}

// Initialize auth handler and service
func authHandlerInit(db *sql.DB, redisClient *redisClient.Client, cfg *config.Config) (*handlers.AuthHandler, ports.AuthService) {
    userRepo := postgres.NewUserRepository(db)
    companyRepo := postgres.NewCompanyRepository(db)
    authService := services.NewAuthService(userRepo, redisClient, cfg, companyRepo)
    return handlers.NewAuthHandler(authService), authService
}

// Initialize profile handler
func profileHandlerInit(db *sql.DB) *handlers.ProfileHandler {
    userRepo := postgres.NewUserRepository(db)
    userService := services.NewProfileService(userRepo)
    return handlers.NewProfileHandler(userService)
}

// Initialize search handler
func searchHandlerInit(db *sql.DB, cfg *config.Config) *handlers.SearchHandler {
    searchRepo, err := elasticsearch.NewSearchRepository(&cfg.ElasticSearch)
    if err != nil {
        log.Fatal(err)
    }
    jobRepo := postgres.NewJobRepository(db)
    searchService := services.NewSearchService(*searchRepo, *jobRepo)
    return handlers.NewSearchHandler(searchService)
}

// Initialize company handler
func companyHandlerInit(db *sql.DB) *handlers.CompanyHandler {
    repo := postgres.NewCompanyRepository(db)
    service := services.NewCompanyService(repo)
    return handlers.NewCompanyHandler(service)
}