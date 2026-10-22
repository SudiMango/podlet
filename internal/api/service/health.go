package service

import (
	"context"

	"github.com/SudiMango/podlet/internal/api/repository"
)

type HealthService struct {
	healthRepository *repository.HealthRepository
}

func NewHealthService(healthRepository *repository.HealthRepository) *HealthService {
	return &HealthService{
		healthRepository: healthRepository,
	}
}

func (s *HealthService) IsDatabaseHealthy(ctx context.Context) (bool, error) {
	return s.healthRepository.IsDatabaseHealthy(ctx)
}
