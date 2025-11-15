package router

import (
	"github.com/keeee21/commit-town/api/controller"
	"github.com/labstack/echo/v4"
)

// SetupRoutes sets up all API routes
func SetupRoutes(e *echo.Echo, healthController *controller.HealthController) {
	// Health check
	e.GET("/health", healthController.Check)

	// Add more routes here as needed
}
