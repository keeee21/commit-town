package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/keeee21/commit-town/api/controller"
	"github.com/keeee21/commit-town/api/db"
	"github.com/keeee21/commit-town/api/repository"
	"github.com/keeee21/commit-town/api/router"
	"github.com/keeee21/commit-town/api/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database, err := db.NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)

	// Initialize usecases
	healthUsecase := usecase.NewHealthUsecase()
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize controllers
	healthController := controller.NewHealthController(healthUsecase)
	userController := controller.NewUserController(userUsecase)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes
	router.SetupRoutes(e, healthController, userController)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
