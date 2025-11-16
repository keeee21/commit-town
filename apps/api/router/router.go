package router

import (
	"github.com/keeee21/commit-town/api/controller"
	"github.com/labstack/echo/v4"
)

// SetupRoutes sets up all API routes
func SetupRoutes(e *echo.Echo, healthController *controller.HealthController, userController *controller.UserController) {
	// Health check
	e.GET("/health", healthController.Check)

	// User routes
	api := e.Group("/api")
	api.POST("/users", userController.UpsertUser)
}
