package usecase

import "context"

// HealthUsecase defines the interface for health check business logic
type HealthUsecase interface {
	Check(ctx context.Context) (string, error)
}

type healthUsecase struct{}

// NewHealthUsecase creates a new health usecase
func NewHealthUsecase() HealthUsecase {
	return &healthUsecase{}
}

func (u *healthUsecase) Check(ctx context.Context) (string, error) {
	return "ok", nil
}
