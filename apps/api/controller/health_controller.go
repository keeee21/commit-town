package controller

import (
	"net/http"

	"github.com/keeee21/commit-town/api/usecase"
	"github.com/labstack/echo/v4"
)

type HealthController struct {
	healthUsecase usecase.HealthUsecase
}

type HealthResponse struct {
	Status string `json:"status"`
}

// NewHealthController creates a new health controller
func NewHealthController(healthUsecase usecase.HealthUsecase) *HealthController {
	return &HealthController{
		healthUsecase: healthUsecase,
	}
}

func (h *HealthController) Check(c echo.Context) error {
	status, err := h.healthUsecase.Check(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, HealthResponse{
		Status: status,
	})
}
